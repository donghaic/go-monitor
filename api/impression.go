package api

import (
	"gt-monitor/common"
	"gt-monitor/common/zap"
	"gt-monitor/services"
	"github.com/valyala/fasthttp"
)

type Imp struct {
	h *services.ImpClkHandler
}

func NewImp(handler *services.ImpClkHandler) *Imp {
	return &Imp{
		handler,
	}
}

func (i *Imp) Handle(ctx *fasthttp.RequestCtx) {
	params, err := parse(ctx)
	if err != nil {
		zap.Get().Error("request url=", string(ctx.RequestURI()))
		_, _ = ctx.WriteString(err.Error())
		return
	}

	response := i.h.Handle(common.Imp, params)
	if 302 == response.Code {
		ctx.Redirect(response.Content, 302)
	} else {
		_, _ = ctx.WriteString(response.Content)
	}
}
