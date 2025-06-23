package handlers

import (
	"fmt"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

type PollAnswer struct {
	Filter   filters.PollAnswer
	Response Response
}

func NewPollAnswer(f filters.PollAnswer, r Response) PollAnswer {
	return PollAnswer{
		Filter:   f,
		Response: r,
	}
}

func (r PollAnswer) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.PollAnswer == nil {
		return false
	}
	return r.Filter == nil || r.Filter(ctx.PollAnswer)
}

func (r PollAnswer) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return r.Response(b, ctx)
}

func (r PollAnswer) Name() string {
	return fmt.Sprintf("pollanswer_%p", r.Response)
}
