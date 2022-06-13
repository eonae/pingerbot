package handlers

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type BotLeavesGroupHandler struct {
	S state.State
}

func (h BotLeavesGroupHandler) Handle(ctx tg.LeaveCtx) error {
	return h.S.ForgetGroup(ctx.Chat)
}
