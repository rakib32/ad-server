package models

import "time"
import "github.com/astaxie/beego/orm"

type Delivery struct {
	DeliveryId            int64 `orm:"pk"`
	Timestamp             time.Time
	Valid                 bool
	ClientIp              string
	UserUid               string `orm:"null"`
	UserMarketid          int64  `orm:"null"`
	UserLocationLatitude  float64
	UserLocationLongitude float64
	Adid                  int64
	AdspaceId             int64
	PlatformId            int64  `orm:"null"`
	DeviceOsId            *int64 `orm:"null"`
	DeviceModelId         *int64 `orm:"null"`
	GeoRegionId           *int64 `orm:"null"`
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(Delivery))
}
