package api

import (
	"gt-monitor/common"
	"gt-monitor/models"
	"gt-monitor/utils"
	"github.com/valyala/fasthttp"
	"strconv"
)

func parse(ctx *fasthttp.RequestCtx) (*models.ImpClkReqInfo, error) {
	params := utils.ParseHttpReq(ctx)
	if 0 == len(params.CreativeId) {
		return params, common.Bad_Req_Creative_Id
	}

	_, err := strconv.Atoi(params.CreativeId)
	if nil != err {
		return params, common.Bad_Req_Creative_Id
	}

	if utils.IsNotEmpty(params.Direct) {
		directNum, err := strconv.Atoi(params.Direct)
		if nil != err {
			return params, common.Bad_Req_Direct
		}

		if 3 < directNum {
			return params, common.Bad_Req_Direct
		}
	}
	return params, nil
}
