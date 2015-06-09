package main

import (
	"github.com/mcoffin/mcoffinbot/yirc"
	"github.com/sorcix/irc"
	"regexp"
)

type RegexHandlerFunc func(*irc.Encoder, []string, *irc.Message) error

type RegexHandler interface {
	HandleRegex(enc *irc.Encoder, matches []string, msg *irc.Message) error
}

func (f RegexHandlerFunc) HandleRegex(enc *irc.Encoder, matches []string, msg *irc.Message) error {
	return f(enc, matches, msg)
}

type YircRegexHandler struct {
	Pattern *regexp.Regexp
	RegexHandler
}

func (self YircRegexHandler) HandleIRC(enc *irc.Encoder, msg *irc.Message, next yirc.Handler) error {
	if msg.Command != irc.PRIVMSG {
		return next.HandleIRC(enc, msg, nil)
	}

	message := msg.Trailing
	matches := self.Pattern.FindAllString(message, -1)
	if matches == nil {
		return next.HandleIRC(enc, msg, nil)
	}

	return self.HandleRegex(enc, matches, msg)
}
