package controllers

import (
	"ad-server/helpers"
	"ad-server/requests"
	"ad-server/responses"
	"ad-server/services"
	"strconv"

	"github.com/astaxie/beego"
)

type AdController struct {
	beego.Controller
}

func (this *AdController) ServeAds() {
	os := this.GetString("os")
	version := this.GetString("version")
	model := this.GetString("model")
	token := this.GetString("token")
	adspaceID := this.GetString("placement_key")
	limit := this.GetString("limit")

	var utility = helpers.Utility{}
	clientip, _ := utility.ExternalIP()

	beego.Debug("client:" + clientip)

	var errorMsgs []responses.Error

	if os == "" {
		buildErrorMessage(&errorMsgs, "Os")
	}
	if version == "" {
		buildErrorMessage(&errorMsgs, "Version")
	}

	if token == "" {
		buildErrorMessage(&errorMsgs, "Token")
	}
	if adspaceID == "" {
		buildErrorMessage(&errorMsgs, "Adspace Id")
	}
	if limit == "" {
		limit = "2" // default:2
	}

	if len(errorMsgs) > 0 {
		this.Data["json"] = errorMsgs
		this.ServeJSON()
	} else {
		requestParams := requests.Ad{}
		i, _ := strconv.ParseInt(adspaceID, 10, 64)
		j, _ := strconv.Atoi(limit)
		requestParams.AdspaceID = i
		requestParams.Model = model
		requestParams.Os = os
		requestParams.Token = token
		requestParams.Version = version
		requestParams.Limit = j
		requestParams.ClientIP = clientip

		adServ := services.AdService{}
		results := adServ.DeliveryAds(requestParams)

		this.Data["json"] = &results
		this.ServeJSON()
	}

}
func buildErrorMessage(errorMsgs *[]responses.Error, fieldName string) {
	errorMsg := responses.Error{}
	errorMsg.Message = fieldName + " is required"
	errorMsg.ErrorType = "Validation"

	*errorMsgs = append(*errorMsgs, errorMsg)
}

func (this *AdController) SaveImpression() {

	adid, _ := strconv.ParseInt(this.Ctx.Input.Param(":adid"), 10, 64)
	deliveryid, _ := strconv.ParseInt(this.Ctx.Input.Param(":deliveryid"), 10, 64)

	impServ := services.ImpressionService{}
	id := impServ.TrackImpression(deliveryid, adid)

	var message = responses.Message{}

	if id > 0 {
		message.Message = "Impression has been tracked successfully."

	} else {
		message.Message = "Failed to track the impression."
	}

	this.Data["json"] = message
	this.ServeJSON()
}

func (this *AdController) SaveClick() {
	adid, _ := strconv.ParseInt(this.Ctx.Input.Param(":adid"), 10, 64)
	deliveryid, _ := strconv.ParseInt(this.Ctx.Input.Param(":deliveryid"), 10, 64)

	clickServ := services.ClickService{}
	id := clickServ.TrackClick(deliveryid, adid)

	var message = responses.Message{}

	if id > 0 {
		message.Message = "Click has been tracked successfully."

	} else {
		message.Message = "Failed to track the Click."
	}

	this.Data["json"] = message
	this.ServeJSON()
}
