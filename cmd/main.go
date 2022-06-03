package main

import (
	"fmt"
	"os"
	"pingerbot/pkg/telegram"
	"time"
)

func parseConfig() telegram.BotConfig {
	token := os.Getenv("BOT_TOKEN")

	interval, err := time.ParseDuration(os.Getenv("POLLING_INTERVAL_MS"))
	if err != nil {
		panic(err)
	}

	return telegram.BotConfig{
		Token:    token,
		Interval: interval,
	}
}

/*
	Предположим, что порядок обработки апдейтов не важен.
	У нас есть очередь обработки апдейтов.
	И есть текущий оффсет до которого все апдейти обработаны.
	Оффсет должен быть равен минимальному значению в очереди.

	GetUpdates возвращает буферизированный канал из которого выпрыгивают апдейты
*/

func main() {
	config := parseConfig()

	bot := telegram.NewBot(config, func(u telegram.Update) {
		fmt.Printf("Handling update:\n%+v\n", u)
	})

	bot.Start()
}
