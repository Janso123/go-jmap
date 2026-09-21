package core

import (
	"context"
	"fmt"
	"net"
	"strings"
)

// lookupSRV is the DNS SRV lookup used by Discover. Tests may replace it.
var lookupSRV = net.DefaultResolver.LookupSRV

// Discover the Session Endpoint of a domain.
// It looks up _jmap._tcp SRV records; if lookup fails or returns no records,
// it falls back to https://<domain>/.well-known/jmap.
func Discover(ctx context.Context, domain string) (string, error) {
	if err := ctx.Err(); err != nil {
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
