package handlers

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type PrivateMessageHandler struct {
	S state.State
}

func (PrivateMessageHandler) Handle(ctx tg.MsgCtx) error {
	return ctx.ReplyTxt("I'm sorry, but i wasn't tought to talk yet")
}
