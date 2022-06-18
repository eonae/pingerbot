package handlers

import (
	"fmt"
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type UserLeavesGroupHandler struct {
	S state.State
}

func (h UserLeavesGroupHandler) Handle(ctx tg.LeaveCtx) error {
	err := h.S.ForgetMember(state.GroupId(ctx.ChatId), ctx.Subject.Username, []string{})
	if err != nil {
		return err
	}

	name := ctx.Subject.Username
	if name == "" {
		name = ctx.Subject.FirstName
	}

	return ctx.SendTxt(fmt.Sprintf("Farewell, mr. %s!", name))
}
