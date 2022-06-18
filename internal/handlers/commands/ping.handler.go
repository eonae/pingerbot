package commands

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
	"strings"
)

type PingCommandHandler struct {
	S state.State
}

func (h PingCommandHandler) Handle(ctx tg.CommandCtx) error {
	members, err := h.S.GetKnownMembers(state.GroupId(ctx.ChatId), ctx.Tags())
	if err != nil {
		return err
	}

	if len(members) == 0 {
		return ctx.ReplyTxt("I don't know anyone yet!")
	}

	return ctx.ReplyTxt(strings.Join(members, " "))
}
