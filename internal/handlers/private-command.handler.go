package handlers

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type PrivateCommandHandler struct {
	S state.State
}

func (h PrivateCommandHandler) Handle(ctx tg.CommandCtx) error {
	return ctx.ReplyTxt("Sorry! No private command yet!")
}
