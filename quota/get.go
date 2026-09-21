package quota

import "github.com/Janso123/go-jmap"

// Get quota details.
type Get struct {
	jmap.Get[Quota]
}

// GetResponse is the result of Quota/get.
type GetResponse = jmap.GetResponse[Quota]
