package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping 测试接口
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
