package handlers

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type BotLeavesGroup struct {
	S state.State
}

func (BotLeavesGroup) Name() string {
	return "BotLeavesGroup"
}

func (BotLeavesGroup) Match(u tg.Update) bool {
	return u.MyChatMember != nil && u.MyChatMember.NewChatMember.Status == "left"
}

func (h BotLeavesGroup) Handle(u tg.Update, ctx tg.Ctx) error {
	return h.S.ForgetGroup(u.MyChatMember.Chat)
}
