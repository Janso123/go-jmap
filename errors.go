package jmap

import "encoding/json/jsontext"

// Request-level error type URIs (RFC 8620 §3.6.1).
const (
	ErrTypeUnknownCapability = "urn:ietf:params:jmap:error:unknownCapability"
	ErrTypeNotJSON           = "urn:ietf:params:jmap:error:notJSON"
	ErrTypeNotRequest        = "urn:ietf:params:jmap:error:notRequest"
	ErrTypeLimit             = "urn:ietf:params:jmap:error:limit"
)

// A RequestError occurs when there is an error with the HTTP request.
// It is also the RFC 7807 problem+json object returned for HTTP failures.
type RequestError struct {
	// The type of request error, eg "urn:ietf:params:jmap:error:limit"
	Type string `json:"type"`

	// The HTTP status code of the response
	Status int `json:"status"`

	// A short, human-readable summary of the problem type
	Title string `json:"title,omitzero"`

	// Detail is a human-readable explanation of this occurrence.
	// JSON null stays null; an unset value omits the key.
	Detail Optional[string] `json:"detail,omitzero"`

	// Instance identifies this specific occurrence (RFC 7807).
	// JSON null stays null; an unset value omits the key.
	Instance Optional[string] `json:"instance,omitzero"`

	// If the error is of type ErrTypeLimit, Limit will contain the name of the
	// limit the request would have exceeded
	Limit Optional[string] `json:"limit,omitzero"`

	// Server-assigned request id when present (problem+json / WebSocket)
	RequestID string `json:"requestId,omitzero"`

	// Extra holds problem+json extension members, including vendor properties.
	Extra map[string]jsontext.Value `json:",embed"`
}

func (e *RequestError) Error() string {
	detail, _ := e.Detail.Value()
	if l, ok := e.Limit.Value(); ok {
		return detail + ": " + l
	}
	if detail != "" {
		return detail
	}
	return e.Type
}

// HTTPError is returned for non-JSON HTTP failures so callers can inspect Status.
type HTTPError struct {
	Status      int
	StatusText  string
	Body        string
	ContentType string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return "HTTP " + e.StatusText
	}
	return "HTTP " + e.StatusText + ": " + e.Body
}

func (e *HTTPError) StatusCode() int { return e.Status }

// MethodErrorType is a method-level error type string (RFC 8620 §3.6.2 + mail).
// MethodError.Type stays string for extension tolerance; use these for compare.
type MethodErrorType string

const (
	MethodErrServerUnavailable               MethodErrorType = "serverUnavailable"
	MethodErrServerFail                      MethodErrorType = "serverFail"
	MethodErrServerPartialFail               MethodErrorType = "serverPartialFail"
	MethodErrUnknownMethod                   MethodErrorType = "unknownMethod"
	MethodErrInvalidArguments                MethodErrorType = "invalidArguments"
	MethodErrInvalidResultReference          MethodErrorType = "invalidResultReference"
	MethodErrForbidden                       MethodErrorType = "forbidden"
	MethodErrAccountNotFound                 MethodErrorType = "accountNotFound"
	MethodErrAccountNotSupportedByMethod     MethodErrorType = "accountNotSupportedByMethod"
	MethodErrAccountReadOnly                 MethodErrorType = "accountReadOnly"
	MethodErrRequestTooLarge                 MethodErrorType = "requestTooLarge"
	MethodErrCannotCalculateChanges          MethodErrorType = "cannotCalculateChanges"
	MethodErrStateMismatch                   MethodErrorType = "stateMismatch"
	MethodErrAlreadyExists                   MethodErrorType = "alreadyExists"
	MethodErrFromAccountNotFound             MethodErrorType = "fromAccountNotFound"
	MethodErrFromAccountNotSupportedByMethod MethodErrorType = "fromAccountNotSupportedByMethod"
	MethodErrAnchorNotFound                  MethodErrorType = "anchorNotFound"
	MethodErrUnsupportedSort                 MethodErrorType = "unsupportedSort"
	MethodErrUnsupportedFilter               MethodErrorType = "unsupportedFilter"
	MethodErrTooManyChanges                  MethodErrorType = "tooManyChanges"
)

