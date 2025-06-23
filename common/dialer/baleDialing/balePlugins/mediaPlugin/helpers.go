package mediaPlugin

import (
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext/handlers"
)

func LoadHandlers(d *ext.Dispatcher, t []rune) {
	uploadMusicCmd := handlers.NewCommand(uploadMusicCommand, uploadMusicHandler)
	downloadMusicCmd := handlers.NewCommand(downloadMusicCommand, downloadMusicHandler)

	uploadMusicCmd.Triggers = append(uploadMusicCmd.Triggers, t...)
	downloadMusicCmd.Triggers = append(downloadMusicCmd.Triggers, t...)

	d.AddHandler(uploadMusicCmd)
	d.AddHandler(downloadMusicCmd)
}
