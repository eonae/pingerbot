package handlers

import (
	"fmt"
	"pingerbot/internal/state"
	"pingerbot/pkg/telegram"
)

type UserJoinsGroup struct {
	S state.State
}

func (UserJoinsGroup) Name() string {
	return "UserJoinsGroup"
}

func (UserJoinsGroup) Match(u telegram.Update) bool {
	return u.Message != nil && u.Message.NewMember != nil
}

func (h UserJoinsGroup) Handle(u telegram.Update, ctx telegram.Ctx) error {
	if u.Message.NewMember.Id == ctx.BotId {
		ctx.Logger.Debug("Skipping message about self")
		return nil
	}

	// We don't work with users that don't have username because
	// there is no way to mention them.
	if u.Message.NewMember.Username == "" {
		msg := telegram.SendMessage{
			ChatId:    u.Message.Chat.Id,
			ParseMode: telegram.Markdown,
			Text: fmt.Sprintf(
				"I can't ping users without username mr.[%s](tg://user?id=%d). Please setup yours!",
				u.Message.NewMember.FirstName,
				u.Message.NewMember.Id,
			),
		}

		_, err := ctx.Actions.SendMessage(msg)

		return err
	}

	err := h.S.RememberMember(u.Message.Chat.Id, *u.Message.NewMember)
	if err != nil {
		return err
	}

	msg := telegram.SendMessage{
		ChatId: u.Message.Chat.Id,
		Text:   fmt.Sprintf("Hi @%s! I know you now!", u.Message.NewMember.Username),
	}

	_, err = ctx.Actions.SendMessage(msg)

	return err
}
