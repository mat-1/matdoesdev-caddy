package matdoesdev

import (
	"crypto/sha1"
	"net/http"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(MatchRandomPaths{})
}

type MatchRandomPaths struct {
	Chance float64 `json:"chance,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (MatchRandomPaths) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.matchers.random_paths",
		New: func() caddy.Module { return new(MatchRandomPaths) },
	}
}

func hash(s string) uint32 {
	h := sha1.New()
	h.Write([]byte(s))
	buf := h.Sum(nil)
	return uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
}

func (m MatchRandomPaths) Match(r *http.Request) bool {
	return (float64(hash(r.URL.Path)))/4294967295.0 < m.Chance
}

func (m *MatchRandomPaths) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			chance, err := strconv.ParseFloat(d.Val(), 64)
			if err != nil {
				return err
			}
			m.Chance = chance
		}
	}
	return nil
}

var (
	_ caddyhttp.RequestMatcher = (*MatchRandomPaths)(nil)
	_ caddyfile.Unmarshaler    = (*MatchRandomPaths)(nil)
)
