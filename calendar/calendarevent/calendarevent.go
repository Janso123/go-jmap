package calendarevent

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
)

func init() {
	jmap.RegisterObject[CalendarEvent](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodQuery |
			jmap.MethodQueryChanges |
			jmap.MethodSet |
			jmap.MethodCopy,
	)
	// CalendarEvent/changes includes updatedProperties; override kit factory.
	jmap.RegisterMethod("CalendarEvent/changes", func() jmap.MethodResponse { return &ChangesResponse{} })
	// Manual method (not part of the standard kit).
	jmap.RegisterMethod("CalendarEvent/parse", newParseResponse)
}

func (CalendarEvent) JMAPType() string { return "CalendarEvent" }

func (CalendarEvent) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }
