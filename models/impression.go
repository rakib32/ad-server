package models

import "github.com/astaxie/beego/orm"

type Impression struct {
	ImpressionId   int64 `orm:"pk"`
	Valid          bool
	ImpressionType int
	DeliveryId     int64
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(Impression))
}
