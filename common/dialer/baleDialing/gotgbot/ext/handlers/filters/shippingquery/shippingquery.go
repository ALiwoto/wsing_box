package shippingquery

import (
	"strings"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

func All(_ *gotgbot.ShippingQuery) bool {
	return true
}

func FromUserID(id int64) filters.ShippingQuery {
	return func(s *gotgbot.ShippingQuery) bool {
		return s.From.Id == id
	}
}

func HasPayloadPrefix(pre string) filters.ShippingQuery {
	return func(s *gotgbot.ShippingQuery) bool {
		return strings.HasPrefix(s.InvoicePayload, pre)
	}
}
