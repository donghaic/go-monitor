package services

import (
	"gt-monitor/common/zap"
	"gt-monitor/models"
	"gt-monitor/utils"
)

type AdvertiserSyncer struct {
	httpCli *utils.HttpConPool
}

func NewAdverSyncer(httpCli *utils.HttpConPool) *AdvertiserSyncer {
	return &AdvertiserSyncer{
		httpCli,
	}
}

// 过滤大点击
// 同步广告主地址
func (s *AdvertiserSyncer) SyncImpClick(creative *models.CreativeInfo, syncAdver *models.SyncAdverInfo, params *models.ImpClkReqInfo) {
	if utils.IsNotEmpty(syncAdver.ServerSyncUrl) && isSyncable(params, creative) {
		ctx, code, err := s.httpCli.ReqGet(syncAdver.ServerSyncUrl)
		syncAdver.ServerSyncResCode = code
		syncAdver.ServerSyncResCtx = ctx
		if nil != err {
			syncAdver.ServerSyncResErr = err.Error()
		}
		// 返回内部过多：例如返回apk包
		if 512 < len(syncAdver.ServerSyncResCtx) {
			syncAdver.ServerSyncResCtx = string([]rune(syncAdver.ServerSyncResCtx)[0:512])
		}
	}
}

func (s *AdvertiserSyncer) SyncConv(url string) (string, int, error) {
	zap.Get().Info("sync conv url=", url)
	return s.httpCli.ReqGet(url)
}

func isSyncable(params *models.ImpClkReqInfo, creative *models.CreativeInfo) bool {
	// 屏蔽大点击
	if utils.IsNotEmpty(params.ActionCode) && 0 < len(creative.WaxExclusionCodes) {
		for _, wecode := range creative.WaxExclusionCodes {
			if utils.Equal(params.ActionCode, wecode) {
				return false
			}
		}
	}
	return true
}
