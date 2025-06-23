package handlers

import (
	"fmt"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

type Poll struct {
	Filter   filters.Poll
	Response Response
}

func NewPoll(f filters.Poll, r Response) Poll {
	return Poll{
		Filter:   f,
		Response: r,
	}
}

func (r Poll) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.Poll == nil {
		return false
	}
	return r.Filter == nil || r.Filter(ctx.Poll)
}

func (r Poll) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return r.Response(b, ctx)
}

func (r Poll) Name() string {
	return fmt.Sprintf("poll_%p", r.Response)
}
