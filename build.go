package main

import (
	"github.com/gin-gonic/gin"
)

var (
	Version   = "Unknown"
	GoVersion = "Unknown"
	GitHash   = "Unknown"
	BuildTime = "Unknown"
	OSArch    = "Unknown"
)

func buildInfo() gin.H {
	return gin.H{
		"lambda": Version,
		"go":     GoVersion,
		"commit": GitHash,
		"built":  BuildTime,
		"arch":   OSArch,
	}
}
