package main

import (
	"./yirc"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sorcix/irc"
	"log"
	"strings"
)

var quotes = map[string]string{}

func quoteHandler(enc *irc.Encoder, msg *irc.Message, next yirc.Handler) error {
	var quot = strings.TrimPrefix(msg.Trailing, "quote ")
	quotes[msg.Name] = quot
	return nil
}

func main() {
	var err error

	nick := flag.String("nick", "mcoffinbot", "desired nickname")
	addr := flag.String("addr", "irc.freenode.net:6667", "server to connect")
	configFile := flag.String("config", "config.toml", "config file")

	flag.Parse()

	var cfg config
	_, err = toml.DecodeFile(*configFile, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var c = yirc.Classic()

	c.UseHandler(yirc.HandlerFunc(func(enc *irc.Encoder, msg *irc.Message, next yirc.Handler) error {
		if msg.Command == irc.JOIN {
			var quote = quotes[msg.Name]
			if quote != "" {
				var greeting = fmt.Sprintf("<%s> %s", msg.Name, quote)
				var greetingMessage = irc.Message{
					Command:  irc.PRIVMSG,
					Params:   []string{msg.Params[0]},
					Trailing: greeting,
				}
				return enc.Encode(&greetingMessage)
			}
			return nil
		} else {
			return next.HandleIRC(enc, msg, nil)
		}
	}))

	var cmdHandler = CommandHandler{
		Lead:     "!",
		Commands: map[string]yirc.Handler{"quote": yirc.HandlerFunc(quoteHandler)},
	}
	c.UseHandler(cmdHandler)

	err = c.ListenAndHandle(*addr, *nick, cfg.Channels)
	if err != nil {
		log.Fatal(err)
	}
}
