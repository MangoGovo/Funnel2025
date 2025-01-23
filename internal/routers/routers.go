package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Init 初始化路由
func Init(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
}
