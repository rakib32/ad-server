package services

import (
	"ad-server/models"

	"github.com/astaxie/beego/orm"
)

type ChargeService struct {
}

func (c *ChargeService) CreateCharge(deliveryID int64, adID int64, cType string) int64 {
	o := orm.NewOrm()
	o.Using("default")

	var chargeID int64

	var adServ = AdService{}
	ad := adServ.GetAdById(adID)

	if (ad.BidType == "CPM" && cType == "impression") || (ad.BidType == "CPC" && cType == "click") {
		//60% for publisher and 40% for our service
		chargeObj := new(models.Charge)

		if ad.BidType == "CPM" {
			chargeObj.ChargeAmount = ad.BidAmount / float64(1000) //CPM: calculate per impression charge bidAmount/1000
		} else {
			chargeObj.ChargeAmount = ad.BidAmount //CPC: calculate per click
		}

		chargeObj.PublisherChage = chargeObj.ChargeAmount * 0.6 // 60% for publisher
		chargeObj.ChargeType = "SERVICE"
		chargeObj.DeliveryId = deliveryID
		id, _ := o.Insert(chargeObj)
		chargeID = id

	}

	return chargeID
}
