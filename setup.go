package blocker

import (

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	plugin.Register("blocker", setup)
}

func setup(c *caddy.Controller) error {
	blocker, err := fileParse(c)
	if err != nil {
		return plugin.Error("blocker", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		blocker.Next = next
		return blocker
	})

	return nil
}

func fileParse(c *caddy.Controller) (*Blocker, error) {
	var file string

	for c.Next() {

		if !c.NextArg() {
			return nil, c.ArgErr()
		}

		if file != "" {
			return nil, c.Errf("multiple files are not supported")
		}

		file = c.Val()

		if len(c.RemainingArgs()) != 0 {
			return nil, c.ArgErr()
		}
	}

	blocker, err := newBlocker(file)
	if err != nil {
		return nil, c.Err(err.Error())
	}

	return blocker, nil
}

