package handlers

import (
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
)

type Response func(b *gotgbot.Bot, ctx *ext.Context) error
