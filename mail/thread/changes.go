package thread

import "github.com/Janso123/go-jmap"

// See RFC8621, Section 3.2.
type Changes struct {
	jmap.Changes[Thread]
}

// ChangesResponse is the result of Thread/changes.
type ChangesResponse = jmap.ChangesResponse
