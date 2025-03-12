package main

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
)

var (
	vmPool = &sync.Pool{
		New: func() any {
			vm := goja.New()
			vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", false))
			return vm
		},
	}
)

func main() {
	loadConfig()

	cfgLock.RLock()
	serviceDir, apiPrefix, listen, dsn := cfg.ServicesDir, cfg.APIPrefix, cfg.Listen, cfg.PostgresDSN
	cfgLock.RUnlock()

	if err := loadDatabase(dsn); err != nil {
		slog.Error("loadDatabase()", "err", err)
	}

	loadServicesAndWatch(serviceDir)

	e := gin.Default()
	e.GET("/version", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, buildInfo()) })
	e.Group(apiPrefix).Any("/*service", func(ctx *gin.Context) {
		ctx.Set("apiPrefix", apiPrefix)
		serviceHandler(ctx)
	})
	e.Run(listen)
}

func serviceHandler(c *gin.Context) {
	serviceName := strings.TrimPrefix(c.Request.URL.Path, c.GetString("apiPrefix"))
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

	vm := vmPool.Get().(*goja.Runtime)
	defer vmPool.Put(vm)

	vm.Set("LambdaHelper", injectorFor(vm, c))

	result, err := vm.RunProgram(service.Program)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	respRaw := result.Export().(string)
	c.String(http.StatusOK, "%s", respRaw)
	c.Abort()
}
