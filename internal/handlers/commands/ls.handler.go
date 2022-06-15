package commands

import (
	"pingerbot/internal/state"
	tg "pingerbot/pkg/telegram"
)

type LsCommandHandler struct {
	S state.State
}

func (h LsCommandHandler) Handle(ctx tg.CommandCtx) error {
	tags := ctx.Tags()

	// Получить из БД список и вывести в виде сообщения. Если приведены тэги, то отфильтровать по ним.

	return nil
}
