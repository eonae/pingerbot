package commands

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type RemovemeCommandHandler struct {
	S state.State
}

func (h RemovemeCommandHandler) Handle(ctx tg.CommandCtx) error {
	return h.S.ForgetMember(ctx.ChatId, "@"+ctx.Message.From.Username, ctx.Tags())
}
