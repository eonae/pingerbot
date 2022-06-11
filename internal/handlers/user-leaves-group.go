package handlers

import (
	"fmt"
	"pingerbot/internal/state"
	"pingerbot/pkg/telegram"
)

type UserLeavesGroup struct {
	S state.State
}

func (UserLeavesGroup) Name() string {
	return "UserJoinsGroup"
}

func (UserLeavesGroup) Match(u telegram.Update) bool {
	return u.Message != nil && u.Message.LeftMember != nil
}

func (h UserLeavesGroup) Handle(u telegram.Update, ctx telegram.Ctx) error {
	err := h.S.ForgetMember(u.Message.Chat.Id, *u.Message.LeftMember)
	if err != nil {
		return err
	}

	name := u.Message.LeftMember.Username
	if name == "" {
		name = u.Message.LeftMember.FirstName
	}

	msg := telegram.SendMessage{
		ChatId: u.Message.Chat.Id,
		Text:   fmt.Sprintf("Farewell, mr. %s!", name),
	}

	_, err = ctx.Actions.SendMessage(msg)

	return err
}
