package services

import (
	"encoding/json"
	"fmt"
	"gt-monitor/common"
	"gt-monitor/common/zap"
	"gt-monitor/config"
	"gt-monitor/models"
	"gt-monitor/services/cache"
	"gt-monitor/services/macro"
	"gt-monitor/utils"
	"strconv"
	"strings"
)

type ImpClkHandler struct {
	config.KafakTopic
	cache  *cache.EntityCacheService
	syncer *AdvertiserSyncer
	queue  *KafkaQueue
}

// New ...
func NewImpClkHandler(kafkaConf config.KafakTopic,
	cache *cache.EntityCacheService,
	adverSyncer *AdvertiserSyncer,
	kqueue *KafkaQueue) *ImpClkHandler {

	return &ImpClkHandler{
		kafkaConf,
		cache,
		adverSyncer,
		kqueue,
	}
}

// 核心业务处理逻辑
// 0. 写本地请求日志
// 1. 根据事件参数查询缓存数据.
// 2. 广告主地址做宏地址处理
// 3. 同步广告主地址
// 4. 日志入kafka
func (h *ImpClkHandler) Handle(eventType common.EventType, params *models.ImpClkReqInfo) models.Response {
	logImpClkEvent(eventType, params)

	advSync := models.SyncAdverInfo{
		Direct: params.Direct,
	}
	topicName := h.getTopic(eventType)
	if 0 == len(params.Direct) {
		go h.queue.SendImpOrClickEvent(topicName, eventType, &advSync, nil, params)

		return models.Response{
			Code:     200,
			Content:  "ok",
			LandPage: "ok",
		}
	}

	//获取广告创意
	creative, err := h.cache.GetCreativeById(params.CreativeId)
	if nil != err {
		zap.Get().Error("GetCreativeById error, creativeId=", params.CreativeId, ", ", err)
		advSync.IsCache = false
		advSync.IsCacheErr = "query cache error"
		go h.queue.SendImpOrClickEvent(topicName, eventType, &advSync, nil, params)

		return models.Response{
			Code:     200,
			Content:  "PARAMETER INVALID [creative_id]",
			LandPage: "PARAMETER INVALID [creative_id]",
		}
	}

	linkSwitch, err := h.cache.GetLinkById(creative.AdvLinkId)
	if err != nil {
		zap.Get().Error("GetLinkById error, linkId=", creative.AdvLinkId, ", ", err)
	}

	params.AdvLinkTag = creative.AdvLinkTag
	params.AdvLinkTagId = creative.AdvLinkId

	if utils.IsEmpty(params.UserId) {
		params.UserId = creative.UserId
	}
	if utils.IsEmpty(params.UserBudgetId) {
		params.UserBudgetId = creative.UserBudgetId
	}

	advSync.IsCache = true
	advSync.IsCacheErr = ""

	if utils.IsEmpty(params.ImpId) {
		params.ImpId = utils.GenUniqueId(creative.CreativeId, params.ReqTimeStamp, utils.Md5(creative.CreativeId))
	} else {
		params.ImpId = strings.Join([]string{strconv.FormatInt(params.ReqTimeStamp, 10), params.ImpId}, "")
	}

	// 宏地址处理，同步广告主地址，日志入kafka
	macro.BuildImpClkAdverSyncInfo(eventType, creative, linkSwitch, &advSync, params)
	go h.syncer.SyncImpClick(creative, &advSync, params)
	go h.queue.SendImpOrClickEvent(topicName, eventType, &advSync, linkSwitch, params)

	return models.Response{
		Code:     advSync.RespnseCode,
		Content:  advSync.ResponseCtx,
		LandPage: advSync.ResponseCtx,
	}
}

func (h *ImpClkHandler) getTopic(eventType common.EventType) string {
	if eventType == common.Imp {
		return h.KafakTopic.Imp
	} else if eventType == common.Click {
		return h.KafakTopic.Click
	}
	panic(fmt.Sprintf("illegel event type: %d", eventType))
}

func logImpClkEvent(eType common.EventType, event *models.ImpClkReqInfo) {
	go func() {
		var logType string
		if eType == common.Imp {
			logType = "imp"
		} else {
			logType = "click"
		}
		data, _ := json.Marshal(event)
		line := fmt.Sprintf("%s %s", logType, string(data))
		zap.GetEvent().Info(line)
	}()
}
