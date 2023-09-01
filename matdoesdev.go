package matdoesdev

import (
	"hash/fnv"
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
	Chance string `json:"string,omitempty"`
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
	chance, err := strconv.ParseFloat(m.Chance, 64)
	if err != nil {
		return false
	}

	return hash(r.URL.Path) < uint32(chance*4294967295)
}

func (m *MatchRandomPaths) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.Args(&m.Chance) {
			return d.ArgErr()
		}
	}
	return nil
}

var (
	_ caddyhttp.RequestMatcher = (*MatchRandomPaths)(nil)
	_ caddyfile.Unmarshaler    = (*MatchRandomPaths)(nil)
)
