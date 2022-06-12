package handlers

import (
	"pingerbot/internal/messages"
	"pingerbot/internal/state"
	"pingerbot/pkg/telegram"

	"github.com/sirupsen/logrus"
)

func rememberIfHasUsername(
	groupId int64,
	user telegram.User,
	s state.State,
	logger *logrus.Entry,
	sender telegram.Sender,
) error {
	if user.Username == "" {
		logger.Debugf("Can't add user %s - no username!", user.FirstName)
		return sender.SendToChat(messages.PleaseAddUsername(user))
	}

	logger.Infof("Remembering user @%s", user.Username)
	return s.RememberMember(groupId, user.Username)
}
