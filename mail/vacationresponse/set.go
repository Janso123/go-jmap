package vacationresponse

import "github.com/Janso123/go-jmap"

// Create, update & modify vacation responses
// https://www.rfc-editor.org/rfc/rfc8621.html#section-8.2
type Set struct {
	jmap.Set[VacationResponse]
}

// SetResponse is the result of VacationResponse/set.
type SetResponse = jmap.SetResponse[VacationResponse]
