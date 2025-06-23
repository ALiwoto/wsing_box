package handlers

import (
	"fmt"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

type InlineQuery struct {
	Filter   filters.InlineQuery
	Response Response
}

func NewInlineQuery(f filters.InlineQuery, r Response) InlineQuery {
	return InlineQuery{
		Filter:   f,
		Response: r,
	}
}

func (i InlineQuery) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return i.Response(b, ctx)
}

func (i InlineQuery) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.InlineQuery == nil {
		return false
	}

	return i.Filter == nil || i.Filter(ctx.InlineQuery)
}

func (i InlineQuery) Name() string {
	return fmt.Sprintf("inlinequery_%p", i.Response)
}
