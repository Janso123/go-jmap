package core

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
)

// lookupSRV is the DNS SRV lookup used by Discover. Tests may replace it.
var lookupSRV = net.DefaultResolver.LookupSRV

func validateDiscoverDomain(domain string) error {
	if domain == "" || strings.ContainsAny(domain, "/\\?#@[]") || strings.Contains(domain, "://") || strings.Contains(domain, ":") {
		return fmt.Errorf("jmap: discover: domain must be a hostname")
	}
	u, err := url.Parse("https://" + domain)
	if err != nil || u.User != nil || u.Host != domain || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("jmap: discover: domain must be a hostname")
	}
	if u.Hostname() == "" {
		return fmt.Errorf("jmap: discover: domain must be a hostname")
	}
	return nil
}

// Discover the Session Endpoint of a domain.
// It looks up _jmap._tcp SRV records; if lookup fails or returns no records,
// it falls back to https://<domain>/.well-known/jmap.
func Discover(ctx context.Context, domain string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if err := validateDiscoverDomain(domain); err != nil {
		return "", err
	}

	_, srvs, err := lookupSRV(ctx, "jmap", "tcp", domain)
	if err == nil && len(srvs) > 0 {
		srv := srvs[0]
		endpoint := strings.Builder{}
		endpoint.WriteString("https://")
		endpoint.WriteString(strings.TrimSuffix(srv.Target, "."))
		if srv.Port > 0 {
			endpoint.WriteString(fmt.Sprintf(":%d", srv.Port))
		}
		endpoint.WriteString("/.well-known/jmap")
		return endpoint.String(), nil
	}

	if err := ctx.Err(); err != nil {
		return "", err
	}

	return "https://" + domain + "/.well-known/jmap", nil
}
