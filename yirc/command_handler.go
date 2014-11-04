package yirc

import (
	"github.com/sorcix/irc"
	"strings"
)

type CommandHandler struct {
	Lead     string
	Commands map[string]Handler
}

func (self CommandHandler) HandleIRC(enc *irc.Encoder, msg *irc.Message, next Handler) error {
	if msg.Command != irc.PRIVMSG {
		return next.HandleIRC(enc, msg, nil)
	}

	var message = msg.Trailing
	if !strings.HasPrefix(message, self.Lead) {
		return next.HandleIRC(enc, msg, nil)
	} else {
		message = strings.TrimPrefix(message, self.Lead)
	}

	args := strings.Split(message, " ")

	var h = self.Commands[args[0]]
	if h == nil {
		return next.HandleIRC(enc, msg, nil)
	}

	msg.Trailing = message
	return h.HandleIRC(enc, msg, next)
}
