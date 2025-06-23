package helpPlugin

import (
	"github.com/ALiwoto/ssg/ssg"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/balePlugins"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/log"
)

func startHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	// senderChat := message.SenderChat
	chat := ctx.EffectiveChat
	if chat == nil || chat.Type != "private" {
		return ext.ContinueGroups
	}

	txt := "Hello!! I'm here to search musics and musical notes for you!\n"
	txt += "I'm still in development, but I will be working soon!"

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

func helpHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	// senderChat := message.SenderChat
	chat := ctx.EffectiveChat
	if chat == nil || chat.Type != "private" {
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

func idHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	user := ctx.EffectiveUser
	chat := ctx.EffectiveChat
	if chat == nil {
		return ext.ContinueGroups
	}

	txt := "-- ID --\n"
	txt += "🔹 Bot ID: " + ssg.ToBase10(bot.Id) + "\n"
	txt += "🔹 User ID: " + ssg.ToBase10(user.Id) + "\n"
	txt += "🔹 Chat ID: " + ssg.ToBase10(chat.Id) + "\n"

	if balePlugins.IsOwner(user.Id) {
		txt += "(You are a cool owner btw!)\n"
	}

	// sEncoded := singingEncoding.StdEncoding.EncodeToString([]byte(txt))
	// myCount := utf8.RuneCountInString(sEncoded)
	// print(sEncoded, myCount)

	// sDecoded, myErr := singingEncoding.StdEncoding.DecodeString(sEncoded)
	// print(sDecoded, myErr)

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
