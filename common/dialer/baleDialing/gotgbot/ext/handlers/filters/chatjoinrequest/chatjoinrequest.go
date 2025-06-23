package chatjoinrequest

import (
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

func All(_ *gotgbot.ChatJoinRequest) bool {
	return true
}

func ChatID(id int64) filters.ChatJoinRequest {
	return func(r *gotgbot.ChatJoinRequest) bool {
		return r.Chat.Id == id
	}
}
