package core

import "github.com/Janso123/go-jmap"

// Echo is Core/echo: an arbitrary JSON object echoed by the server (RFC 8620 §4.1).
type Echo map[string]any

func (e Echo) Name() string { return "Core/echo" }

func (e Echo) Requires() []jmap.URI { return []jmap.URI{jmap.CoreURI} }

func newEcho() jmap.MethodResponse {
	e := Echo{}
	return &e
}
