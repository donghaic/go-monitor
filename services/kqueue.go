package services

import (
	"gt-monitor/common"
	"gt-monitor/common/kafka"
	"gt-monitor/common/zap"
	"gt-monitor/models"
	"time"
)

type KafkaQueue struct {
	producer kafka.Producer
}

func NewKqueue(producer kafka.Producer) *KafkaQueue {
	return &KafkaQueue{producer: producer}
}

func (q *KafkaQueue) SendImpOrClickEvent(topic string, eventType common.EventType, syncAdv *models.SyncAdverInfo,
	link *models.LinkSwitch, params *models.ImpClkReqInfo) {

	event := models.ImpClkEvent{
		RequestId:    params.RequestId,
		TimeStamp:    params.TimeStamp,
		ReqTimeStamp: params.ReqTimeStamp,
		ReqTimeStr:   time.Now().Format("2006-01-02 15:04:05"),

		ImpId:      params.ImpId,
		AdxImpId:   params.AdxImpId,
		MediaImpId: params.MediaImpId,

		AdvId:      params.AdvId,
		ProductId:  params.ProductId,
		CampaignId: params.CampaignId,
		AdGroupId:  params.AdGroupId,
		CreativeId: params.CreativeId,

		MediaId: params.MediaId,
		AdPosId: params.AdPosId,

		Price: params.Price,

		AdvLinkTag:   params.AdvLinkTag,
		AdvLinkTagId: params.AdvLinkTagId,

		DevType:  params.DevType,
		Os:       params.Os,
		ConnType: params.ConnType,
		Ip:       params.Ip,
		ReqIp:    params.ReqIp,
		Ua:       params.Ua,
		ReqUa:    params.ReqUa,
		ReqUrl:   params.ReqUrl,

		Idfa:     params.Idfa,
		IdfaMd5:  params.IdfaMd5,
		IdfaSha1: params.IdfaSha1,

		AndroidId:     params.AndroidId,
		AndroidIdMd5:  params.AndroidIdMd5,
		AndroidIdSha1: params.AndroidIdSha1,

		Imei:     params.Imei,
		ImeiMd5:  params.ImeiMd5,
		ImeiSha1: params.ImeiSha1,

		DeviceId: params.DeviceId,

		Mac:     params.Mac,
		MacMd5:  params.MacMd5,
		MacSha1: params.MacSha1,

		CurAdv: params.CurAdv,
		CurAdx: params.CurAdx,

		Make:         params.Make,
		Model:        params.Model,
		Bundle:       params.Bundle,
		AppName:      params.AppName,
		MaterialType: params.MaterialType,
		TaId:         params.TaId,
		TagId:        params.TagId,
		CrowdId:      params.CrowdId,

		ActionCode: params.ActionCode,

		CallBack: params.CallBack,
		Direct:   params.Direct,

		S1: params.S1,
		S2: params.S2,
		S3: params.S3,
		S4: params.S4,
		S5: params.S5,

		UserId:       params.UserId,
		UserBudgetId: params.UserBudgetId,

		SyncToAdv: *syncAdv,
	}
	if nil != link {
		event.LinkSwInfo = *link
		event.AdvLinkTag = "switched"
		event.AdvLinkTagId = link.AfLinkId
	}

	if 0 == len(event.AdvLinkTagId) {
		event.AdvLinkTagId = "0"
	}

	if eventType == common.Imp {
		event.LogType = "impression"
	} else if eventType == common.Click {
		event.LogType = "click"
	}

	err := q.producer.SendKeyedMessage(topic, event.ReqIp, event)
	if err != nil {
		zap.Get().Error("send conv msg error")
	}
}

func (q *KafkaQueue) sendConvEvent(topic string, event *models.ConvEvent) {
	err := q.producer.SendKeyedMessage(topic, event.ReqIp, event)
	if err != nil {
		zap.Get().Error("send conv msg error")
	}
}
