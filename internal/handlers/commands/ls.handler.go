package commands

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type LsCommandHandler struct {
	S state.State
}

func (h LsCommandHandler) Handle(ctx tg.CommandCtx) error {
	tags := ctx.Tags()

	list, err := h.S.ListGroupMembers(ctx.ChatId, tags)
	if err != nil {
		return err
	}

	err = ctx.ReplyTxt(list.String())

	if err != nil {
		return err
	}

	return nil
}
