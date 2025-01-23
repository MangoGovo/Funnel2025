package response

import (
	"funnel/internal/apiException"
	"funnel/pkg/log"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Json 返回json格式数据
func Json(c *gin.Context, httpStatusCode int, code int, msg string, data any) {
	c.JSON(httpStatusCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// JsonSuccess 返回成功json格式数据
func JsonSuccess(c *gin.Context, data any) {
	Json(c, http.StatusOK, 200, "OK", data)
}

// JsonError 返回错误json格式数据
func JsonError(c *gin.Context, code int, msg string) {
	Json(c, http.StatusOK, code, msg, nil)
}

// AbortWithException 用于返回自定义错误信息
func AbortWithException(c *gin.Context, apiError *apiException.Error, err error) {
	logError(c, apiError, err)
	_ = c.AbortWithError(200, apiError) //nolint:errcheck
}

// logError 记录错误日志
func logError(c *gin.Context, apiErr *apiException.Error, err error) {
	// 构建日志字段
	logFields := []zap.Field{
		zap.Int("error_code", apiErr.Code),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("ip", c.ClientIP()),
		zap.Error(err), // 记录原始错误信息
	}
	log.GetLogFunc(apiErr.Level)(apiErr.Msg, logFields...)
}
