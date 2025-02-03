package schedule

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	c = cron.New()
)

// Register 注册并运行全局定时任务并在优雅启停中处理定时任务的关闭
func Register(spec string, cmd func()) {
	_, err := c.AddFunc(spec, cmd)
	if err != nil {
		zap.L().Error("全局定时任务注册失败", zap.Error(err))
		return
	}
}

func Start() {
	c.Start()
	zap.L().Info("全局定时任务启动成功")
}

func Stop() {
	c.Stop()
	zap.L().Info("全局定时任务已关闭")

}
