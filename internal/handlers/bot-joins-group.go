package handlers

import (
	"pingerbot/internal/messages"
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type BotJoinsGroup struct {
	S state.State
}

func (BotJoinsGroup) Name() string {
	return "BotJoinsGroup"
}

func (BotJoinsGroup) Match(u tg.Update) bool {
	return u.MyChatMember != nil && u.MyChatMember.NewChatMember.Status != "left"
}

func (h BotJoinsGroup) Handle(u tg.Update, ctx tg.Ctx) (err error) {
	err = h.S.RememberGroup(u.MyChatMember.Chat)
	if err != nil {
		return err
	}

	msg := tg.OutgoingMessage{
		ChatId:    u.MyChatMember.Chat.Id,
		ParseMode: tg.Markdown,
		Text: `Hi all!

I'm mr.Pinger! If you put **/ping** into you message, i will notify everyone.

_Please note. I don't know people in this chat yet. But I remember everyone who writes something._`,
	}

	_, err = ctx.Actions.SendMessage(msg)
	if err != nil {
		return err
	}

	if u.MyChatMember.From.Username == "" {
		msg := messages.AddUsername(u.MyChatMember.Chat.Id, u.MyChatMember.From)
		_, err = ctx.Actions.SendMessage(msg)

		return err
	}

	err = h.S.RememberMember(u.MyChatMember.Chat.Id, u.MyChatMember.From)

	return err
}
