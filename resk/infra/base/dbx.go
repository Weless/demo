package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/tietang/props/kvs"
	"joeytest.com/resk/infra"
)

var db *dbx.Database

func DbxDatabase() *dbx.Database {
	return db
}

type DbxDatabaseStarter struct {
	infra.BaseStarter
}

func (d *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {
	conf := ctx.Props()
	// 数据库配置
	settings := dbx.Settings{}
	err := kvs.Unmarshal(conf, &settings, "mysql")
	if err != nil {
		panic(err)
	}
	settings.Options["parseTime"] = "true"
	//logrus.Infof("%+v\n", settings)
	logrus.Info("mysql.conn url:", settings.ShortDataSourceName())
	db, err = dbx.Open(settings)
	if err != nil {
		panic(err)
	}
}
