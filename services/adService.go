package services

import (
	"ad-server/models"
	"ad-server/requests"
	"math/rand"
	"sort"
	"time"

	"strconv"

	"ad-server/responses"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ByScore []AdInfo

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }

type AdService struct {
}

type AdInfo struct {
	Adid     int64
	Weight   int64
	Priority int64
	Score    int64
}

//# ad  selection algorithm
//# output=json&os=android&version=6.0&model=SGS6&token=1234&placement_key=tokenkey&limit=3
//# 1. Find all the active ads in the System. Then filter the results based on OS and Version (I.E if os=Android then it will list ads which are targeted for Android)
//# 2. Now we will have the filtered results. Initially set adPriority=0 for all ads.
//# 3. Check if there is any model target for some ads, if yes then increment the adPriority to make ad score high.(High score high probability. iOS or Android app has to sent model as param)
//# 4. Check if there is any adspace target for some ads, if yes then increment the adPriority to make ad score high.(High score high probability)
//# 5. Then check if there is any geo target for some ads. if yes then increment the adPriority to make ad score high.(High score high probability. From the user ip address we can locate the region of the ad request.)
//# 6. Now for each ad calculate weight = (maxImpressions+(adPriority*baseWeight) - servedImpressions)/(adEndtime-adStartTime)  (baseWeight = 500 impresion weight )
//# 7. Then assign score for each ad, score = (weight/totalAllAdsWeight)*100
//# 8. if limit = 1 then select top 5 ads (sort the results by score desc).
//# 8. 			Generate random number to select winning ad with range 1 to 5
//# 9. ELSE
//# 10.			select top n ads (order by score desc). where n=limit*2 (Select results 2 times of the limit )
//# 11.         shuffle the results based on random number.
//#	12.			then select top n results where n=Limit
//# 13.
//# 14.END
//# 15.Create delivery in delivery table for each ad
//# 16.Send the adList to the requester.
//# Future improvement:
//#	    adspace historical analysis score can be added
//#     global cache or hasmaps to hold the wieght of each ads
//#     updated cache if there is any new ad
//# 	Concurrency handler
//# AFTER LINE# 10:Split n ads into 2 blocks.from first block genereate random numbers to select winning ads (80% ads will served from first block and 20% from 2nd block).From 2nd block genereate random numbers to select winning ads (20% For example if Limit is 4 then 3 ads will be selected from first block and 1 //# ad will be selected for 2nd block)
func (c *AdService) DeliveryAds(params requests.Ad) []responses.Ad {
	var results []responses.Ad
	var adInfos []AdInfo

	qb, _ := orm.NewQueryBuilder("mysql")

	versionID := getVersionId(params.Version, params.Os)

	// Construct query object
	qb = qb.Select("ad.adid",
		"ROUND(((ad.ad_max_impression - IFNULL(ad_delivery_summary.impression_count, 0)) / DATEDIFF(end_time, start_time))) AS weight", "1 AS priority", "0 AS score"). //ad_delivery_summary table for quick stats
		From("ad").
		InnerJoin("platform").On("ad.platform_id = platform.platform_id").
		LeftJoin("ad_delivery_summary").On("ad.adid = ad_delivery_summary.adid").
		Where("current_timestamp() >= ad.start_time").
		And("current_timestamp() <= ad.end_time").
		And("ad.state = 'ACTIVE'").
		And("LOWER(platform.platform_name) = LOWER(?)")

	if versionID > 0 {
		qb = qb.And("ad.os_min_version_id <= ?").
			And("ad.os_max_version_id >= ?")
	}

	sql := qb.String()
	o := orm.NewOrm()
	o.Using("default")

	if versionID > 0 {
		o.Raw(sql, params.Os, versionID, versionID).QueryRows(&adInfos)
	} else {
		o.Raw(sql).QueryRows(&adInfos)
	}

	if len(adInfos) > 0 {
		//Not Implemented model and geo target
		//Implemented adspace target
		checkAdspaceTarget(&adInfos, params.AdspaceID)

		// calculate score
		adInfos = calculateWeightedScore(adInfos)

		//sort the result by score
		sort.Sort(ByScore(adInfos))
		var adList []int64

		// select top n ads
		for index, value := range adInfos {
			adList = append(adList, value.Adid)

			if params.Limit == 1 && index >= 5 {
				break
			} else if (params.Limit > 1) && (index >= (params.Limit * 2)) {
				break
			}
		}

		adInfos = nil

		//suffle the results and select top n ads where n=Limit
		shuffle(adList)

		minValue := params.Limit
		len := len(adList)

		if minValue > len {
			minValue = len
		}

		for i := 0; i < minValue; i++ {
			deliveryID := createDelivery(adList[i], params)
			appendResults(&results, deliveryID, adList[i])
		}
	}

	return results
}

