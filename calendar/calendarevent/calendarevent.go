package calendarevent

import "git.sr.ht/~rockorager/go-jmap"

func init() {
	jmap.RegisterMethod("CalendarEvent/get", newGetResponse)
	jmap.RegisterMethod("CalendarEvent/changes", newChangesResponse)
	jmap.RegisterMethod("CalendarEvent/set", newSetResponse)
	jmap.RegisterMethod("CalendarEvent/copy", newCopyResponse)
	jmap.RegisterMethod("CalendarEvent/query", newQueryResponse)
	jmap.RegisterMethod("CalendarEvent/queryChanges", newQueryChangesResponse)
	jmap.RegisterMethod("CalendarEvent/parse", newParseResponse)
}
