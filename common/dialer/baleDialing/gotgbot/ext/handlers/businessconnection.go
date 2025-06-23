package handlers

import (
	"fmt"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

type BusinessConnection struct {
	Filter   filters.BusinessConnection
	Response Response
}

func (bc BusinessConnection) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	return ctx.BusinessConnection != nil && bc.Filter(ctx.BusinessConnection)
}

func (bc BusinessConnection) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return bc.Response(b, ctx)
}

func (bc BusinessConnection) Name() string {
	return fmt.Sprintf("businessconnection_%p", bc.Response)
}
