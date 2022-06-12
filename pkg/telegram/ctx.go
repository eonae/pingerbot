package telegram

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Ctx struct {
	BotId   int64
	BotName string
	ChatId  int64
	Logger  *logrus.Entry
	api     Api
}

type Sender interface {
	SendToChat(msg OutgoingMessage) error
}

func (ctx Ctx) SendToChat(msg OutgoingMessage) error {
	msg.ChatId = ctx.ChatId
	_, err := ctx.api.SendMessage(msg)
	return err
}

type MsgCtx struct {
	Ctx
	Message  IncomingMessage
	entities map[string][]string
}

type CommandCtx struct {
	MsgCtx
	Command string
}

func CreateMessageCtx(base Ctx, msg IncomingMessage) MsgCtx {
	entities := make(map[string][]string)

	fmt.Println(msg.Entities)

	for _, e := range msg.Entities {
		key := e.Type
		value := msg.Text[e.Offset : e.Offset+e.Length]

		_, ok := entities[key]
		if ok {
			entities[key] = append(entities[key], value)
		} else {
			entities[key] = []string{value}
		}
	}

	return MsgCtx{base, msg, entities}
}

func (ctx MsgCtx) entitiesOfType(t string) []string {
	mentions, ok := ctx.entities[t]
	if ok {
		return mentions
	}
	return make([]string, 0)
}

func (ctx MsgCtx) Commands() []string {
	return ctx.entitiesOfType("bot_command")
}

func (ctx MsgCtx) Mentions() []string {
	return ctx.entitiesOfType("mention")
}

func (ctx MsgCtx) Tags() []string {
	return ctx.entitiesOfType("hashtag")
}

func (ctx MsgCtx) Reply(msg OutgoingMessage) error {
	msg.ChatId = ctx.ChatId
	msg.ReplyTo = ctx.Message.Id
	_, err := ctx.api.SendMessage(msg)
	return err
}

type JoinCtx struct {
	Ctx
	Chat    Chat
	Actor   User
	Subject User
}

type LeaveCtx JoinCtx
