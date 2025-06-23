package helpPlugin

import (
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers"
)

func LoadHandlers(d *ext.Dispatcher, t []rune) {
	startCmd := handlers.NewCommand(startCommand, startHandler)
	helpCmd := handlers.NewCommand(helpCommand, helpHandler)
	idCmd := handlers.NewCommand(idCommand, idHandler)

	startCmd.Triggers = append(startCmd.Triggers, t...)
	helpCmd.Triggers = append(helpCmd.Triggers, t...)
	idCmd.Triggers = append(idCmd.Triggers, t...)

	d.AddHandler(startCmd)
	d.AddHandler(helpCmd)
	d.AddHandler(idCmd)
}
