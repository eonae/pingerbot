package telegram

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Base context that will be embedded to more specific ones
type Ctx struct {
	BotId   int64
	BotName string
	ChatId  int64
	Logger  *logrus.Entry
	api     Api
}

type Sender interface {
	SendToChat(msg MsgContent) error
}

func (ctx Ctx) SendToChat(msg MsgContent) error {
	_, err := ctx.api.SendMessage(OutgoingMessage{
		MsgContent: msg,
		ChatId:     ctx.ChatId,
	})
	return err
}

func (ctx Ctx) SendTxt(text string) error {
	return ctx.SendToChat(MsgContent{Text: text})
}

func (ctx Ctx) SendMD(text string) error {
	return ctx.SendToChat(MsgContent{Text: text, ParseMode: Markdown})
}

type MsgCtx struct {
	Ctx
	Message  IncomingMessage
	entities map[string][]string
}

func CreateMessageCtx(base Ctx, msg IncomingMessage) MsgCtx {
	runes := []rune(msg.Text)
	entities := make(map[string][]string)

	fmt.Println(msg.Entities)

	for _, e := range msg.Entities {
		key := e.Type
		value := string(runes[e.Offset : e.Offset+e.Length])

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

func (ctx MsgCtx) Reply(msg MsgContent) error {
	_, err := ctx.api.SendMessage(OutgoingMessage{
		MsgContent: msg,
		ChatId:     ctx.ChatId,
		ReplyTo:    ctx.Message.Id,
	})
	return err
}

func (ctx MsgCtx) ReplyTxt(text string) error {
	return ctx.Reply(MsgContent{Text: text})
}

func (ctx MsgCtx) ReplyMD(text string) error {
	return ctx.Reply(MsgContent{Text: text, ParseMode: Markdown})
}

type CommandCtx struct {
	MsgCtx
	Command string
}

type JoinCtx struct {
	Ctx
	Chat    Chat
	Actor   User
	Subject User
}

type LeaveCtx JoinCtx
