package pollanswer

import (
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

func All(_ *gotgbot.PollAnswer) bool {
	return true
}

func Id(id string) filters.PollAnswer {
	return func(p *gotgbot.PollAnswer) bool {
		return p.PollId == id
	}
}

func FromUserId(id int64) filters.PollAnswer {
	return func(p *gotgbot.PollAnswer) bool {
		return p.User.Id == id
	}
}
