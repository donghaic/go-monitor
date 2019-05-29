package macro

import (
	"gt-monitor/common"
	"gt-monitor/common/zap"
	"gt-monitor/models"
	"gt-monitor/utils"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func BuildImpClkAdverSyncInfo(eventType common.EventType, creative *models.CreativeInfo, link *models.LinkSwitch,
	synAdverInfo *models.SyncAdverInfo, params *models.ImpClkReqInfo) {
	adverLink := ""
	if eventType == common.Imp {
		adverLink = creative.AdvImpLink
	} else if eventType == common.Click {
		adverLink = creative.AdvClkLink
	}
	synAdverInfo.ServerSyncSourceUrl = adverLink

	//judge switch
	if nil != link {
		if eventType == common.Imp {
			adverLink = link.AfLinkImp
		} else if eventType == common.Click {
			adverLink = link.AfLinkClk
		}
		synAdverInfo.ServerSyncSourceUrl = adverLink
	}

	if utils.IsNotEmpty(adverLink) {
		adverLink = replacedClickTLMacroAndFunc(adverLink, params)
	}

	if utils.Equal(common.Direct_1_302, params.Direct) {
		synAdverInfo.RespnseCode = 302
		synAdverInfo.ResponseCtx = adverLink
	}

	if utils.Equal(common.Direct_2_S2S, params.Direct) {
		synAdverInfo.RespnseCode = 200
		synAdverInfo.ResponseCtx = "ok"
		synAdverInfo.ServerSyncUrl = adverLink
	}

	if utils.Equal(common.Direct_3_302_S2S, params.Direct) {
		if eventType == common.Imp {
			synAdverInfo.RespnseCode = 200
			synAdverInfo.ResponseCtx = "ok"
		} else {
			synAdverInfo.RespnseCode = 302
			synAdverInfo.ResponseCtx = creative.LandingPage
			synAdverInfo.ServerSyncUrl = adverLink
		}
	}
}

func replacedClickTLMacroAndFunc(urlin string, params *models.ImpClkReqInfo) string {
	relink := urlin
	//本系统生成的点击ID
	relink = strings.Replace(relink, "__CLICKID__", params.ImpId, -1)
	//__CTS__ 毫秒
	relink = strings.Replace(relink, "__CTS__", strconv.FormatInt(params.ReqTimeStamp/1000000, 10), -1)
	//__CTSYMDHMS__ 年月日
	relink = strings.Replace(relink, "__CTSYMDHMS__", time.Now().Format("2006-01-02 15:04:05"), -1)
	//__CTST__ 秒
	relink = strings.Replace(relink, "__CTST__", strconv.FormatInt(params.ReqTimeStamp/1000000000, 10), -1)
	//__IP__
	if 0 == len(params.Ip) {
		relink = strings.Replace(relink, "__IP__", params.ReqIp, -1)
	} else {
		relink = strings.Replace(relink, "__IP__", params.Ip, -1)
	}
	//__CHANNEL__
	relink = strings.Replace(relink, "__CHANNEL__", params.MediaId, -1)
	//__S1__
	relink = strings.Replace(relink, "__S1__", params.S1, -1)
	//__S2__
	relink = strings.Replace(relink, "__S2__", params.S2, -1)
	//__S3__
	relink = strings.Replace(relink, "__S3__", params.S3, -1)
	//__S4__
	relink = strings.Replace(relink, "__S4__", params.S4, -1)
	//__S5__
	relink = strings.Replace(relink, "__S5__", params.S5, -1)
	//__UA__
	relink = strings.Replace(relink, "__UA__", url.QueryEscape(params.Ua), -1)
	//__IDFA__
	relink = strings.Replace(relink, "__IDFA__", params.Idfa, -1)
	// //__IDFAMD5__
	// relink = strings.Replace(relink, "__IDFAMD5__", params.Idfa, -1)
	//__ANDROIDID__
	relink = strings.Replace(relink, "__ANDROIDID__", params.AndroidId, -1)

	//__IMEI__
	relink = strings.Replace(relink, "__IMEI__", params.Imei, -1)

	//__IMEIMD5__
	relink = strings.Replace(relink, "__IMEIMD5__", params.ImeiMd5, -1)

	//__OS__
	relink = strings.Replace(relink, "__OS__", params.Os, -1)

	return replacedClickTLFunction(relink)
	// return relink
}

func replacedClickTLFunction(urlin string) string {
	relink := urlin
	if strings.Contains(urlin, "=heyman,") {
		u, err := url.Parse(urlin)
		if nil == err {
			uvalues, err := url.ParseQuery(u.RawQuery)
			if nil == err {
				for key, value := range uvalues {
					if len(value) > 0 {
						if strings.Contains(value[0], "heyman,") {
							newv := functionRepPro(key, value[0])
							uvalues.Del(key)
							uvalues.Set(key, newv)
						}
					}
				}
			}
			chaneg := uvalues.Encode()
			u.RawQuery = chaneg

			relink = u.String()
		}
	}
	return relink
}

func functionRepPro(key, value string) string {
	vs := strings.Split(value, "heyman,")
	if len(vs) > 1 {
		strs := vs[1]

		for {
			if !strings.Contains(strs, "[") {
				break
			}

			pstart := strings.LastIndex(strs, "[")
			pend := strings.Index(strs, "]")

			params := strs[pstart+1 : pend]

			fstart := strings.LastIndex(strs[:pstart], "[")
			fname := strs[:pstart][fstart+1:]

			execvalue := ""
			switch fname {
			case "MD5":
				execvalue = utils.Md5(params)
			case "SHA1":
				execvalue = utils.Sha1(params)
			case "UPSTR":
				execvalue = strings.ToUpper(params)
			case "DOWNSTR":
				execvalue = strings.ToLower(params)
			case "RSASIGN":
				execvalue = utils.RsaSign(params)
			}

			strstmp := strs[:pstart][:strings.LastIndex(strs[:pstart], "[")+1]
			strstmp = strings.Join([]string{strstmp, execvalue}, "")

			strstmp = strings.Join([]string{strstmp, strs[pend+1:]}, "")

			strs = strstmp
		}

		return strs
	}

	return ""
}

// ------------------------
//
// conversion
func BaiduExpec(linkin, akey string) string {
	if (strings.Contains(linkin, "{{ATYPE}}")) && (strings.Contains(linkin, "{{AVALUE}}")) {
		linkin = strings.Replace(linkin, "{{ATYPE}}", "activate", -1)
		linkin = strings.Replace(linkin, "{{AVALUE}}", "0", -1)

		sign := utils.Md5(strings.Join([]string{linkin, akey}, ""))

		linkin = strings.Join([]string{linkin, "&sign=", strings.ToLower(sign)}, "")

		// linkin = strings.Replace(linkin, "==", `%3D%3D`, -1)

		return linkin
	} else {
		return linkin
	}
}

func ReplacedClickTLMacroAndFunc(url string, reqParams *models.ConvReqParams, clklog models.ClickEntity) string {
	relink := url

	relink = strings.Replace(relink, "__CCID__", clklog.MediaImpId, -1)
	relink = strings.Replace(relink, "__IDFA__", clklog.Idfa, -1)
	relink = strings.Replace(relink, "__EVENTUNIX__", strconv.FormatInt(reqParams.EventTs/1000000000, 10), -1)

	// app类型；取值为 android或ios（联盟Android为unionandroid）
	if strings.Contains(clklog.ReqUrl, "app_type=ios") {
		relink = strings.Replace(relink, "__MUID__", clklog.IdfaMd5, -1)
		relink = strings.Replace(relink, "__APPTYPE__", "ios", -1)
	} else {
		relink = strings.Replace(relink, "__MUID__", clklog.ImeiMd5, -1)
		if strings.Contains(clklog.ReqUrl, "app_type=unionandroid") {
			relink = strings.Replace(relink, "__APPTYPE__", "unionandroid", -1)
		} else {
			relink = strings.Replace(relink, "__APPTYPE__", "android", -1)
		}
	}

	return relink
}

const (
	NetEaseMusicCb = "https://ad-effect.music.163.com/ad/action?"
	// NetEaseMusicCb     = "https://qa-qtm.igame.163.com/ad/action?"
	NetEaseMusicSKey   = "YU3B0LK80kMcG0c319Qu6RA4464504181v8Vs64Yk2sW7C5561275U1188XB4R91"
	NetEaseMusicAppKey = "hds8b8aWN4QJ6426C1J5LusE2285D4M373R3lq3YM1WW2792cyk0yKd52y42B1Jx"
	NetEaseMusicSource = "16"
)

type SpecialConvURL struct{}

func (s *SpecialConvURL) Special(reqParams *models.ConvReqParams, clklog *models.ClickEntity) string {
	speial := strings.Contains(clklog.ReqUrl, "maxwell=electromagnetism")

	if speial {
		zap.Get().Info(clklog.ImpId, " Need special process, url=", clklog.ReqUrl)
		return s.NetEaseMusic(reqParams, clklog)
	}

	return ""
}

func (s *SpecialConvURL) NetEaseMusic(reqParams *models.ConvReqParams, clklog *models.ClickEntity) string {
	neteasekeys := []string{"appkey", "clickid", "event", "platType", "source", "timestamp"}
	neteasevalues := []string{NetEaseMusicAppKey, clklog.MediaImpId, "0", "0", NetEaseMusicSource, strconv.FormatInt(time.Now().Unix(), 10)}

	link := []string{}
	signsrc := []string{NetEaseMusicSKey}
	for i, k := range neteasekeys {
		signsrc = append(signsrc, k)
		signsrc = append(signsrc, neteasevalues[i])
		signsrc = append(signsrc, NetEaseMusicSKey)

		link = append(link, strings.Join([]string{k, neteasevalues[i]}, "="))
	}

	sign := utils.Md5(strings.Join(signsrc, ""))

	return strings.Join([]string{NetEaseMusicCb, strings.Join(link, "&"), "&sign=", strings.ToUpper(sign)}, "")
}
