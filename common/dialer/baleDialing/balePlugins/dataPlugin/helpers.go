package dataPlugin

import (
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers"
)

func LoadHandlers(d *ext.Dispatcher, t []rune) {
	d.AddHandler(handlers.NewMessage(dataMessageFilter, dataMessageHandler))
}
