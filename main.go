package main

import (
	"funnel/internal/midwares"
	"funnel/internal/routers"
	"funnel/pkg/config"
	_ "funnel/pkg/log"
	"funnel/pkg/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// TODO Health Check

	// 如果配置文件中开启了调试模式
	if !config.Config.GetBool("server.debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	routers.Init(r)

	server.Run(r, ":"+config.Config.GetString("server.port"))
}
