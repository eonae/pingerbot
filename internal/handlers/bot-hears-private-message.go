package handlers

import (
	"pingerbot/internal/state"
	"pingerbot/pkg/telegram"
)

type BotHearsPrivateMessage struct {
	S state.State
}

func (BotHearsPrivateMessage) Name() string {
	return "BotHearsPrivateMessage"
}

func (BotHearsPrivateMessage) Match(u telegram.Update) bool {
	return u.Message != nil && u.Message.Chat.Type == "private"
}

func (BotHearsPrivateMessage) Handle(u telegram.Update, ctx telegram.Ctx) error {
	_, err := ctx.Actions.SendMessage(telegram.SendMessage{
		ChatId:  u.Message.Chat.Id,
		Text:    "I'm sorry, but i wasn't tought to talk yet",
		ReplyTo: u.Message.MessageId,
	})

	return err
}
