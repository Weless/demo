package resk

import (
	"joeytest.com/resk/infra"
	"joeytest.com/resk/infra/base"
)

// 手动注册以管理注册顺序
func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
}
