package handlers

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type BotLeavesGroup struct {
	S state.State
}

func (h BotLeavesGroup) Handle(ctx tg.LeaveCtx) error {
	return h.S.ForgetGroup(ctx.Chat)
}
