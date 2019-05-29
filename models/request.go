package models

/**
 * 展示点击 HTTP 请求参数信息
 */
type ImpClkReqInfo struct {
	RequestId string
	TimeStamp string

	ReqTimeStamp int64

	AdxImpId   string
	ImpId      string
	MediaImpId string

	AdvId      string
	ProductId  string
	CampaignId string
	AdGroupId  string
	CreativeId string

	MediaId string
	AdPosId string

	AdvLinkTag   string
	AdvLinkTagId string

	Price string

	DevType  string
	Os       string
	ConnType string
	Ip       string
	ReqIp    string
	Ua       string
	ReqUa    string
	ReqUrl   string

	Idfa     string
	IdfaMd5  string
	IdfaSha1 string

	AndroidId     string
	AndroidIdMd5  string
	AndroidIdSha1 string

	Imei     string
	ImeiMd5  string
	ImeiSha1 string

	DeviceId string

	Mac     string
	MacMd5  string
	MacSha1 string

	CurAdv string
	CurAdx string

	Make         string
	Model        string
	Bundle       string
	AppName      string
	MaterialType string
	TaId         string
	TagId        string
	CrowdId      string

	ActionCode string

	CallBack string
	Direct   string

	S1 string
	S2 string
	S3 string
	S4 string
	S5 string

	UserId       string
	UserBudgetId string
}

/**
 * 转化 HTTP 请求参数信息
 */
type ConvReqParams struct {
	ImpId string
	ClkTs string

	EventName string
	EventType string //独立 重复
	EventTs   int64
	EventDay  string
	EventHour int
	EventMin  int
	EventSec  int

	ReqUrl string
	ReqIp  string
	ReqUa  string
}
