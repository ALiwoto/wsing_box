package balePlugins

import "slices"

func IsOwner(id int64) bool {
	return slices.Contains(OwnersId, id)
}

func IsInsideBot(id int64) bool {
	for i := range BotPairs {
		if BotPairs[i].InsideBotId == id {
			return true
		}
	}

	return false
}

func IsOutsideBot(id int64) bool {
	for i := range BotPairs {
		if BotPairs[i].OutsideBotId == id {
			return true
		}
	}

	return false
}

func GetInsidePair(id int64) int64 {
	for i := range BotPairs {
		if BotPairs[i].OutsideBotId == id {
			return BotPairs[i].InsideBotId
		}
	}

	return 0
}

func GetOutsidePair(id int64) int64 {
	for i := range BotPairs {
		if BotPairs[i].InsideBotId == id {
			return BotPairs[i].OutsideBotId
		}
	}

	return 0
}

func GetMyChats(botId int64) []int64 {
	for i := range BotPairs {
		if BotPairs[i].InsideBotId == botId || BotPairs[i].OutsideBotId == botId {
			return BotPairs[i].ChatIds
		}
	}

	return []int64{}
}

func IsInCorrectChat(botId, chatId int64) bool {
	return slices.Contains(GetMyChats(botId), chatId)
}
