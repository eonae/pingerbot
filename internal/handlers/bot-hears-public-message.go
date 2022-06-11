package handlers

import (
	"fmt"
	"pingerbot/internal/state"
	"pingerbot/pkg/telegram"
	"strings"
)

type BotHearsPublicMessage struct {
	S state.State
}

func (BotHearsPublicMessage) Name() string {
	return "BotHearsPublicMessage"
}

func (BotHearsPublicMessage) Match(u telegram.Update) bool {
	return u.Message != nil && u.Message.Chat.Type != "private"
}

func (h BotHearsPublicMessage) Handle(u telegram.Update, ctx telegram.Ctx) error {
	if u.Message.From.Id == ctx.BotId {
		ctx.Logger.Debug("Skipping message from self")
		return nil
	}

	groupId := u.Message.Chat.Id

	// We don't work with users that don't have username because
	// there is no way to mention them.
	if u.Message.From.Username == "" {
		msg := telegram.SendMessage{
			ChatId:    u.Message.Chat.Id,
			ParseMode: telegram.Markdown,
			Text:      fmt.Sprintf("I can't ping users without username mr.%s. Please setup yours!", u.Message.From.FirstName),
		}

		_, err := ctx.Actions.SendMessage(msg)

		return err
	}

	ctx.Logger.Debugf("Remembering user @%s", u.Message.From.Username)
	err := h.S.RememberMember(groupId, u.Message.From)
	if err != nil {
		return err
	}

	ctx.Logger.Debug("Looping through entities...")
	for _, e := range u.Message.Entities {
		switch e.Type {
		case "bot_command":
			cmd := u.Message.Text[e.Offset : e.Offset+e.Length]
			if cmd == "/ping" {
				members, err := h.S.GetKnownMembers(groupId)
				if err != nil {
					return err
				}

				mentions := make([]string, 0)

				for _, username := range members {
					mentions = append(mentions, "@"+username)
				}

				_, err = ctx.Actions.SendMessage(telegram.SendMessage{
					ChatId:  groupId,
					ReplyTo: u.Message.MessageId,
					Text:    strings.Join(mentions, " "),
				})
				return err
			}
		case "mention":
			username := u.Message.Text[e.Offset : e.Offset+e.Length]
			if username == ctx.BotName {
				_, err := ctx.Actions.SendMessage(telegram.SendMessage{
					ChatId:  u.Message.Chat.Id,
					Text:    "You can talk to me in private if you want.",
					ReplyTo: u.Message.MessageId,
				})
				return err
			}
		}
	}

	return nil
}
