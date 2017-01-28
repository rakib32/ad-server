package main

import (
	_ "ad-server/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root123@/ad_server_db?charset=utf8")
}

func main() {
	orm.Debug = true
	beego.Run()
}
