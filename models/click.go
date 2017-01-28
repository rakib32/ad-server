package models

import "github.com/astaxie/beego/orm"

type Click struct {
	ClickId    int64 `orm:"pk"`
	Valid      bool
	DeliveryId int64
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(Click))
}
