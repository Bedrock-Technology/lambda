package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

var (
	vmPool = &sync.Pool{
		New: func() any {
			vm := goja.New()
			vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", false))
			for k, v := range injections {
				vm.Set(k, v)
			}
			return vm
		},
	}
)

func setLogLevel() {
	var l slog.Level
	l.UnmarshalText([]byte(cfg.LogLevel))
	slog.SetLogLoggerLevel(l)
}

func main() {
	loadConfig()
	setLogLevel()
	loadServicesAndWatch(cfg.ServicesDir)

	e := gin.Default()
	e.GET("/version", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, buildInfo()) })
	e.Group(cfg.APIPrefix).Any("/*service", serviceHandler)
	e.Run(cfg.Listen)
}

type rawRequest struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func serviceHandler(c *gin.Context) {
	serviceName := strings.TrimPrefix(c.Request.URL.Path, cfg.APIPrefix)
	serviceName = strings.TrimPrefix(serviceName, "/")
	serviceName += ".js"

	slog.Debug("serviceHandler()", "serviceName", serviceName)

	servicesLock.RLock()
	service, ok := services[serviceName]
	servicesLock.RUnlock()

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": "Service not found"})
		return
	}

	body := ""
	if c.Request.Body != nil {
		body = string(lo.Must(io.ReadAll(c.Request.Body)))
	}

	r := rawRequest{
		Method:  c.Request.Method,
		Path:    c.Request.URL.Path,
		Query:   c.Request.URL.Query(),
		Headers: c.Request.Header,
		Body:    body,
	}

	vm := vmPool.Get().(*goja.Runtime)
	defer vmPool.Put(vm)

	vm.Set("req", r)

	result, err := vm.RunProgram(service.Program)
	if err != nil {
		resp := make(gin.H)

		errMsg := err.Error()
		if ex, ok := err.(*goja.Exception); ok {
			errMsg = ex.Value().String()
		}

		if isJSON([]byte(errMsg)) {
			var x any
			json.Unmarshal([]byte(errMsg), &x)
			resp["err"] = x
		} else {
			resp["err"] = errMsg
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		return
	}

	respRaw := result.Export().(string)
	if isJSON([]byte(respRaw)) {
		c.Writer.Header().Set("Content-Type", "application/json")
	}
	c.String(http.StatusOK, "%s", respRaw)
	c.Abort()
}

func isJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}
