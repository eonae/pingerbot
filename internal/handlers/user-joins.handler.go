package handlers

import (
	"fmt"
	"pingerbot/internal/messages"
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type UserJoinsGroupHandler struct {
	S state.State
}

func (h UserJoinsGroupHandler) Handle(ctx tg.JoinCtx) error {
	if ctx.Subject.Id == ctx.BotId {
		ctx.Logger.Debug("Skipping message about self")
		return nil
	}

	// We don't work with users that don't have username because
	// there is no way to mention them.
	if ctx.Subject.Username == "" {
		err := ctx.SendToChat(messages.PleaseAddUsername(ctx.Subject))

		return err
	}

	err := h.S.RememberMember(state.GroupId(ctx.ChatId), ctx.Subject.Username, []string{})
	if err != nil {
		return err
	}

	return ctx.SendTxt(fmt.Sprintf("Hi @%s! I know you now!", ctx.Subject.Username))
}
