package services

import (
	"ad-server/models"
	"fmt"

	"github.com/astaxie/beego/orm"
)

type ImpressionService struct {
}

func (c *ImpressionService) TrackImpression(deliveryID int64, adID int64) int64 {
	o := orm.NewOrm()
	o.Using("default")

	//created impression
	impressionObj := new(models.Impression)
	impressionObj.DeliveryId = deliveryID
	impressionObj.ImpressionType = 1 //1: Impression 2: Viewable Impression
	impressionObj.Valid = true

	impressionID, err := o.Insert(impressionObj)

	if err == nil {
		//create or update ad_delivery summary
		deliverySummaryObj := models.AdDeliverySummary{Adid: adID}

		if o.Read(&deliverySummaryObj, "Adid") == nil {
			deliverySummaryObj.ImpressionCount++

			if num, err := o.Update(&deliverySummaryObj, "ImpressionCount"); err == nil {
				fmt.Println(num)
			}
		} else {
			deliverySummaryObj := new(models.AdDeliverySummary)
			deliverySummaryObj.Adid = adID
			deliverySummaryObj.ClickCount = 0
			deliverySummaryObj.ImpressionCount = 1
			o.Insert(deliverySummaryObj)
		}

		// create charge
		var cService = ChargeService{}
		cService.CreateCharge(deliveryID, adID, "impression")

	}

	return impressionID
}
