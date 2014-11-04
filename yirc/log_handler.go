package yirc

import (
	"github.com/sorcix/irc"
	"log"
)

var LogHandler Handler = HandlerFunc(handleLog)

func handleLog(enc *irc.Encoder, msg *irc.Message, next Handler) error {
	log.Println(msg)
	return next.HandleIRC(enc, msg, nil)
}
