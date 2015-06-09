package main

import (
	"fmt"
	"github.com/sorcix/irc"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
)

func handleUrlIRC(enc *irc.Encoder, matches []string, msg *irc.Message) error {
	for _, url := range matches {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("URL handler: GET '%s' returned status code: %d", url, resp.StatusCode)
		}

		node, err := html.Parse(resp.Body)
		if err != nil {
			return err
		}

		var search func(*html.Node) (string, error)
		search = func(n *html.Node) (string, error) {
			if n.Type == html.ElementNode && n.Data == "title" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						return c.Data, nil
					}
				}
			}
			for r := n.FirstChild; r != nil; r = r.NextSibling {
				ret, e := search(r)
				if e != nil {
					return "", e
				} else if len(ret) > 0 {
					return ret, nil
				}
			}
			return "", nil
		}
		title, err := search(node)
		if err != nil {
			return err
		}

		outMsg := irc.Message{
			Command:  irc.PRIVMSG,
			Params:   []string{msg.Params[0]},
			Trailing: fmt.Sprintf("%s -- %s", url, title),
		}
		err = enc.Encode(&outMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

func newUrlCommand() (*YircRegexHandler, error) {
	compiled, err := regexp.Compile("(https?|ftp)://[^\\s/$.?#].[^\\s]*")
	if err != nil {
		return nil, err
	}

	return &YircRegexHandler{
		Pattern:      compiled,
		RegexHandler: RegexHandlerFunc(handleUrlIRC),
	}, nil
}
