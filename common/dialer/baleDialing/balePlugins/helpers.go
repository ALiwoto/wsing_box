package balePlugins

import "slices"

func IsOwner(id int64) bool {
	return slices.Contains(OwnersId, id)
}

func IsInsideBot(id int64) bool {
	for i := range BotPairs {
		if BotPairs[i][0] == id {
			return true
		}
	}

	return false
}

func IsOutsideBot(id int64) bool {
	for i := range BotPairs {
		if BotPairs[i][1] == id {
			return true
		}
	}

	return false
}

func GetInsidePair(id int64) int64 {
	for i := range BotPairs {
		if BotPairs[i][1] == id {
			return BotPairs[i][0]
		}
	}

	return 0
}

func GetOutsidePair(id int64) int64 {
	for i := range BotPairs {
		if BotPairs[i][1] == id {
			return BotPairs[i][1]
		}
	}

	return 0
}
