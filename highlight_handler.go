package main

import (
	"./yirc"
	"fmt"
	"github.com/sorcix/irc"
	"strings"
	"sync"
	"time"
)

type HighlightHandler struct {
	mutex    sync.Mutex
	tracking map[string][]*irc.Message
}

func NewHighlightHandler() *HighlightHandler {
	return &HighlightHandler{
		tracking: map[string][]*irc.Message{},
	}
}

func (self *HighlightHandler) stopTracking(nick string) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.tracking[nick] = nil
}

func (self *HighlightHandler) HandleIRC(enc *irc.Encoder, msg *irc.Message, next yirc.Handler) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if msg.Command == irc.QUIT || msg.Command == irc.PART {
		self.tracking[msg.Prefix.Name] = []*irc.Message{}
		go func() {
			time.Sleep(60 * time.Second)
			self.stopTracking(msg.Prefix.Name)
		}()
	} else if msg.Command == irc.PRIVMSG {
		for name, highlights := range self.tracking {
			if strings.Contains(msg.Trailing, name) {
				highlights = append(highlights, msg)
				self.tracking[name] = highlights
			}
		}
	} else if msg.Command == irc.JOIN {
		if self.tracking[msg.Prefix.Name] != nil {
			for _, hMsg := range self.tracking[msg.Prefix.Name] {
				outMsg := irc.Message{
					Command:  irc.PRIVMSG,
					Params:   []string{msg.Prefix.Name},
					Trailing: fmt.Sprintf("[%s]<%s> %s", hMsg.Params[0], hMsg.Prefix.Name, hMsg.Trailing),
				}
				err := enc.Encode(&outMsg)
				if err != nil {
					return err
				}
			}
		}
	}
	return next.HandleIRC(enc, msg, nil)
}
