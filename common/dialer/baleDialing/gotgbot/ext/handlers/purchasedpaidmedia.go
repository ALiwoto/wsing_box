package handlers

import (
	"fmt"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers/filters"
)

type PurchasedPaidMedia struct {
	Filter   filters.PurchasedPaidMedia
	Response Response
}

func NewPurchasedPaidMedia(f filters.PurchasedPaidMedia, r Response) PurchasedPaidMedia {
	return PurchasedPaidMedia{
		Filter:   f,
		Response: r,
	}
}

func (r PurchasedPaidMedia) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.PreCheckoutQuery == nil {
		return false
	}
	return r.Filter == nil || r.Filter(ctx.PurchasedPaidMedia)
}

func (r PurchasedPaidMedia) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return r.Response(b, ctx)
}

func (r PurchasedPaidMedia) Name() string {
	return fmt.Sprintf("purchasedpaidmedia_%p", r.Response)
}
