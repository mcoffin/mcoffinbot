package main

import (
	"./yirc"
	"github.com/sorcix/irc"
	"strings"
)

type CommandFunc func(*irc.Encoder, string, []string, *irc.Message) error

func (f CommandFunc) HandleCommand(enc *irc.Encoder, command string, args []string, msg *irc.Message) error {
	return f(enc, command, args, msg)
}

type Command interface {
	HandleCommand(enc *irc.Encoder, command string, args []string, msg *irc.Message) error
}

type CommandHandler struct {
	Lead     string
	Commands map[string]Command
}

func (self CommandHandler) HandleIRC(enc *irc.Encoder, msg *irc.Message, next yirc.Handler) error {
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

	return h.HandleCommand(enc, args[0], args[1:], msg)
}
