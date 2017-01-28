package models

import "github.com/astaxie/beego/orm"

type AdResource struct {
	AdResourceId int64 `orm:"pk"`
	ResourceHtml string
	BannerLink   string
	Width        string
	Height       string
	Adid         string
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(AdResource))
}
