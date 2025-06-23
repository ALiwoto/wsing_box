package baleDialing

import (
	"errors"
	"sync"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/option"
)

// github.com/sagernet/sing-box/protocol/socks/baleDialing

func NewBaleDialerContainer(opts option.DialerOptions) (*BaleDialerContainer, error) {
	dialerContainer := &BaleDialerContainer{
		Bots:     &BaleBotPairsContainer{},
		connPool: ssg.NewSafeMap[string, BaleConn](),
	}

	allPairs := opts.BaleConfig.BotPairs
	for index := range len(allPairs) {
		currentConfig := allPairs[index]
		currentPair := &BaleBotPair{
			InsideBot: &BaleBotContainer{
				MyConfig:   currentConfig.Inside,
				BaleConfig: opts.BaleConfig,
				lock:       &sync.Mutex{},
			},
			OutsideBot: &BaleBotContainer{
				MyConfig:   currentConfig.Outside,
				BaleConfig: opts.BaleConfig,
				lock:       &sync.Mutex{},
			},
		}

		anyBot := false
		if opts.BaleConfig.IsInside {
			if err := createBotInstance(opts.BaleConfig, currentPair.InsideBot); err != nil {
				return nil, err
			}
			anyBot = true
		}

		if opts.BaleConfig.IsOutside {
			if err := createBotInstance(opts.BaleConfig, currentPair.OutsideBot); err != nil {
				return nil, err
			}
			anyBot = true
		}

		if !anyBot {
			return nil, errors.New("no bale bots are configured")
		}

		dialerContainer.Bots.Pairs = append(
			dialerContainer.Bots.Pairs,
			currentPair,
		)
	}

	return dialerContainer, nil
}

func SubString(target string, from, to int) string {
	value := ""
	i := -1
	for _, current := range target {
		i++
		if i < from {
			continue
		}
		if i >= to {
			return value
		}

		value += string(current)
	}

	return value
}

func MakeChunks(target string, totalLen, maxLen int) []string {
	allChunks := []string{}
	for i := 0; i < totalLen; i += maxLen {
		end := min(i+maxLen, totalLen)

		chunk := SubString(target, i, end)
		allChunks = append(allChunks, chunk)
	}

	return allChunks
}

func createBotInstance(
	config *option.BaleConfiguration,
	container *BaleBotContainer,
) error {
	token := container.MyConfig.BotToken
	if len(token) == 0 {
		return errors.New("bot token is empty")
	}

	bot, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		RequestOpts: &gotgbot.RequestOpts{
			Timeout: 6 * gotgbot.DefaultTimeout,
			APIURL:  config.APIUrl,
		},
	})
	if err != nil {
		return err
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		MaxRoutines: MaxGoRoutines,
	})

	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{})
	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates:    true,
		EnableWebhookDeletion: false,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: 6 * gotgbot.DefaultTimeout,
				APIURL:  config.APIUrl,
			},
		},
	})
	if err != nil {
		return err
	}

	container.Bot = bot
	container.Dispatcher = dispatcher
	container.Updater = updater
	// logging.Info(fmt.Sprintf("%s has started | ID: %d", bot.Username, bot.Id))

	loadAllHandlers(dispatcher, cmdPrefixes)
	return nil
}

func loadAllHandlers(d *ext.Dispatcher, triggers []rune) {
	// pingPlugin.LoadHandlers(d, triggers)
	// helpPlugin.LoadHandlers(d, triggers)
	// mediaPlugin.LoadHandlers(d, triggers)
	// balePlugin.LoadHandlers(d, triggers)
}
