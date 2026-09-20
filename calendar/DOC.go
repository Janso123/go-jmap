// Package calendar will implement JMAP Calendars and JSCalendar 2.0 types for Pelton.
//
// Matrix status goal: Partial until these drafts become RFCs with assigned numbers.
//
// Pinned IETF revisions (verified on datatracker 2026-09-19):
//
//   - JMAP Calendars: draft-ietf-jmap-calendars-29
//     Published 16 September 2026.
//     https://datatracker.ietf.org/doc/draft-ietf-jmap-calendars/29/
//
//   - JSCalendar 2.0: draft-ietf-calext-jscalendarbis-20
//     Published 14 September 2026.
//     https://datatracker.ietf.org/doc/html/draft-ietf-calext-jscalendarbis-20
//
// JSContact remains RFC 9553 (see package contacts/jscontact). Calendar event data
// in this tree follows the JSCalendar bis draft above, not the original JSCalendar
// subset in RFC 9553 alone.
package calendar
