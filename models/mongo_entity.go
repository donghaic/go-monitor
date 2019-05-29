package models

// 点击
type ClickEntity struct {
	RequestId string `bson:"request_id"`
	TimeStamp string `bson:"timestamp"`

	ReqTimeStamp int64 `bson:"req_timestamp"`

	ImpId      string `bson:"impid"`
	AdxImpId   string `bson:"adx_impid"`
	MediaImpId string `bson:"media_impid"`

	AdvId      string `bson:"adver_id"`
	ProductId  string `bson:"product_id"`
	CampaignId string `bson:"campaign_id"`
	AdGroupId  string `bson:"adgroup_id"`
	CreativeId string `bson:"creative_id"`

	MediaId string `bson:"media_id"`
	AdPosId string `bson:"adpos_id"`

	Price string `bson:"price"`

	AdvLinkTag   string `bson:"advlink_tag"`
	AdvLinkTagId string `bson:"advlink_tag_id"`

	DevType  string `bson:"device_type"`
	Os       string `bson:"os"`
	ConnType string `bson:"connection_type"`
	Ip       string `bson:"ip"`
	ReqIp    string `bson:"req_ip"`
	Ua       string `bson:"ua"`
	ReqUa    string `bson:"req_ua"`
	ReqUrl   string `bson:"req_url"`

	Idfa     string `bson:"idfa"`
	IdfaMd5  string `bson:"idfa_md5"`
	IdfaSha1 string `bson:"idfa_sha1"`

	AndroidId     string `bson:"android_id"`
	AndroidIdMd5  string `bson:"android_id_md5"`
	AndroidIdSha1 string `bson:"impandroid_id_sha1"`

	Imei     string `bson:"imei"`
	ImeiMd5  string `bson:"imei_md5"`
	ImeiSha1 string `bson:"imei_sha1"`

	DeviceId string `bson:"device_id"`

	Mac     string `bson:"mac"`
	MacMd5  string `bson:"mac_md5"`
	MacSha1 string `bson:"mac_sha1"`

	CurAdv string `bson:"cur_adv"`
	CurAdx string `bson:"cur_adx"`

	Make         string `bson:"make"`
	Model        string `bson:"model"`
	Bundle       string `bson:"bundle"`
	AppName      string `bson:"app_name"`
	MaterialType string `bson:"material_type"`
	TaId         string `bson:"taid"`
	TagId        string `bson:"tagid"`
	CrowdId      string `bson:"crowdid"`

	CallBack string `bson:"callback"`
	Direct   string `bson:"direct"`

	S1 string `bson:"s1"`
	S2 string `bson:"s2"`
	S3 string `bson:"s3"`
	S4 string `bson:"s4"`
	S5 string `bson:"s5"`

	UserId       string `bson:"user_id"`
	UserBudgetId string `bson:"user_budgetid"`

	SyncToAdv SyncAdverInfo `bson:"sync_adv"`
	LinkSw    LinkSwitch    `bson:"link_switch"`
}

//
type ConvEntity struct {
	ImpId      string `bson:"impid"`
	MediaImpId string `bson:"media_impid"`

	AdvId      string `bson:"adver_id"`
	ProductId  string `bson:"product_id"`
	CampaignId string `bson:"campaign_id"`
	AdGroupId  string `bson:"adgroup_id"`
	CreativeId string `bson:"creative_id"`

	MediaId string `bson:"media_id"`

	EventName string `bson:"event_name"`
	EventType string `bson:"event_type"`
	EventTs   int64  `bson:"event_ts"`
	EventDay  string `bson:"event_day"`
	EventHour int    `bson:"event_hour"`
	EventMin  int    `bson:"event_min"`
	EventSec  int    `bson:"event_sec"`

	ReqUrl string `bson:"req_url"`
	ReqIp  string `bson:"req_ip"`
	ReqUa  string `bson:"req_ua"`

	CbInfo ChannelCb `bson:"callback_info"` //渠道回调信息

	ClickInfo ClickEntity `bson:"click_info"` //点击日志
}

//渠道回调信息
type ChannelCb struct {
	CbUrl        string `bson:"cburl"`
	CbResCode    int    `bson:"cb_rescode"`
	CbResContent string `bson:"cb_rescontent"`
	CbErrInfo    string `bson:"cb_errinfo"`
}
