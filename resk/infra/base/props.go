package base

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"joeytest.com/resk/infra"
)

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	return props
}

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	logrus.Info("初始化配置完成")
}
