package mediaPlugin

import (
	"os"
	"strings"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/balePlugins"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/log"
)

func uploadMusicHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	user := ctx.EffectiveUser
	chat := ctx.EffectiveChat
	if chat == nil || user == nil || !balePlugins.IsOwner(user.Id) {
		return ext.ContinueGroups
	}

	// 0-index is the command itself
	_, name, found := strings.Cut(message.Text, " ")
	if !found {
		txt := "Usage: /uploadMusic name"
		_, err := message.Reply(
			bot,
			txt,
			&gotgbot.SendMessageOpts{},
		)
		if err != nil {
			log.Error("Failed to send start txt reply:", err)
			return ext.EndGroups
		}
		return ext.EndGroups
	}

	if strings.Contains(name, "..") {
		// prevent path traversal vulnerability
		return ext.EndGroups
	}

	if !strings.HasSuffix(name, ".txt") {
		name += ".txt"
	}

	targetFile, err := os.Open("./musics/" + name)
	if err != nil {
		txt := "Failed to send music: " + err.Error()
		_, err := message.Reply(
			bot,
			txt,
			&gotgbot.SendMessageOpts{},
		)
		if err != nil {
			log.Error("Failed to send start txt reply:", err)
			return ext.EndGroups
		}
		return ext.EndGroups
	}

	_, err = bot.SendDocument(
		chat.Id,
		gotgbot.InputFileByReader("test.txt", targetFile),
		&gotgbot.SendDocumentOpts{
			// ReplyParameters: &gotgbot.ReplyParameters{
			// 	MessageId: message.MessageId,
			// },
		},
	)
	if err != nil {
		log.Error("Failed to send start txt reply:", err)
		return ext.EndGroups
	}

	// don't let another handlers get executed
	return ext.EndGroups
}

func downloadMusicHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	// senderChat := message.SenderChat
	chat := ctx.EffectiveChat
	if chat == nil {
		return ext.ContinueGroups
	}

	txt := "-- help section --\n"
	txt += "I'm here to search musics and musical notes for you!\n"
	txt += "Commands: \n"
	txt += "  /id: shows ID numbers\n"

	_, err := message.Reply(
		bot,
		txt,
		&gotgbot.SendMessageOpts{},
	)

	if err != nil {
		log.Error("Failed to send start txt reply:", err)
		return ext.EndGroups
	}

	// don't let another handlers get executed
	return ext.EndGroups
}
