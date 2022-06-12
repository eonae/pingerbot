package handlers

import (
	"fmt"
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
	"strings"
)

type BotHearsPublicCommand struct {
	S state.State
}

func (h BotHearsPublicCommand) Handle(ctx tg.CommandCtx) error {
	from := ctx.Message.From
	if from.Id == ctx.BotId {
		ctx.Logger.Debug("Skipping message from self")
		return nil
	}

	switch ctx.Command {
	case "/ping":
		members, err := h.S.GetKnownMembers(ctx.ChatId)
		if err != nil {
			return err
		}

		if len(members) == 0 {
			return ctx.Reply(tg.OutgoingMessage{Text: "I don't know anyone yet!"})
		}

		return ctx.Reply(tg.OutgoingMessage{Text: strings.Join(members, " ")})
	case "/add":
		mentions := ctx.Mentions()

		fmt.Println(mentions)

		if len(mentions) == 0 {
			return ctx.Reply(tg.OutgoingMessage{Text: "Please provide some usernames!"})
		}

		for _, mention := range mentions {
			if string(mention[0]) != "@" {
				ctx.Logger.Warn("Can't remember user %s - not username!", mention)
				continue
			}

			// TODO: Batch add
			ctx.Logger.Infof("Remembering user %s", mention)

			// Strip @ because it will be added inside of RememberMember method
			err := h.S.RememberMember(ctx.ChatId, mention[1:])
			if err != nil {
				ctx.Logger.Error(err)
			}
		}

	case "/addme":
		return rememberIfHasUsername(ctx.ChatId, ctx.Message.From, h.S, ctx.Logger, ctx)
	default:
		ctx.Logger.Warnf("Unknown command: %s", ctx.Command)
	}

	return nil
}
