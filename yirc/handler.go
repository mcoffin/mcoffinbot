package yirc

import (
	"github.com/sorcix/irc"
)

type HandlerFunc func(*irc.Encoder, *irc.Message, Handler) error

func (f HandlerFunc) HandleIRC(dec *irc.Encoder, m *irc.Message, next Handler) error {
	return f(dec, m, next)
}

type Handler interface {
	HandleIRC(*irc.Encoder, *irc.Message, Handler) error
}
