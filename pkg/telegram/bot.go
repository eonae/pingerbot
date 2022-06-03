package telegram

import (
	"fmt"
	"time"
)

// Bot configuration parameters
type BotConfig struct {
	Token    string
	Interval time.Duration
}

// Bot itself
type Bot struct {
	api      Api
	interval time.Duration
	handler  UpdateHandler
	offset   int
}

// Function for processing updates
type UpdateHandler func(Update)

// Bot constructor/initializer
func NewBot(config BotConfig, handler UpdateHandler) Bot {
	api := createApi(config.Token)
	return Bot{
		api:      api,
		interval: config.Interval,
		handler:  handler,
		offset:   0,
	}
}

// Start polling
func (b *Bot) Start() {
	fmt.Println("Bot started...")
	for {
		updates, err := b.api.GetUpdates(b.offset)
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, u := range updates {
			b.handler(u)
			b.offset = u.UpdateId + 1
		}
	}
}
