package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"funnel/pkg/redis"
	"funnel/pkg/schedule"
	"go.uber.org/zap"
)

// Run 运行 http 服务器
func Run(handler http.Handler, addr string) {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
	}

	// 启动服务器协程
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("服务遇到了异常", zap.Error(err))
		}
	}()

	// 阻塞并监听结束信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("正在关闭服务...")

	// 关闭服务器（5秒超时时间）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("服务关闭失败", zap.Error(err))
	}

	// 关闭 Redis 客户端
	if err := redis.Client.Close(); err != nil {
		zap.L().Error("Redis客户端关闭失败", zap.Error(err))
	}

	// 关闭定时任务
	schedule.Stop()

	zap.L().Info("服务已关闭")
}
