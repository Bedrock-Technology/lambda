package main

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
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

	r := rawRequest{
		Method:  c.Request.Method,
		Path:    c.Request.URL.Path,
		Headers: c.Request.Header,
		Body:    string(lo.Must(io.ReadAll(c.Request.Body))),
	}

	vm := service.VM
	for k, v := range injections {
		vm.Set(k, v)
	}
	vm.Set("req", r)

	result, err := vm.RunProgram(service.Program)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	respRaw := result.Export().(string)
	c.String(http.StatusOK, "%s", respRaw)
	c.Abort()
}
