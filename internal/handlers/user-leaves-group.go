package handlers

import (
	"fmt"
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type UserLeavesGroup struct {
	S state.State
}

func (h UserLeavesGroup) Handle(ctx tg.LeaveCtx) error {
	err := h.S.ForgetMember(ctx.ChatId, ctx.Subject.Username)
	if err != nil {
		return err
	}

	name := ctx.Subject.Username
	if name == "" {
		name = ctx.Subject.FirstName
	}

	return ctx.SendToChat(tg.OutgoingMessage{
		Text: fmt.Sprintf("Farewell, mr. %s!", name),
	})
}
