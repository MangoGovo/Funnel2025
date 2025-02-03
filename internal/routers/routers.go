package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
}
