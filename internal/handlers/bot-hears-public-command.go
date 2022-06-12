package handlers

import (
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
	default:
		ctx.Logger.Warnf("Unknown command: %s", ctx.Command)
	}

	return nil
}
