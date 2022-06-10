package main

import (
	"database/sql"
	"os"
	"pingerbot/internal/handlers"
	"pingerbot/internal/state"
	"pingerbot/pkg/telegram"
	"time"

	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type StorageConfig struct {
	Url string
}

type AppConfig struct {
	Bot     telegram.BotConfig
	Storage StorageConfig
}

func parseConfig() AppConfig {
	token := os.Getenv("BOT_TOKEN")

	timeout, err := time.ParseDuration(os.Getenv("LONG_POLLING_TIMEOUT"))
	if err != nil {
		panic(err)
	}

	connstr := os.Getenv("DB_CONNECTION")

	return AppConfig{
		Bot: telegram.BotConfig{
			Token:   token,
			Timeout: timeout,
		},
		Storage: StorageConfig{
			Url: connstr,
		},
	}
}

func configureLogger() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	configureLogger()

	config := parseConfig()

	bot := telegram.NewBot(config.Bot)

	db, err := sql.Open("postgres", config.Storage.Url)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.SetMaxOpenConns(10)

	state := state.New(db)

	bot.AddHandler(handlers.BotJoinsGroup{S: state})
	bot.AddHandler(handlers.BotLeavesGroup{S: state})
	bot.AddHandler(handlers.UserJoinsGroup{S: state})
	bot.AddHandler(handlers.UserLeavesGroup{S: state})
	bot.AddHandler(handlers.BotHearsPrivateMessage{S: state})
	bot.AddHandler(handlers.BotHearsPublicMessage{S: state})

	bot.Start()
}
