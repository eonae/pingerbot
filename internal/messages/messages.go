package messages

import (
	"fmt"
	"pingerbot/pkg/telegram"
)

func AddUsername(chatId int64, u telegram.User) telegram.OutgoingMessage {
	return telegram.OutgoingMessage{
		ChatId:    chatId,
		ParseMode: telegram.Markdown,
		Text: fmt.Sprintf(
			"I can't ping users without username mr.[%s](tg://user?id=%d). Please setup yours!",
			u.FirstName,
			u.Id,
		),
	}
}
