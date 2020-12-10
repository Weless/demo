package base

import (
	"github.com/sirupsen/logrus"
	"joeytest.com/resk/infra"
)

var apiInitializerRegister *infra.InitializerRegister = new(infra.InitializerRegister)

// 注册WEB API初始化对象
func RegisterApi(ai infra.Initializer) {
	apiInitializerRegister.Register(ai)
}

// 获取注册的web api初始化对象
func GetApiInitializers() []infra.Initializer {
	return apiInitializerRegister.Initializers
}

type WebApiStarter struct {
	infra.BaseStarter
}

func (w *WebApiStarter) Setup(ctx infra.StarterContext) {
	logrus.Info("apiInitializerRegister:", apiInitializerRegister)
	logrus.Info(len(apiInitializerRegister.Initializers))
	for _, v := range GetApiInitializers() {
		logrus.Info("initializer:", v)
		v.Init()
	}
}
