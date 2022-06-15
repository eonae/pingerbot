package commands

import (
	"pingerbot/internal/messages"
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type AddmeCommandHandler struct {
	S state.State
}

func (h AddmeCommandHandler) Handle(ctx tg.CommandCtx) error {
	user := ctx.Message.From
	if user.Username == "" {
		ctx.Logger.Debugf("Can't add user %s - no username!", user.FirstName)
		return ctx.SendToChat(messages.PleaseAddUsername(user))
	}

	ctx.Logger.Infof("Remembering user @%s", user.Username)
	return h.S.RememberMember(ctx.ChatId, user.Username, ctx.Tags())
}
