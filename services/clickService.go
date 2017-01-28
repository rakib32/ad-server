package services

import (
	"ad-server/models"
	"fmt"

	"github.com/astaxie/beego/orm"
)

type ClickService struct {
}

func (c *ClickService) TrackClick(deliveryID int64, adID int64) int64 {
	o := orm.NewOrm()
	o.Using("default")

	//created impression
	clickObj := new(models.Click)
	clickObj.DeliveryId = deliveryID
	clickObj.Valid = true

	clickID, err := o.Insert(clickObj)

	if err == nil {
		//create or update ad_delivery summary
		deliverySummaryObj := models.AdDeliverySummary{Adid: adID}

		if o.Read(&deliverySummaryObj, "Adid") == nil {
			deliverySummaryObj.ClickCount++

			if num, err := o.Update(&deliverySummaryObj, "ClickCount"); err == nil {
				fmt.Println(num)
			}
		} else {
			deliverySummaryObj := new(models.AdDeliverySummary)
			deliverySummaryObj.Adid = adID
			deliverySummaryObj.ClickCount = 1
			deliverySummaryObj.ImpressionCount = 0
			o.Insert(deliverySummaryObj)
		}

		// create charge
		var cService = ChargeService{}
		cService.CreateCharge(deliveryID, adID, "click")

	}

	return clickID
}
