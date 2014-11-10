package yirc

import (
	"fmt"
	"github.com/sorcix/irc"
	"log"
)

type middleware struct {
	Handler Handler
	Next    *middleware
}

func voidMiddlewareHandler(enc *irc.Encoder, msg *irc.Message, next Handler) error {
	return nil
}

func (self middleware) HandleIRC(enc *irc.Encoder, msg *irc.Message, next Handler) error {
	if self.Next == nil {
		return self.Handler.HandleIRC(enc, msg, HandlerFunc(voidMiddlewareHandler))
	} else {
		return self.Handler.HandleIRC(enc, msg, self.Next)
	}
}

type Client struct {
	*irc.Conn
	Handlers []Handler
	nick     string
}

func Classic() Client {
	return Client{
		Handlers: []Handler{LogHandler, PingHandler},
	}
}

func buildMiddleware(handlers []Handler) middleware {
	if len(handlers) == 1 {
		return middleware{
			Handler: handlers[0],
		}
	} else {
		var next = buildMiddleware(handlers[1:])
		return middleware{
			Handler: handlers[0],
			Next:    &next,
		}
	}
}

func (self *Client) UseHandler(h Handler) {
	self.Handlers = append(self.Handlers, h)
}

func (self *Client) ListenAndHandle(addr string, nick string, channels []string) error {
	var err error

	// Open connection
	self.Conn, err = irc.Dial(addr)
	if err != nil {
		return err
	}
	defer self.Close()

	// Sign up with a nickname
	var nickMessage = irc.Message{
		Command: irc.NICK,
		Params:  []string{nick},
	}
	err = self.Encode(&nickMessage)
	if err == nil {
		self.nick = nick
	} else {
		return err
	}

	// Tell the server about us
	var userMessage = irc.Message{
		Command: irc.USER,
		Params:  []string{nick, "0", "*", fmt.Sprintf(":%s", nick)},
	}
	err = self.Encode(&userMessage)
	if err != nil {
		return err
	}

	// Join channels
	var joinMessage = irc.Message{
		Command: irc.JOIN,
	}
	for _, c := range channels {
		joinMessage.Params = []string{fmt.Sprintf(":%s", c)}
		err = self.Encode(&joinMessage)
		if err != nil {
			return err
		}
	}

	// Convert handlers to middleware stack
	var middlewareStack = buildMiddleware(self.Handlers)

	// Receive messages and respond to them
	for err == nil {
		msg, err := self.Decode()
		// Passing nil because the middleware stack already handles it
		if err == nil {
			err = middlewareStack.HandleIRC(&self.Encoder, msg, nil)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return err
}
