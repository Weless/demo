package base

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	iris_recover "github.com/kataras/iris/v12/middleware/recover"
	"github.com/sirupsen/logrus"
	"joeytest.com/resk/infra"
	"time"
)

var irisApplication *iris.Application

func Iris() *iris.Application {
	return irisApplication
}

type IrisServerStarter struct {
	infra.BaseStarter
}

func (i *IrisServerStarter) Init(ctx infra.StarterContext) {
	// 创建iris application实例
	irisApplication = initIris()
	// 日志组件配置和扩展
	logger := irisApplication.Logger()
	logger.Install(logrus.StandardLogger())
}

func (i *IrisServerStarter) Start(ctx infra.StarterContext) {
	// 把路有信息打印到控制台方便查看
	routes := Iris().GetRoutes()
	for _, r := range routes {
		logrus.Info(r.Trace)
	}
	// 启动Iris
	port := ctx.Props().GetDefault("app.server.port", "8000")
	err := Iris().Run(iris.Addr(":" + port))
	if err != nil {
		logrus.Error(err)
	}
}

func (i *IrisServerStarter) StartBlocking() bool {
	return true
}

func initIris() *iris.Application {
	app := iris.New()
	// 主要中间件的配置：recover，日志输出中间件的定义
	app.Use(iris_recover.New())
	cfg := logger.Config{
		Status: true,
		IP:     true,
		Method: true,
		Path:   true,
		Query:  true,
		LogFunc: func(endTime time.Time, latency time.Duration,
			status, ip, method, path string, message interface{},
			headerMessage interface{}) {
			app.Logger().Infof("| %s | %s | %s | %s | %s | %s | %s | %s |",
				endTime.Format("2006-01-02 15:04:05"),
				latency.String(),
				status,
				ip,
				method,
				path,
				headerMessage,
				message,
			)
		},
	}
	app.Use(logger.New(cfg))
	return app
}
