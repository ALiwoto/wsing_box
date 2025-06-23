package ext

import "github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"

type Handler interface {
	// CheckUpdate checks whether the update should handled by this handler.
	CheckUpdate(b *gotgbot.Bot, ctx *Context) bool
	// HandleUpdate processes the update.
	HandleUpdate(b *gotgbot.Bot, ctx *Context) error
	// Name gets the handler name; used to differentiate handlers programmatically. Names should be unique.
	Name() string
}
