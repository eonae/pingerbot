package commands

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type AddCommandHandler struct {
	S state.State
}

func (h AddCommandHandler) Handle(ctx tg.CommandCtx) error {
	mentions, tags := ctx.Mentions(), ctx.Tags()

	if len(mentions) == 0 {
		return ctx.ReplyTxt("Please provide some usernames!")
	}

	for _, mention := range mentions {
		if string(mention[0]) != "@" {
			ctx.Logger.Warn("Can't remember user %s - not username!", mention)
			continue
		}

		// TODO: Batch add
		ctx.Logger.Infof("Remembering user %s", mention)

		// Strip @ because it will be added inside of RememberMember method
		err := h.S.RememberMember(state.GroupId(ctx.ChatId), mention[1:], tags)
		if err != nil {
			ctx.Logger.Error(err)
		}
	}

	return nil
}
