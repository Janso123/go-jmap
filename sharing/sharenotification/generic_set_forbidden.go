//go:build destroyonlyset

package sharenotification

import "github.com/Janso123/go-jmap"

// Compile probe: ShareNotification/set is destroy-only (RFC 9670 §3.3).
var genericSetForbidden = jmap.Set[ShareNotification]{}
