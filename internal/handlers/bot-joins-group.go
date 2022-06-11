package handlers

import (
	"fmt"
	"pingerbot/internal/state"
	"pingerbot/pkg/telegram"
)

type BotJoinsGroup struct {
	S state.State
}

func (BotJoinsGroup) Name() string {
	return "BotJoinsGroup"
}

func (BotJoinsGroup) Match(u telegram.Update) bool {
	return u.MyChatMember != nil && u.MyChatMember.NewChatMember.Status != "left"
}

func (h BotJoinsGroup) Handle(u telegram.Update, ctx telegram.Ctx) (err error) {
	err = h.S.RememberGroup(u.MyChatMember.Chat)
	if err != nil {
		return err
	}

	msg := telegram.SendMessage{
		ChatId:    u.MyChatMember.Chat.Id,
		ParseMode: telegram.Markdown,
		Text: `Hi all!

I'm mr.Pinger! If you put **/ping** into you message, i will notify everyone.

_Please note. I don't know people in this chat yet. But I remember everyone who writes something._`,
	}

	_, err = ctx.Actions.SendMessage(msg)
	if err != nil {
		return err
	}

	if u.MyChatMember.From.Username == "" {
		msg := telegram.SendMessage{
			ChatId:    u.MyChatMember.Chat.Id,
			ParseMode: telegram.Markdown,
			Text: fmt.Sprintf(
				"I can't ping users without username mr.[%s](tg://user?id=%d). Please setup yours!",
				u.MyChatMember.From.FirstName,
				u.MyChatMember.From.Id,
			),
		}

		_, err = ctx.Actions.SendMessage(msg)

		return err
	}

	err = h.S.RememberMember(u.MyChatMember.Chat.Id, u.MyChatMember.From)

	return err
}
