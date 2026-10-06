package core

import "github.com/Janso123/go-jmap"

func init() {
	jmap.RegisterMethod("Core/echo", newEcho)
}

// Core is the urn:ietf:params:jmap:core capability. Package jmap registers it.
type Core = jmap.Core
