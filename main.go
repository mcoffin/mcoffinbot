package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/mcoffin/mcoffinbot/yirc"
	"github.com/robertkrimen/otto"
	"github.com/sorcix/irc"
	"io/ioutil"
	"log"
	"regexp"
	"sync"
)

func serverLoop(cfg server, nick string, c yirc.Client) error {
	return c.ListenAndHandle(cfg.Address, nick, cfg.Channels)
}

func newRegexCommand(pattern string, script string) (*YircRegexHandler, error) {
	source, err := ioutil.ReadFile(script)
	if err != nil {
		return nil, err
	}

	vm := otto.New()
	handler := RegexHandlerFunc(func(enc *irc.Encoder, matches []string, msg *irc.Message) error {
		err := vm.Set("irc", func(call otto.FunctionCall) otto.Value {
			raw := call.Argument(0).String()
			msg := irc.ParseMessage(raw)
			err := enc.Encode(msg)
			if err == nil {
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		})
		if err != nil {
			return err
		}
		_, err = vm.Run(source)
		if err != nil {
			return err
		}

		jsArgs := []interface{}{msg.Prefix.Name, msg.Params[0], matches}
		_, err = vm.Call("handleRegex", nil, jsArgs...)
		if err != nil {
			return err
		}
		return nil
	})

	var compiled *regexp.Regexp
	compiled, err = regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &YircRegexHandler{
		Pattern:      compiled,
		RegexHandler: handler,
	}, nil
}

func newJSCommand(script string) (Command, error) {
	source, err := ioutil.ReadFile(script)
	if err != nil {
		return nil, err
	}

	vm := otto.New()
	return CommandFunc(func(enc *irc.Encoder, cmd string, args string, msg *irc.Message) error {
		err := vm.Set("irc", func(call otto.FunctionCall) otto.Value {
			raw := call.Argument(0).String()
			msg := irc.ParseMessage(raw)
			err := enc.Encode(msg)
			if err == nil {
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		})
		if err != nil {
			return err
		}
		_, err = vm.Run(source)
		if err != nil {
			return err
		}
		jsArgs := []interface{}{cmd, msg.Prefix.Name, msg.Params[0], args}
		_, err = vm.Call(cmd, nil, jsArgs...)
		if err != nil {
			return err
		}
		return nil
	}), nil
}

func main() {
	var err error

	nick := flag.String("nick", "mcoffinbot", "desired nickname")
	configFile := flag.String("config", "config.toml", "config file")

	flag.Parse()

	var cfg config
	_, err = toml.DecodeFile(*configFile, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var c = yirc.Classic()

	c.UseHandler(NewHighlightHandler())

	commandMap := map[string]Command{}
	for _, cmdCfg := range cfg.Commands {
		commandMap[cmdCfg.Name], err = newJSCommand(cmdCfg.Script)
		if err != nil {
			log.Fatal(err)
		}
	}

	var cmdHandler = CommandHandler{
		Lead:     cfg.CommandPrefix,
		Commands: commandMap,
	}
	c.UseHandler(cmdHandler)

	for _, cmdCfg := range cfg.RegexCommands {
		var rh *YircRegexHandler
		rh, err = newRegexCommand(cmdCfg.Pattern, cmdCfg.Script)
		if err != nil {
			log.Fatal(err)
		}
		c.UseHandler(rh)
	}
	urlHandler, err := newUrlCommand()
	if err != nil {
		log.Fatal(err)
	}
	c.UseHandler(urlHandler)

	var wg sync.WaitGroup
	for _, server := range cfg.Servers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := serverLoop(server, *nick, c)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	wg.Wait()
}
