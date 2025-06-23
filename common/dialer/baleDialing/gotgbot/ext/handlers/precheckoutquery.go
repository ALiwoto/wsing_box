package handlers

import (
	"fmt"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

type PreCheckoutQuery struct {
	Filter   filters.PreCheckoutQuery
	Response Response
}

func NewPreCheckoutQuery(f filters.PreCheckoutQuery, r Response) PreCheckoutQuery {
	return PreCheckoutQuery{
		Filter:   f,
		Response: r,
	}
}

func (r PreCheckoutQuery) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.PreCheckoutQuery == nil {
		return false
	}
	return r.Filter == nil || r.Filter(ctx.PreCheckoutQuery)
}

func (r PreCheckoutQuery) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return r.Response(b, ctx)
}

func (r PreCheckoutQuery) Name() string {
	return fmt.Sprintf("precheckoutquery_%p", r.Response)
}
