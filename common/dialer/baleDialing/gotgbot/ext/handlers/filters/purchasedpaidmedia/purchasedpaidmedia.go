package purchasedpaidmedia

import (
	"strings"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

func All(_ *gotgbot.PaidMediaPurchased) bool {
	return true
}

func FromUserID(id int64) filters.PurchasedPaidMedia {
	return func(s *gotgbot.PaidMediaPurchased) bool {
		return s.From.Id == id
	}
}

func HasPayloadPrefix(pre string) filters.PurchasedPaidMedia {
	return func(s *gotgbot.PaidMediaPurchased) bool {
		return strings.HasPrefix(s.PaidMediaPayload, pre)
	}
}
