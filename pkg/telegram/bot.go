package telegram

import (
	"time"

	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
)

// Bot configuration parameters
type BotConfig struct {
	Token   string
	Timeout time.Duration
}

// Bot itself
type Bot struct {
	api      Api
	timeout  time.Duration
	handlers []Handler
	offset   int64
}

type Ctx struct {
	BotId   int64
	BotName string
	Actions Api
	Logger  *logrus.Entry
}

type Handler interface {
	Name() string
	Match(u Update) bool
	Handle(u Update, ctx Ctx) error
}

// Bot constructor/initializer
func NewBot(config BotConfig) Bot {
	api := NewApi(config.Token)
	return Bot{
		api:      api,
		timeout:  config.Timeout,
		handlers: []Handler{},
		offset:   0,
	}
}

func (b *Bot) AddHandler(h Handler) {
	b.handlers = append(b.handlers, h)
}

// Start polling
func (b *Bot) Start() {

	me, err := b.api.GetMe()
	if err != nil {
		logrus.Fatal("Failed to get bot profile!")
	}

	logger := logrus.WithFields(logrus.Fields{
		"botId":   me.Id,
		"botName": "@" + me.Username,
	})

	logger.Infof("BotApplication started! Processing updates...")

	for {
		logger.Debugf("Polling with timeout %s", b.timeout)
		updates, err := b.api.GetUpdates(b.offset, b.timeout)
		if err != nil {
			logrus.Error("Failed to get updates!", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, u := range updates {
			logger := logrus.WithField("updateId", u.UpdateId)
			ctx := Ctx{
				BotId:   me.Id,
				BotName: "@" + me.Username,
				Actions: b.api,
				Logger:  logger,
			}

			logrus.Debugf("Processing update:%# v", pretty.Formatter(u))

			handler, ok := b.getHandler(u)
			if !ok {
				logger.Debug("No handler found. Update is skipped")
				b.offset = u.UpdateId + 1
				continue
			}

			logger.Debug("Handler found:", handler.Name())
			ctx.Logger = ctx.Logger.WithField("handler", handler.Name())

			err := handler.Handle(u, ctx)
			if err != nil {
				logger.Error("Failed to process update", err)
			} else {
				logger.Debug("Update handled successfully")
			}
			b.offset = u.UpdateId + 1
		}
	}
}

func (b Bot) getHandler(u Update) (Handler, bool) {
	for _, h := range b.handlers {
		if h.Match(u) {
			return h, true
		}
	}
	return nil, false
}
