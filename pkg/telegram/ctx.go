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

// Message context

type MsgCtx struct {
	Ctx
	Message  IncomingMessage
	entities map[string][]string
}

func (base Ctx) ToMsgContext(msg IncomingMessage) MsgCtx {
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

func (ctx MsgCtx) entitiesOfType(t string) []string {
	mentions, ok := ctx.entities[t]
	if ok {
		return mentions
	}
	return make([]string, 0)
}

// Command context

type CommandCtx struct {
	MsgCtx
	Command string
}

func (base MsgCtx) ToCommandCtx(cmd string) CommandCtx {
	return CommandCtx{
		MsgCtx:  base,
		Command: cmd,
	}
}

//

type JoinCtx struct {
	Ctx
	Chat    Chat
	Actor   User
	Subject User
}

func ctxFromJoinEvent(ctx Ctx, data JoinLeave) JoinCtx {
	return JoinCtx{
		Ctx:     ctx,
		Chat:    data.Chat,
		Actor:   data.From,
		Subject: data.NewChatMember.User,
	}
}

func ctxFromLeaveEvent(ctx Ctx, data JoinLeave) LeaveCtx {
	return LeaveCtx(ctxFromJoinEvent(ctx, data))
}

func joinCtxFromMessage(ctx Ctx, data IncomingMessage) JoinCtx {
	return JoinCtx{
		Ctx:     ctx,
		Chat:    data.Chat,
		Actor:   data.From,
		Subject: *data.NewMember,
	}
}

func leaveCtxFromMessage(ctx Ctx, data IncomingMessage) LeaveCtx {
	return LeaveCtx{
		Ctx:     ctx,
		Chat:    data.Chat,
		Actor:   data.From,
		Subject: *data.LeftMember,
	}
}

type LeaveCtx JoinCtx
