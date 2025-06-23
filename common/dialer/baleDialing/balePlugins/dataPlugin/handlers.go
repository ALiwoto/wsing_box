package dataPlugin

import (
	"strings"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/balePlugins"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/log"
)

func dataMessageHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	isInsideBot := balePlugins.IsInsideBot(bot.Id)
	isOutsideBot := balePlugins.IsOutsideBot(bot.Id)

	if isInsideBot && strings.HasPrefix(message.Text, balePlugins.InsidePreData) {
		return ext.ContinueGroups
	}

	if isOutsideBot && strings.HasPrefix(message.Text, balePlugins.OutsidePreData) {
		return ext.ContinueGroups
	}

	if !isInsideBot && !isOutsideBot {
		// perhaps this is a helper bot or something?
		return ext.ContinueGroups
	}

	processMessageLock.Lock()
	defer processMessageLock.Unlock()

	messageKey := ssg.ToBase10(chat.Id) + "_" + ssg.ToBase10(message.MessageId)
	if processedMessages.Exists(messageKey) {
		return ext.ContinueGroups
	}
	processedMessages.Set(messageKey, true)

	_, _ = bot.DeleteMessage(chat.Id, message.MessageId, nil)

	// process rest of the messages here...
	// The format of data is:
	// F-<conn_id> [DATA] or F-<conn_id> [CMD]
	// S-<conn_id> [DATA] or S-<conn_id> [CMD]
	// but only data gets handled here.
	firstPart, incomingData, found := strings.Cut(message.Text, " ")
	if !found {
		log.Error("got invalid data: " + message.Text)
		return ext.ContinueGroups
	}
	connId := strings.Split(firstPart, "-")[1]

	if isInsideBot {
		connPtr := balePlugins.BaleConnectionsPool.Get(connId)
		if connPtr == nil {
			log.Error("non-existing connections are not handled yet...")
			return ext.ContinueGroups
		}

		conn := *connPtr
		_, _ = conn.Write([]byte(incomingData))
	} else if isOutsideBot {
		// not handled yet... do something here I guess
		log.Error("not handled yet...")
	}

	// don't let another handlers get executed
	return ext.ContinueGroups
}

func dataMessageFilter(msg *gotgbot.Message) bool {
	return strings.HasPrefix(msg.Text, balePlugins.InsidePreData) ||
		strings.HasPrefix(msg.Text, balePlugins.OutsidePreData)
}
