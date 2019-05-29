package common

import "errors"

type EventType string

const (
	Click EventType = "click"
	Imp   EventType = "impression"
	Conv  EventType = "conversion"
)

const (
	Direct_0_No_Found = "0" // 返回200，无此值，不同步广告主链接
	Direct_1_302      = "1" // 返回302，返回宏替换之后的广告主链接
	Direct_2_S2S      = "2" // 返回200，服务器同步宏替换之后的广告主链接
	Direct_3_302_S2S  = "3" // 返回302，302地址为广告主的产品页，服务器同步宏替换的广告主链接
)

var (
	Bad_Req_Creative_Id = errors.New("illegal param creative_id")
	Bad_Req_Direct      = errors.New("illegal param direct")
)
