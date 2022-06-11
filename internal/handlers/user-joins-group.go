package handlers

import (
	"fmt"
	"pingerbot/internal/messages"
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type UserJoinsGroup struct {
	S state.State
}

func (UserJoinsGroup) Name() string {
	return "UserJoinsGroup"
}

func (UserJoinsGroup) Match(u tg.Update) bool {
	return u.Message != nil && u.Message.NewMember != nil
}

func (h UserJoinsGroup) Handle(u tg.Update, ctx tg.Ctx) error {
	if u.Message.NewMember.Id == ctx.BotId {
		ctx.Logger.Debug("Skipping message about self")
		return nil
	}

	// We don't work with users that don't have username because
	// there is no way to mention them.
	if u.Message.NewMember.Username == "" {
		msg := messages.AddUsername(u.Message.Chat.Id, *u.Message.NewMember)
		_, err := ctx.Actions.SendMessage(msg)

		return err
	}

	err := h.S.RememberMember(u.Message.Chat.Id, *u.Message.NewMember)
	if err != nil {
		return err
	}

	msg := tg.OutgoingMessage{
		ChatId: u.Message.Chat.Id,
		Text:   fmt.Sprintf("Hi @%s! I know you now!", u.Message.NewMember.Username),
	}

	_, err = ctx.Actions.SendMessage(msg)

	return err
}
