package models

// 展示点击kafka日志
type ImpClkEvent struct {
	LogType   string `json:"logtype"`
	RequestId string `json:"request_id"`
	TimeStamp string `json:"timestamp"`

	ReqTimeStamp int64  `json:"req_timestamp"`
	ReqTimeStr   string `json:"req_time_str"`

	ImpId      string `json:"impid"`
	AdxImpId   string `json:"adx_impid"`
	MediaImpId string `json:"media_impid"`

	AdvId      string `json:"adver_id"`
	ProductId  string `json:"product_id"`
	CampaignId string `json:"campaign_id"`
	AdGroupId  string `json:"adgroup_id"`
	CreativeId string `json:"creative_id"`

	MediaId string `json:"media_id"`
	AdPosId string `json:"adpos_id"`

	Price string `json:"price"`

	AdvLinkTag   string `json:"advlink_tag"`
	AdvLinkTagId string `json:"advlink_tag_id"`

	DevType  string `json:"device_type"`
	Os       string `json:"os"`
	ConnType string `json:"connection_type"`
	Ip       string `json:"ip"`
	ReqIp    string `json:"req_ip"`
	Ua       string `json:"ua"`
	ReqUa    string `json:"req_ua"`
	ReqUrl   string `json:"req_url"`

	Idfa     string `json:"idfa"`
	IdfaMd5  string `json:"idfa_md5"`
	IdfaSha1 string `json:"idfa_sha1"`

	AndroidId     string `json:"android_id"`
	AndroidIdMd5  string `json:"android_id_md5"`
	AndroidIdSha1 string `json:"impandroid_id_sha1"`

	Imei     string `json:"imei"`
	ImeiMd5  string `json:"imei_md5"`
	ImeiSha1 string `json:"imei_sha1"`

	DeviceId string `json:"device_id"`

	Mac     string `json:"mac"`
	MacMd5  string `json:"mac_md5"`
	MacSha1 string `json:"mac_sha1"`

	CurAdv string `json:"cur_adv"`
	CurAdx string `json:"cur_adx"`

	Make         string `json:"make"`
	Model        string `json:"model"`
	Bundle       string `json:"bundle"`
	AppName      string `json:"app_name"`
	MaterialType string `json:"material_type"`
	TaId         string `json:"taid"`
	TagId        string `json:"tagid"`
	CrowdId      string `json:"crowdid"`

	CallBack string `json:"callback"`
	Direct   string `json:"direct"`

	LinkSwInfo LinkSwitch `json:"link_switch"`

	SyncToAdv SyncAdverInfo `json:"sync_adv"`

	ActionCode string `json:"action"`

	S1 string `json:"s1"`
	S2 string `json:"s2"`
	S3 string `json:"s3"`
	S4 string `json:"s4"`
	S5 string `json:"s5"`

	UserId       string `json:"user_id"`
	UserBudgetId string `json:"user_budgetid"`
}

type ConvEvent struct {
	ImpId      string `json:"impid"`
	MediaImpId string `json:"media_impid"`

	AdvId      string `json:"adver_id"`
	ProductId  string `json:"product_id"`
	CampaignId string `json:"campaign_id"`
	AdGroupId  string `json:"adgroup_id"`
	CreativeId string `json:"creative_id"`

	MediaId string `json:"media_id"`

	EventName string `json:"event_name"`
	EventType string `json:"event_type"`
	EventTs   int64  `json:"event_ts"`
	EventDay  string `json:"event_day"`
	EventHour int    `json:"event_hour"`
	EventMin  int    `json:"event_min"`
	EventSec  int    `json:"event_sec"`

	ReqUrl string `json:"req_url"`
	ReqIp  string `json:"req_ip"`
	ReqUa  string `json:"req_ua"`

	CbInfo ChannelCb `json:"callback_info"` //渠道回调信息
}

// 同步广告主信息
type SyncAdverInfo struct {
	IsCache    bool   `json:"is_cache"`
	IsCacheErr string `json:"is_cache_err"`

	Direct      string `json:"direct"`
	RespnseCode int    `json:"reponse_code"`
	ResponseCtx string `json:"response_ctx"`

	ServerSyncSourceUrl string `json:"server_sync_sourceurl"`
	ServerSyncUrl       string `json:"server_sync_url"`
	ServerSyncResCode   int    `json:"server_sync_rescode"`
	ServerSyncResCtx    string `json:"server_sync_resctx"`
	ServerSyncResErr    string `json:"server_sync_reserr"`
}
