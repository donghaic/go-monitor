package services

import (
	"encoding/json"
	"fmt"
	"gt-monitor/common/zap"
	"gt-monitor/config"
	"gt-monitor/models"
	"gt-monitor/services/dao"
	"gt-monitor/services/delayed"
	"gt-monitor/services/macro"
	"gt-monitor/utils"
	"time"
)

// 转化业务处理
type ConvHandler struct {
	config.KafakTopic
	dao          *dao.MongoDao
	syncer       *AdvertiserSyncer
	queue        *KafkaQueue
	delayedTaskQ *delayed.DelayedTaskQueue
	macro        *macro.SpecialConvURL
}

// New ...
func NewConvHandler(kafkaConf config.KafakTopic,
	dao *dao.MongoDao, adverSyncer *AdvertiserSyncer, kqueue *KafkaQueue, delayedTaskQ *delayed.DelayedTaskQueue) *ConvHandler {

	return &ConvHandler{
		kafkaConf,
		dao,
		adverSyncer,
		kqueue,
		delayedTaskQ,
		&macro.SpecialConvURL{},
	}
}

// 核心业务处理逻辑
// 0. 写本地请求日志
// 1. 查询点击日志，如果没查到点击，延迟处理转化
// 2. 广告主地址做宏地址处理
// 3. 同步广告主地址
// 4. 日志入kafka
func (h *ConvHandler) Handle(params *models.ConvReqParams) *models.Response {
	go logConvEvent(params)

	clkEvent, err := h.dao.FindClickByImpId(params.ImpId, params.ClkTs)
	if err != nil {
		// 延迟处理转化
		_ = h.delayedTaskQ.Put(delayed.Task{Cnt: 0, Params: params})
		return &models.Response{Code: 200, Content: "OK"}
	}

	go h.DoHandle(params, clkEvent)

	return &models.Response{Code: 200, Content: "OK"}

}

func (h *ConvHandler) DoHandle(params *models.ConvReqParams, clkEvent *models.ClickEntity) {
	oldConvEntity, _ := h.dao.FindConvByImpId(params.ImpId)
	if oldConvEntity != nil && utils.IsNotEmpty(oldConvEntity.ImpId) {
		params.EventType = "duplicated"
		go h.saveDBAndSendKafka(params, clkEvent, &models.ChannelCb{})
	}
	//回调处理
	cblink := clkEvent.CallBack
	if utils.IsEmpty(cblink) {
		cblink = h.macro.Special(params, clkEvent)
	}
	if utils.IsEmpty(cblink) {
		zap.Get().Info(" There is no callback link for this event. impId=", params.ImpId)
		go h.saveDBAndSendKafka(params, clkEvent, &models.ChannelCb{})
	} else {
		cbLink := macro.BaiduExpec(cblink, clkEvent.S2)
		cbUrl := macro.ReplacedClickTLMacroAndFunc(cbLink, params, *clkEvent)
		hstart := time.Now().Unix()
		resCbCont, resCbCode, err := h.syncer.SyncConv(cbUrl)

		errCont := ""
		if nil != err {
			errCont = err.Error()
			msg := fmt.Sprintf("sync conv callback url=%s time=%d error,", cbUrl, time.Now().Unix()-hstart)
			zap.Get().Info(msg, err)
		}

		go h.saveDBAndSendKafka(params, clkEvent, &models.ChannelCb{
			CbUrl:        cbUrl,
			CbResCode:    resCbCode,
			CbResContent: resCbCont,
			CbErrInfo:    errCont,
		})
	}
}

func (h *ConvHandler) saveDBAndSendKafka(reqParams *models.ConvReqParams, clickEntity *models.ClickEntity, cbinfo *models.ChannelCb) {
	go h.saveConvToMongoLogDB(reqParams, clickEntity, cbinfo)
	go h.sendToKafka(reqParams, clickEntity, cbinfo)
	go func() {
		err := h.dao.UpsertReport(clickEntity, reqParams)
		if err != nil {
			dd, _ := json.Marshal(reqParams)
			zap.Get().Error("upsert  report error", string(dd), err)
		}
	}()
}

//
func (h *ConvHandler) sendToKafka(params *models.ConvReqParams, clklog *models.ClickEntity, cbinfo *models.ChannelCb) {
	h.queue.sendConvEvent(h.KafakTopic.Conv, &models.ConvEvent{
		ImpId:      params.ImpId,
		MediaImpId: clklog.MediaImpId,

		AdvId:      clklog.AdvId,
		ProductId:  clklog.ProductId,
		CampaignId: clklog.CampaignId,
		AdGroupId:  clklog.AdGroupId,
		CreativeId: clklog.CreativeId,

		MediaId: clklog.MediaId,

		EventName: params.EventName,
		EventType: params.EventType,
		EventTs:   params.EventTs,
		EventDay:  params.EventDay,
		EventHour: params.EventHour,
		EventMin:  params.EventMin,
		EventSec:  params.EventSec,

		ReqUrl: params.ReqUrl,
		ReqIp:  params.ReqIp,
		ReqUa:  params.ReqUa,

		CbInfo: *cbinfo,
		//ClickInfo: *clklog,
	})
}

//
func (h *ConvHandler) saveConvToMongoLogDB(params *models.ConvReqParams, clickEntity *models.ClickEntity, cbinfo *models.ChannelCb) {
	error := h.dao.SaveConv(&models.ConvEntity{
		ImpId:      params.ImpId,
		MediaImpId: clickEntity.MediaImpId,

		AdvId:      clickEntity.AdvId,
		ProductId:  clickEntity.ProductId,
		CampaignId: clickEntity.CampaignId,
		AdGroupId:  clickEntity.AdGroupId,
		CreativeId: clickEntity.CreativeId,

		MediaId: clickEntity.MediaId,

		EventName: params.EventName,
		EventType: params.EventType,
		EventTs:   params.EventTs,
		EventDay:  params.EventDay,
		EventHour: params.EventHour,
		EventMin:  params.EventMin,
		EventSec:  params.EventSec,

		ReqUrl: params.ReqUrl,
		ReqIp:  params.ReqIp,
		ReqUa:  params.ReqUa,

		CbInfo:    *cbinfo,
		ClickInfo: *clickEntity,
	})
	if error != nil {
		zap.Get().Error("save conv to log db error")
	}
}

func logConvEvent(event *models.ConvReqParams) {
	data, _ := json.Marshal(event)
	line := fmt.Sprintf("conversion %s", string(data))
	zap.GetEvent().Info(line)
}
