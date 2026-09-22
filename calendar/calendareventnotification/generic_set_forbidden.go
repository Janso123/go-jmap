//go:build destroyonlyset

package calendareventnotification

import "github.com/Janso123/go-jmap"

// Compile probe: CalendarEventNotification/set is destroy-only (draft-calendars-29 §7.3).
var genericSetForbidden = jmap.Set[CalendarEventNotification]{}
