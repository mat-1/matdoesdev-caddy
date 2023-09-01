package matdoesdev

import (
	"hash/fnv"
	"net/http"

	"github.com/caddyserver/caddy/caddyfile"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(MatchRandomPaths{})
}

type MatchRandomPaths struct {
	Chance float32 `json:"number,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (MatchRandomPaths) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.matchers.random_paths",
		New: func() caddy.Module { return new(MatchRandomPaths) },
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (m MatchRandomPaths) Match(r *http.Request) bool {
	return hash(r.URL.Path) < uint32(m.Chance*4294967295)
}

func (m *MatchRandomPaths) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.Next() {
			if !d.Args(&m.Chance) {
				return d.ArgErr()
			}
		}
		if d.NextBlock(0) {
			return d.Err("malformed random_paths matcher: blocks are not supported")
		}
	}
	return nil
}

var (
	_ caddyhttp.RequestMatcher = (*MatchRandomPaths)(nil)
)
