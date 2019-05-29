package api

import (
	"gt-monitor/common"
	"gt-monitor/common/zap"
	"gt-monitor/services"
	"github.com/valyala/fasthttp"
)

type Click struct {
	handler *services.ImpClkHandler
}

func NewClick(handler *services.ImpClkHandler) *Click {
	return &Click{
		handler,
	}
}

func (c *Click) Handle(ctx *fasthttp.RequestCtx) {
	params, err := parse(ctx)
	if err != nil {
		zap.Get().Error("request url=", string(ctx.RequestURI()))
		_, _ = ctx.WriteString(err.Error())
		return
	}

	response := c.handler.Handle(common.Click, params)
	if 302 == response.Code {
		ctx.Redirect(response.Content, 302)
	} else {
		_, _ = ctx.WriteString(response.Content)
	}
}
