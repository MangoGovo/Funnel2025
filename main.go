package main

import (
	"funnel/internal/routers"
	"funnel/pkg/config"
	_ "funnel/pkg/log"
	"funnel/pkg/schedule"
	"funnel/pkg/server"
	"github.com/gin-gonic/gin"
)

var debug = config.Config.GetBool("server.debug")
var port = ":" + config.Config.GetString("server.port")

func main() {
	// TODO Health Check
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	schedule.Start()
	server.Run(routers.R, port)
}
