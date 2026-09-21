package thread

import "github.com/Janso123/go-jmap"

// See RFC8621, Section 3.1.
type Get struct {
	jmap.Get[Thread]
}

// GetResponse is the result of Thread/get.
type GetResponse = jmap.GetResponse[Thread]
