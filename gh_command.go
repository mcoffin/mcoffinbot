package main

import (
	"fmt"
	"github.com/sorcix/irc"
	"net/http"
)

func ghCommandHandler(enc *irc.Encoder, command string, args string, msg *irc.Message) error {
	url := "http://github.com/" + args
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		outMsg := irc.Message{
			Command:  irc.PRIVMSG,
			Params:   []string{msg.Params[0]},
			Trailing: url,
		}
		err = enc.Encode(&outMsg)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("!gh: GET '%s' returned status code: %d", url, resp.StatusCode)
	}

	return nil
}
