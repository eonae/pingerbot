package handlers

import (
	"pingerbot/internal/messages"
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type BotJoinsGroup struct {
	S state.State
}

func (h BotJoinsGroup) Handle(ctx tg.JoinCtx) (err error) {
	err = h.S.RememberGroup(ctx.Chat)
	if err != nil {
		return err
	}

	err = ctx.SendToChat(messages.Welcome)
	if err != nil {
		return err
	}

	if ctx.Actor.Username == "" {
		return ctx.SendToChat(messages.PleaseAddUsername(ctx.Actor))
	}

	return h.S.RememberMember(ctx.Chat.Id, ctx.Actor.Username)
}
