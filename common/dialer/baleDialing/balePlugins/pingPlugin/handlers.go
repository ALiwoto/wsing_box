package pingPlugin

import (
	"strings"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/balePlugins"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/log"
)

func musicalCheckHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	if !balePlugins.IsInsideBot(bot.Id) {
		return ext.ContinueGroups
	}

	message := ctx.EffectiveMessage
	// senderChat := message.SenderChat
	// chat := ctx.EffectiveChat

	_, err := message.Reply(
		bot,
		MusicalInOnCommand+" ["+ssg.ToBase10(bot.Id)+"]",
		&gotgbot.SendMessageOpts{},
	)

	if err != nil {
		log.Error("Failed to send musicalInOn command:", err)
		return ext.EndGroups
	}

	// don't let another handlers get executed
	return ext.EndGroups
}

func musicalCheckFilter(msg *gotgbot.Message) bool {
	return strings.HasPrefix(msg.Text, MusicalCheckCommand)
}

func musicalInPairHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	if !balePlugins.IsOutsideBot(bot.Id) {
		return ext.ContinueGroups
	}

	message := ctx.EffectiveMessage
	// senderChat := message.SenderChat
	// chat := ctx.EffectiveChat

	incomingId := ssg.ToInt64(ssg.Split(message.Text, "[", "]")[1])
	if balePlugins.GetInsidePair(bot.Id) != incomingId {
		return ext.ContinueGroups
	}

	_, err := message.Reply(
		bot,
		MusicalOutOnCommand+" ["+ssg.ToBase10(bot.Id)+"]",
		&gotgbot.SendMessageOpts{},
	)

	if err != nil {
		log.Error("Failed to send musicalInOn command:", err)
		return ext.EndGroups
	}

	// don't let another handlers get executed
	return ext.EndGroups
}

func musicalInPairFilter(msg *gotgbot.Message) bool {
	return strings.HasPrefix(msg.Text, MusicalInOnCommand)
}
