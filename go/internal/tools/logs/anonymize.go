package logs

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

var (
	ipv4RE = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	macRE  = regexp.MustCompile(`\b(?:[0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}\b`)
	ipv6RE = regexp.MustCompile(`\b(?:[0-9A-Fa-f]{1,4}:){2,7}[0-9A-Fa-f]{1,4}\b`)
)

// anonymizer scrubs identifying details from log text with consistent
// per-value tokens so logs stay readable (same IP -> same token).
type anonymizer struct {
	mu       sync.Mutex
	hostOnce sync.Once
	hostname string
	ipTokens map[string]string
	ip6Seen  map[string]string
	macTok   map[string]string
}

func newAnonymizer() *anonymizer {
	return &anonymizer{
		ipTokens: map[string]string{},
		ip6Seen:  map[string]string{},
		macTok:   map[string]string{},
	}
}

func (a *anonymizer) token(m map[string]string, value, prefix string) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	if t, ok := m[value]; ok {
		return t
	}
	t := fmt.Sprintf("[%s%d]", prefix, len(m)+1)
	m[value] = t
	return t
}

func (a *anonymizer) hostnameFor(ctx context.Context, s *mcp.Server) string {
	a.hostOnce.Do(func() {
		data, err := s.GraphQLQuery(ctx, `query { info { os { hostname } } }`, nil)
		if err == nil {
			if info, ok := data["info"].(map[string]interface{}); ok {
				if osi, ok := info["os"].(map[string]interface{}); ok {
					a.hostname, _ = osi["hostname"].(string)
				}
			}
		}
	})
	return a.hostname
}

// scrub applies hostname/IP/MAC anonymization when the user enabled it.
func (a *anonymizer) scrub(ctx context.Context, s *mcp.Server, text string) string {
	if host := a.hostnameFor(ctx, s); host != "" {
		text = strings.ReplaceAll(text, host, "Tower")
	}

	text = ipv4RE.ReplaceAllStringFunc(text, func(ip string) string {
		if ip == "127.0.0.1" || ip == "0.0.0.0" || strings.HasPrefix(ip, "255.") {
			return ip
		}
		return a.token(a.ipTokens, ip, "ip")
	})

	text = macRE.ReplaceAllStringFunc(text, func(mac string) string {
		return a.token(a.macTok, strings.ToLower(mac), "mac")
	})

	text = ipv6RE.ReplaceAllStringFunc(text, func(s string) string {
		// Avoid false positives on time-like strings (10:24:50): require a hex
		// letter or a double colon to treat it as an IPv6 address.
		lower := strings.ToLower(s)
		if !strings.Contains(s, "::") && !strings.ContainsAny(lower, "abcdef") {
			return s
		}
		return a.token(a.ip6Seen, s, "ip6")
	})

	return text
}
