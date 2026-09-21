package vacationresponse

import "github.com/Janso123/go-jmap"

// Get vacation response details
// https://www.rfc-editor.org/rfc/rfc8621.html#section-8.1
type Get struct {
	jmap.Get[VacationResponse]
}

// GetResponse is the result of VacationResponse/get.
type GetResponse = jmap.GetResponse[VacationResponse]
