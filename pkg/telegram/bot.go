package telegram

import (
	"errors"
	"pingerbot/pkg/helpers"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
)

const (
	WrongBotErr = helpers.Error("Command for another bot!")
	ParseCmdErr = helpers.Error("Couldn't parse command!")
)

// Bot configuration parameters
type BotConfig struct {
	Token   string
	Timeout time.Duration
}

type MessageHandler interface {
	Handle(ctx MsgCtx) error
}

type CommandHandler interface {
	Handle(ctx CommandCtx) error
}

type JoinHandler interface {
	Handle(ctx JoinCtx) error
}

type LeaveHandler interface {
	Handle(ctx LeaveCtx) error
}

type Handlers struct {
	PrivateMessages MessageHandler
	SelfJoin        JoinHandler
	SelfLeave       LeaveHandler
	UserJoin        JoinHandler
	UserLeave       LeaveHandler
	PublicCommands  map[string]CommandHandler
	PrivateCommands map[string]CommandHandler
}

// Bot itself
type Bot struct {
	offset   int64
	timeout  time.Duration
	api      Api
	handlers Handlers
}

// Bot constructor/initializer
func NewBot(config BotConfig, handlers Handlers) Bot {
	api := NewApi(config.Token)
	return Bot{
		offset:   0,
		timeout:  config.Timeout,
		api:      api,
		handlers: handlers,
	}
}

// Start polling
func (b *Bot) Start() {

	me, err := b.api.GetMe()
	if err != nil {
		logrus.Fatal("Failed to get bot profile!")
	}

	botId, botName := me.Id, "@"+me.Username

	logger := logrus.WithFields(logrus.Fields{
		"botId":   botId,
		"botName": botName,
	})

	ctx := Ctx{
		BotId:   botId,
		BotName: botName,
		api:     b.api,
		Logger:  logger,
	}

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
			if u.MyChatMember != nil {
				ctx.ChatId = u.MyChatMember.Chat.Id
			} else if u.Message != nil {
				ctx.ChatId = u.Message.Chat.Id
			} else {
				logger.Error("Failed to process update: ", u)
				b.offset = u.UpdateId + 1
				continue
			}

			err := b.handle(u, ctx)
			if err != nil {
				logger.Error("Failed to process update", err)
			} else {
				logger.Debug("Update handled successfully")
			}
			b.offset = u.UpdateId + 1
		}
	}
}

func (b Bot) handle(u Update, ctx Ctx) error {
	logger := logrus.WithField("updateId", u.UpdateId)

	logger.Debugf("Processing update:%# v", pretty.Formatter(u))

	if u.MyChatMember != nil {
		if u.MyChatMember.NewChatMember.Status == "left" {
			return b.handlers.SelfLeave.Handle(ctxFromLeaveEvent(ctx, *u.MyChatMember))
		}

		return b.handlers.SelfJoin.Handle(ctxFromJoinEvent(ctx, *u.MyChatMember))
	}

	if u.Message == nil {
		return errors.New("no handler found")
	}

	if u.Message.NewMember != nil {
		return b.handlers.UserJoin.Handle(joinCtxFromMessage(ctx, *u.Message))
	}

	if u.Message.LeftMember != nil {
		return b.handlers.UserLeave.Handle(leaveCtxFromMessage(ctx, *u.Message))
	}

	msgCtx := ctx.ToMsgContext(*u.Message)

	switch len(msgCtx.Commands()) {
	case 0:
		// No checks needed that it's private message because
		// Privacy mode is on.
		if b.handlers.PrivateMessages == nil {
			logger.Debugf("No handler for private messages")
		}

		return b.handlers.PrivateMessages.Handle(msgCtx)
	case 1:
		cmd, err := parseCmd(ctx.BotName, msgCtx.Commands()[0])
		switch err {
		case ParseCmdErr:
			return err
		case WrongBotErr:
			logger.Debugf("Skipping command because it's for some other bot")
		}

		var handlers map[string]CommandHandler

		if msgCtx.Message.Chat.Type == "private" {
			handlers = b.handlers.PrivateCommands
		} else {
			handlers = b.handlers.PrivateCommands
		}

		if handlers == nil {
			ctx.Logger.Warnf("Unknown command: %s", cmd)
			return nil
		}

		handler, ok := handlers[cmd]
		if ok {
			return handler.Handle(msgCtx.ToCommandCtx(cmd))
		}

		ctx.Logger.Warnf("Unknown command: %s", cmd)
		return nil
	default:
		err := msgCtx.ReplyTxt("More that one command contains in message!")
		if err != nil {
			return err
		}

		return errors.New("more that one command contains in message")
	}
}

func parseCmd(botName string, cmd string) (string, error) {
	parts := strings.Split(cmd, "@")
	switch len(parts) {
	case 1:
		return cmd, nil
	case 2:
		if "@"+parts[1] == botName {
			return parts[0], nil
		}
		return "", WrongBotErr
	default:
		return "", ParseCmdErr
	}
}
