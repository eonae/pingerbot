package handlers

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type BotHearsPrivateMessage struct {
	S state.State
}

func (BotHearsPrivateMessage) Name() string {
	return "BotHearsPrivateMessage"
}

func (BotHearsPrivateMessage) Match(u tg.Update) bool {
	return u.Message != nil && u.Message.Chat.Type == "private"
}

func (BotHearsPrivateMessage) Handle(u tg.Update, ctx tg.Ctx) error {
	_, err := ctx.Actions.SendMessage(tg.NewReply(*u.Message, tg.OutgoingMessage{
		Text: "I'm sorry, but i wasn't tought to talk yet",
	}))

	return err
}
