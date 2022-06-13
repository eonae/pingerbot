package messages

import (
	"fmt"
	"pingerbot/pkg/telegram"
)

var Welcome = telegram.MsgContent{
	ParseMode: telegram.Markdown,
	Text: `Hi all!

I'm mr.Pinger! If you put **/ping** into you message, i will notify everyone.

_Please note. I don't know people in this chat yet. But I remember everyone who writes something._`,
}

func PleaseAddUsername(u telegram.User) telegram.MsgContent {
	return telegram.MsgContent{
		ParseMode: telegram.Markdown,
		Text: fmt.Sprintf(
			"I can't ping users without username mr.[%s](tg://user?id=%d). Please setup yours!",
			u.FirstName,
			u.Id,
		),
	}
}
