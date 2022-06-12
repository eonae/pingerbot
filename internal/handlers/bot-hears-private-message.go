package handlers

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type BotHearsPrivateMessage struct {
	S state.State
}

func (BotHearsPrivateMessage) Handle(ctx tg.MsgCtx) error {
	return ctx.Reply(tg.OutgoingMessage{
		Text: "I'm sorry, but i wasn't tought to talk yet",
	})
}
