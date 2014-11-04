package yirc

import (
	"github.com/sorcix/irc"
)

var PingHandler Handler = HandlerFunc(handlePing)

func handlePing(enc *irc.Encoder, msg *irc.Message, next Handler) error {
	if msg.Command == irc.PING {
		var pongMessage = irc.Message{
			Command: irc.PONG,
			Params:  []string{},
		}
		return enc.Encode(&pongMessage)
	}
	return next.HandleIRC(enc, msg, nil)
}
