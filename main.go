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

type channel struct {
	Name string `toml:"name"`
}

type config struct {
	Channels []channel `toml:"channel"`
}

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

	var c = yirc.Client{}

	// Printing middleware handler
	c.UseHandler(yirc.HandlerFunc(func(enc *irc.Encoder, msg *irc.Message, next yirc.Handler) error {
		fmt.Println(msg)
		next.HandleIRC(enc, msg, nil)
		return nil
	}))

	c.UseHandler(yirc.PingHandler)

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

	var cmdHandler = yirc.CommandHandler{
		Lead:     "!",
		Commands: map[string]yirc.Handler{"quote": yirc.HandlerFunc(quoteHandler)},
	}
	c.UseHandler(cmdHandler)

	var channels = make([]string, 0, len(cfg.Channels))
	for _, ch := range cfg.Channels {
		channels = append(channels, ch.Name)
	}

	err = c.ListenAndHandle(*addr, *nick, channels)
	if err != nil {
		log.Fatal(err)
	}
}
