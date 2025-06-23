package handlers

import (
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
)

type Named struct {
	// Custom name to identify handler by
	CustomName string
	// Inlined version of parent handler to inherit methods.
	ext.Handler
}

func (n Named) Name() string {
	return n.CustomName
}

func NewNamedhandler(name string, handler ext.Handler) Named {
	return Named{
		CustomName: name,
		Handler:    handler,
	}
}
