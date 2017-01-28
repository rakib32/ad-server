package models

import "github.com/astaxie/beego/orm"

type AdDeliverySummary struct {
	DeliverySummaryId int64 `orm:"pk"`
	Adid              int64
	ImpressionCount   int
	ClickCount        int64
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(AdDeliverySummary))
}
