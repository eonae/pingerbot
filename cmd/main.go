package main

import (
	"database/sql"
	"os"
	"pingerbot/internal/handlers"
	"pingerbot/internal/handlers/commands"
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

	db, err := sql.Open("postgres", config.Storage.Url)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.SetMaxOpenConns(10)

	state := state.New(db)

	bot := telegram.NewBot(config.Bot, telegram.Handlers{
		PrivateMessages: handlers.PrivateMessageHandler{S: state},
		SelfJoin:        handlers.BotJoinsGroupHandler{S: state},
		SelfLeave:       handlers.BotLeavesGroupHandler{S: state},
		UserJoin:        handlers.UserJoinsGroupHandler{S: state},
		UserLeave:       handlers.UserLeavesGroupHandler{S: state},
		PublicCommands: map[string]telegram.CommandHandler{
			"/add":      commands.AddCommandHandler{S: state},
			"/addme":    commands.AddmeCommandHandler{S: state},
			"/remove":   commands.RemoveCommandHandler{S: state},
			"/removeme": commands.RemovemeCommandHandler{S: state},
			"/ping":     commands.PingCommandHandler{S: state},
		},
	})

	bot.Start()
}
