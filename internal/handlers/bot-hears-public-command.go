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
		members, err := h.S.GetKnownMembers(ctx.ChatId, ctx.Tags())
		if err != nil {
			return err
		}

		if len(members) == 0 {
			return ctx.Reply(tg.OutgoingMessage{Text: "I don't know anyone yet!"})
		}

		return ctx.Reply(tg.OutgoingMessage{Text: strings.Join(members, " ")})
	case "/add":
		mentions, tags := ctx.Mentions(), ctx.Tags()

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
			err := h.S.RememberMember(ctx.ChatId, mention[1:], tags)
			if err != nil {
				ctx.Logger.Error(err)
			}
		}
	case "/remove":
		mentions, tags := ctx.Mentions(), ctx.Tags()

		if len(mentions) == 0 {
			return ctx.Reply(tg.OutgoingMessage{Text: "Please provide some usernames!"})
		}

		for _, mention := range mentions {
			ctx.Logger.Infof("Forgetting user %s", mention)

			err := h.S.ForgetMember(ctx.ChatId, mention, tags)
			if err != nil {
				ctx.Logger.Error(err)
			}
		}
	case "/addme":
		return rememberIfHasUsername(ctx.ChatId, ctx.Message.From, h.S, ctx.Logger, ctx, ctx.Tags())
	case "/removeme":
		return h.S.ForgetMember(ctx.ChatId, "@"+ctx.Message.From.Username, ctx.Tags())
	default:
		ctx.Logger.Warnf("Unknown command: %s", ctx.Command)
	}

	return nil
}
