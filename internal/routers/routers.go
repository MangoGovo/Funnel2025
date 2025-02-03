package routers

import (
	"net/http"

	"funnel/internal/midwares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var R = gin.Default()

// init 初始化路由
func init() {
	R.Use(cors.Default())
	R.Use(midwares.ErrHandler())
	R.NoMethod(midwares.HandleNotFound)
	R.NoRoute(midwares.HandleNotFound)
	R.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
}
