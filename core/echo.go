package core

import "github.com/Janso123/go-jmap"

// The Core/echo method
type Echo struct {
	Hello string
}

func (e Echo) Name() string {
	return "Core/echo"
}

func (e Echo) Requires() []jmap.URI { return nil }

func newEcho() jmap.MethodResponse {
	return &Echo{}
}
