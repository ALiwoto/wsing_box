package pingPlugin

import (
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers"
)

func LoadHandlers(d *ext.Dispatcher, t []rune) {
	musicalCheckCmd := handlers.NewMessage(musicalCheckFilter, musicalCheckHandler)
	musicalInPairCmd := handlers.NewMessage(musicalInPairFilter, musicalInPairHandler)

	d.AddHandler(musicalCheckCmd)
	d.AddHandler(musicalInPairCmd)
}