// ErrStateMismatch is a sentinel for errors.Is against stateMismatch method errors.
var ErrStateMismatch = &MethodError{Type: string(MethodErrStateMismatch)}

// A MethodError is returned when an error occurred while the server was
// processing a method. Instead of the Response of that method, a MethodError
// invocation will be in it's place
type MethodError struct {
	// The type of error that occurred. Always present
	Type string `json:"type,omitzero"`

	// Description is available on some method errors (notably,
	// invalidArguments)
	Description Optional[string] `json:"description,omitzero"`

	// Extra holds RFC 8620 §3.6.2 extension properties.
	Extra map[string]jsontext.Value `json:",embed"`
}

func (m *MethodError) Error() string {
	if d, ok := m.Description.Value(); ok {
		return m.Type + ": " + d
	}
	return m.Type
}

// Is reports whether target is a *MethodError with the same Type.
func (m *MethodError) Is(target error) bool {
	t, ok := target.(*MethodError)
	return ok && m.Type == t.Type
}

func newMethodError() MethodResponse { return &MethodError{} }

// SetErrorType is a set-error type string (Core + Mail + Submission + Sieve).
// SetError.Type stays string for extension tolerance; use these for compare.
type SetErrorType string

const (
	SetErrForbidden         SetErrorType = "forbidden"
	SetErrOverQuota         SetErrorType = "overQuota"
	SetErrTooLarge          SetErrorType = "tooLarge"
	SetErrRateLimit         SetErrorType = "rateLimit"
	SetErrNotFound          SetErrorType = "notFound"
	SetErrInvalidPatch      SetErrorType = "invalidPatch"
	SetErrWillDestroy       SetErrorType = "willDestroy"
	SetErrInvalidProperties SetErrorType = "invalidProperties"
	SetErrSingleton         SetErrorType = "singleton"
	SetErrMailboxHasChild   SetErrorType = "mailboxHasChild"
	SetErrMailboxHasEmail   SetErrorType = "mailboxHasEmail"
	SetErrBlobNotFound      SetErrorType = "blobNotFound"
	SetErrTooManyKeywords   SetErrorType = "tooManyKeywords"
	SetErrTooManyMailboxes  SetErrorType = "tooManyMailboxes"
	SetErrForbiddenFrom     SetErrorType = "forbiddenFrom"
	SetErrInvalidEmail      SetErrorType = "invalidEmail"
	SetErrTooManyRecipients SetErrorType = "tooManyRecipients"
	SetErrNoRecipients      SetErrorType = "noRecipients"
	SetErrInvalidRecipients SetErrorType = "invalidRecipients"
	SetErrForbiddenMailFrom SetErrorType = "forbiddenMailFrom"
	SetErrForbiddenToSend   SetErrorType = "forbiddenToSend"
	SetErrCannotUnsend      SetErrorType = "cannotUnsend"
	SetErrAlreadyExists     SetErrorType = "alreadyExists"
	SetErrInvalidSieve      SetErrorType = "invalidSieve"
	SetErrSieveIsActive     SetErrorType = "sieveIsActive"
	SetErrMDNAlreadySent    SetErrorType = "mdnAlreadySent"
	SetErrUnknownDataType   SetErrorType = "unknownDataType"
)

// A SetError is returned in set calls for individual record changes
type SetError struct {
	Type              string             `json:"type,omitzero"`
	Description       Optional[string]   `json:"description,omitzero"`
	Properties        Optional[[]string] `json:"properties,omitzero"`
	ExistingID        ID                 `json:"existingId,omitzero"`
	NotFound          []ID               `json:"notFound,omitzero"`
	MaxSize           *UnsignedInt       `json:"maxSize,omitzero"`
	MaxRecipients     *UnsignedInt       `json:"maxRecipients,omitzero"`
	InvalidRecipients []string           `json:"invalidRecipients,omitzero"`

	// Extra holds RFC 8620 §5.3 extension properties ("Other properties MAY").
	Extra map[string]jsontext.Value `json:",embed"`
}

func (s *SetError) Error() string {
	if d, ok := s.Description.Value(); ok {
		return s.Type + ": " + d
	}
	return s.Type
}
