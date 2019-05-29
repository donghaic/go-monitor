package utils

import (
	"gt-monitor/models"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"time"
)

func ParseHttpReq(ctx *fasthttp.RequestCtx) *models.ImpClkReqInfo {
	args := ctx.URI().QueryArgs()

	reqParams := models.ImpClkReqInfo{
		RequestId:    string(args.Peek("request_id")),
		TimeStamp:    string(args.Peek("timestamp")),
		ReqTimeStamp: time.Now().UnixNano(),
		AdxImpId:     string(args.Peek("adx_impid")),
		ImpId:        string(args.Peek("impid")),
		MediaImpId:   string(args.Peek("cid")),

		AdvId:      string(args.Peek("adver_id")),
		ProductId:  string(args.Peek("product_id")),
		CampaignId: string(args.Peek("campaign_id")),
		AdGroupId:  string(args.Peek("adgroup_id")),
		CreativeId: string(args.Peek("creative_id")),

		MediaId: string(args.Peek("media_id")),
		AdPosId: string(args.Peek("adpos_id")),

		Price: string(args.Peek("price")),

		DevType:  string(args.Peek("devicetype")),
		Os:       string(args.Peek("os")),
		ConnType: string(args.Peek("connectiontype")),
		Ip:       string(args.Peek("ip")),
		Ua:       string(args.Peek("ua")),

		Idfa:          string(args.Peek("idfa")),
		IdfaMd5:       string(args.Peek("idfa_md5")),
		IdfaSha1:      string(args.Peek("idfa_sha1")),
		AndroidId:     string(args.Peek("android_id")),
		AndroidIdMd5:  string(args.Peek("android_id_md5")),
		AndroidIdSha1: string(args.Peek("android_id_sha1")),
		Imei:          string(args.Peek("imei")),
		ImeiMd5:       string(args.Peek("imei_md5")),
		ImeiSha1:      string(args.Peek("imei_sha1")),
		DeviceId:      string(args.Peek("deviceID")),
		Mac:           string(args.Peek("mac")),
		MacMd5:        string(args.Peek("mac_md5")),
		MacSha1:       string(args.Peek("mac_sha1")),

		CurAdv: string(args.Peek("cur_adv")),
		CurAdx: string(args.Peek("cur_adx")),

		ActionCode: string(args.Peek("action_code")),

		Make:         string(args.Peek("make")),
		Model:        string(args.Peek("model")),
		AppName:      string(args.Peek("app_name")),
		MaterialType: string(args.Peek("materialType")),
		TaId:         string(args.Peek("taid")),
		TagId:        string(args.Peek("tagId")),
		Bundle:       string(args.Peek("bundle")),
		CrowdId:      string(args.Peek("crowdId")),
		CallBack:     string(args.Peek("callback")),
		Direct:       string(args.Peek("direct")),

		S1: string(args.Peek("s1")),
		S2: string(args.Peek("s2")),
		S3: string(args.Peek("s3")),
		S4: string(args.Peek("s4")),
		S5: string(args.Peek("s5")),

		UserId:       string(args.Peek("user_id")),
		UserBudgetId: string(args.Peek("user_budgetid")),
	}

	reqParams.ReqIp = FastGetIp(ctx)
	reqParams.ReqUa = string(ctx.UserAgent())
	reqParams.ReqUrl = string(ctx.RequestURI())

	// app类型；取值为 android或ios（联盟Android为unionandroid）
	if 0 == strings.Compare("zuc", string(args.Peek("zuc"))) {
		if strings.Contains(reqParams.ReqUrl, "app_type=ios") {
			reqParams.MediaImpId = string(args.Peek("click_id"))
			reqParams.Idfa = string(args.Peek("muid"))
			reqParams.IdfaMd5 = string(args.Peek("muid"))
		} else {
			reqParams.MediaImpId = string(args.Peek("click_id"))
			reqParams.Imei = string(args.Peek("muid"))
			reqParams.ImeiMd5 = string(args.Peek("muid"))
		}
	}

	if "102" == string(args.Peek("media_id")) {
		priceStr := strings.Split(reqParams.ReqUrl, "price=")
		if 2 == len(priceStr) {
			priceStrTmps := strings.Split(priceStr[1], "&")

			// fmt.Println("MoMo media 102.", priceStrTmps[0], reqParams.Price)

			reqParams.Price = priceStrTmps[0]
		}
	}

	_, err := strconv.ParseInt(reqParams.TimeStamp, 10, 64)
	if nil != err {
		reqParams.TimeStamp = strconv.FormatInt(reqParams.ReqTimeStamp, 10)
	}

	// https://noticecpc.ad-mex.com/clk?mmid=1245&cid=__CID__&idfa=__IDFA__&ip=__IP__&ua=__UA__&callback=__CALLBACK_URL__
	// 47.94.112.89
	// https://monitor.clk?mmid=1245&cid=123456789&idfa=__IDFA__&ip=__IP__&ua=__UA__&callback=__CALLBACK_URL__
	// adver_id=3&product_id=2&campaign_id=3&adgroup_id=766&creative_id=3019&media_id=2
	if 0 != len(string(args.Peek("mmid"))) && 0 == strings.Compare("1245", string(args.Peek("mmid"))) {
		reqParams.AdvId = "3"
		reqParams.ProductId = "2"
		reqParams.CampaignId = "3"
		reqParams.AdGroupId = "766"
		reqParams.CreativeId = "3019"
		reqParams.MediaId = "2"
	}

	return &reqParams
}
