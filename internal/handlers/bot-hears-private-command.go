package handlers

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type BotHearsPrivateCommand struct {
	S state.State
}

func (h BotHearsPrivateCommand) Handle(ctx tg.CommandCtx) error {
	return ctx.Reply(tg.OutgoingMessage{Text: "Sorry! No private command yet!"})
}
