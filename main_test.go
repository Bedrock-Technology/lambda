package main

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkVm(b *testing.B) {
	loadConfig()
	setLogLevel()
	loadServicesAndWatch(cfg.ServicesDir)

	gin.SetMode(gin.ReleaseMode)

	e := gin.Default()
	e.GET("/version", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, buildInfo()) })
	e.Group(cfg.APIPrefix).Any("/*service", serviceHandler)

	runRequest(b, e, http.MethodGet, "/services/utils/dump")
}

func runRequest(B *testing.B, r *gin.Engine, method, path string) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}

type mockWriter struct {
	headers http.Header
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		http.Header{},
	}
}

func (m *mockWriter) Header() (h http.Header) {
	return m.headers
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockWriter) WriteHeader(int) {}
