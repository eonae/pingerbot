package handlers

import (
	"pingerbot/internal/state"
	"pingerbot/pkg/telegram"
)

type BotLeavesGroup struct {
	S state.State
}

func (BotLeavesGroup) Name() string {
	return "BotLeavesGroup"
}

func (BotLeavesGroup) Match(u telegram.Update) bool {
	return u.MyChatMember != nil && u.MyChatMember.NewChatMember.Status == "left"
}

func (h BotLeavesGroup) Handle(u telegram.Update, ctx telegram.Ctx) error {
	return h.S.ForgetGroup(u.MyChatMember.Chat)
}
