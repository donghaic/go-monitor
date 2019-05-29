package api

import (
	"gt-monitor/models"
	"gt-monitor/services"
	"gt-monitor/utils"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"time"
)

type Conversion struct {
	handler *services.ConvHandler
}

func NewConversion(handler *services.ConvHandler) *Conversion {
	return &Conversion{handler,}
}

func (c *Conversion) Handle(ctx *fasthttp.RequestCtx) {
	args := ctx.URI().QueryArgs()
	reqParams := models.ConvReqParams{
		ImpId:     string(args.Peek("clickid")),
		ClkTs:     string(args.Peek("t")),
		EventName: string(args.Peek("event_name")),
	}

	if 0 == len(reqParams.ImpId) || 19 > len(reqParams.ImpId) {
		ctx.Response.Header.Set("Content-Type", "text/plain;charset=utf-8")
		ctx.Response.Header.SetStatusCode(http.StatusNotFound)
		_, _ = ctx.Write([]byte("not found"))
		return
	}

	reqParams.ClkTs = reqParams.ImpId[:19]
	_, err := strconv.ParseInt(reqParams.ClkTs, 10, 64)
	if nil != err {
		ctx.Response.Header.Set("Content-Type", "text/plain;charset=utf-8")
		ctx.Response.Header.SetStatusCode(http.StatusNotFound)
		_, _ = ctx.Write([]byte("not found"))
		return
	}

	reqParams.EventType = "normal"
	reqParams.EventTs = time.Now().UnixNano()
	reqParams.EventDay = time.Now().Format("2006-01-02")
	reqParams.EventHour = time.Now().Hour()
	reqParams.EventMin = time.Now().Minute()
	reqParams.EventSec = time.Now().Second()

	reqParams.ReqUrl = string(ctx.RequestURI())
	reqParams.ReqIp = utils.FastGetIp(ctx)
	reqParams.ReqUa = string(ctx.UserAgent())

	res := c.handler.Handle(&reqParams)

	//返回结果
	if 302 == res.Code {
		ctx.Redirect(res.Content, 302)
	} else {
		_, _ = ctx.WriteString("recevied")
	}
}
