package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Ad struct {
	Adid                    int64 `orm:"pk"`
	AdName                  string
	StartTime               time.Time
	EndTime                 time.Time
	LinkUrl                 string
	State                   string
	AdMaxImpression         string
	CampaignId              int64
	AccountId               int64
	FormatId                int64
	CurrencyId              int64
	PreviewUrl              string
	ThirdPartyImpressionUrl string
	ActionText              string
	BidType                 string
	BidAmount               float64
	DailyMaxImpression      int64
	DailyMaxClick           int64
	ThirdPartyClickTracker  int64
	PlatformId              int64
	OsMinVersionId          int64
	OsMaxVersionId          int64
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(Ad))
}
