package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Ping 测试接口
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
