package testx

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"joeytest.com/resk/infra"
	"joeytest.com/resk/infra/base"
)

func init() {
	// 获取程序运行文件所在 的路径
	file := kvs.GetCurrentFilePath("../../brun/config.ini", 1)
	// 加载和解析配置文件
	conf := ini.NewIniFileConfigSource(file)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	app := infra.New(conf)
	app.Start()
}