func appendResults(results *[]responses.Ad, deliveryID int64, adid int64) {
	var tmpAd = responses.Ad{}
	var image = responses.Image{}

	var adServ = AdService{}
	ad := adServ.GetAdById(adid)

	resource := getResouceByAdId(adid)

	tmpAd.Adid = adid
	tmpAd.ActionText = ad.ActionText
	tmpAd.AdName = ad.AdName
	tmpAd.AppStoreURL = ad.PreviewUrl
	t := strconv.FormatInt(deliveryID, 10)
	a := strconv.FormatInt(adid, 10)
	tmpAd.ImpressionURL = "http://localhost:8080/api/v1/ad/impression/" + a + "/" + t
	tmpAd.ClickURL = "http://localhost:8080/api/v1/ad/click/" + a + "/" + t

	var banner = responses.ImageDetails{}
	banner.Url = resource.BannerLink
	banner.Height = resource.Height
	banner.Width = resource.Width

	image.Banner = banner
	tmpAd.Images = image

	*results = append(*results, tmpAd)
}

func (c *AdService) GetAdById(adid int64) models.Ad {
	var ad models.Ad
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Ad))
	qs.Filter("adid", adid).One(&ad)
	return ad
}

func getResouceByAdId(adid int64) models.AdResource {
	var resource models.AdResource
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.AdResource))
	qs.Filter("adid", adid).One(&resource)
	return resource
}

func createDelivery(adID int64, params requests.Ad) int64 {
	o := orm.NewOrm()
	o.Using("default")

	deliveryObj := new(models.Delivery)
	deliveryObj.Adid = adID
	deliveryObj.AdspaceId = params.AdspaceID
	deliveryObj.ClientIp = params.ClientIP
	deliveryObj.GeoRegionId = nil
	deliveryObj.UserMarketid = 1

	if params.Os == "android" {
		deliveryObj.PlatformId = int64(200)
	} else {
		deliveryObj.PlatformId = int64(100)
	}

	deliveryID, _ := o.Insert(deliveryObj)

	return deliveryID
}

func shuffle(arr []int64) {
	for i := range arr {
		t := time.Now()
		rand.Seed(int64(t.Nanosecond())) // no shuffling without this line
		j := rand.Intn(i + 1)
		beego.Debug(j)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func calculateWeightedScore(adInfos []AdInfo) []AdInfo {
	var baseImpWeight int64 = 500
	var totalWeight int64

	for _, value := range adInfos {
		totalWeight = totalWeight + value.Weight
	}

	for i, value := range adInfos {
		value.Score = int64(float64(value.Weight+(baseImpWeight*value.Priority)) / float64(totalWeight) * 100)
		adInfos[i].Score = value.Score
	}

	return adInfos
}

func checkAdspaceTarget(adInfos *[]AdInfo, adspaceID int64) {
	var adids []int64
	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("ad_target_adspace.adid").
		From("ad_target_adspace").
		Where("ad_target_adspace.adspace_id = ?")

	sql := qb.String()
	o := orm.NewOrm()
	o.Using("default")

	o.Raw(sql, adspaceID).QueryRow(&adids)

	if len(adids) > 0 {
		for _, adid := range adids {
			for _, value := range *adInfos {
				if value.Adid == adid {
					value.Priority = value.Priority + 1
				}
			}
		}
	}

}

func getVersionId(version string, os string) int {
	var versionID = 0
	var platformID = 0

	if os == "android" {
		platformID = 200 // Android
	} else {
		platformID = 100 //iOS
	}

	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("device_os.os_id").
		From("device_os").
		Where("device_os.version = ?").
		And("device_os.platform_id = ?")

	sql := qb.String()
	o := orm.NewOrm()
	o.Using("default")
	o.Raw(sql, version, platformID).QueryRow(&versionID)
	return versionID
}
