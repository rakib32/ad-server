package models

import "github.com/astaxie/beego/orm"

type Charge struct {
	ChargeId       int64 `orm:"pk"`
	ChargeAmount   float64
	DeliveryId     int64
	ChargeType     string
	PublisherChage float64
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(Charge))
}
