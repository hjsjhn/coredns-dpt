package encoding20

import (
	"strings"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("encoding20", setup) }

func setup(c *caddy.Controller) error {
	sid, err := sidParse(c)
	if err != nil {
		return plugin.Error("0x20 encoding", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Sid{Next: next, Data: sid}
	})

	return nil
}

func sidParse(c *caddy.Controller) (string, error) {
	// Use hostname as the default
	sid := ""
	i := 0
	for c.Next() {
		if i > 0 {
			return sid, plugin.ErrOnce
		}
		i++
		args := c.RemainingArgs()
		if len(args) > 0 {
			sid = strings.Join(args, " ")
		}
	}
	return sid, nil
}
