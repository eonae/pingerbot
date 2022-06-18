package commands

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type RemoveCommandHandler struct {
	S state.State
}

func (h RemoveCommandHandler) Handle(ctx tg.CommandCtx) error {
	mentions, tags := ctx.Mentions(), ctx.Tags()

	if len(mentions) == 0 {
		return ctx.ReplyTxt("Please provide some usernames!")
	}

	for _, mention := range mentions {
		ctx.Logger.Infof("Forgetting user %s", mention)

		err := h.S.ForgetMember(state.GroupId(ctx.ChatId), mention, tags)
		if err != nil {
			ctx.Logger.Error(err)
		}
	}

	return nil
}
