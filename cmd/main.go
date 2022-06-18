package main

import (
	"context"
	"os"
	"pingerbot/internal/handlers"
	"pingerbot/internal/handlers/commands"
	"pingerbot/internal/state"
	"pingerbot/pkg/telegram"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
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

	db, err := pgxpool.Connect(context.Background(), config.Storage.Url)
	if err != nil {
		panic(err)
	}

	defer db.Close()

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
			"/ls":       commands.LsCommandHandler{S: state},
		},
	})

	bot.Start()
}
