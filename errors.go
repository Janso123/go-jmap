package jmap

// Request-level error type URIs (RFC 8620 §3.6.1).
const (
	ErrTypeUnknownCapability = "urn:ietf:params:jmap:error:unknownCapability"
	ErrTypeNotJSON           = "urn:ietf:params:jmap:error:notJSON"
	ErrTypeNotRequest        = "urn:ietf:params:jmap:error:notRequest"
	ErrTypeLimit             = "urn:ietf:params:jmap:error:limit"
)

// A RequestError occurs when there is an error with the HTTP request.
type RequestError struct {
	// The type of request error, eg "urn:ietf:params:jmap:error:limit"
	Type string `json:"type"`

	// The HTTP status code of the response
	Status int `json:"status"`

	// A short, human-readable summary of the problem type
	Title string `json:"title,omitzero"`

	// The description of the error
	Detail string `json:"detail"`

	// If the error is of type ErrTypeLimit, Limit will contain the name of the
	// limit the request would have exceeded
	Limit *string `json:"limit,omitzero"`

	// Server-assigned request id when present (problem+json / WebSocket)
	RequestID string `json:"requestId,omitzero"`
}

func (e *RequestError) Error() string {
	if e.Limit != nil {
		return e.Detail + ": " + *e.Limit
	}
	if e.Detail != "" {
		return e.Detail
	}
	return e.Type
}

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
	Type string `json:"type,omitempty"`

	// Description is available on some method errors (notably,
	// invalidArguments)
	Description *string `json:"description,omitempty"`
}

func (m *MethodError) Error() string {
	if m.Description != nil {
		return m.Type + ": " + *m.Description
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
	SetErrInvalidScript     SetErrorType = "invalidScript"
	SetErrScriptIsActive    SetErrorType = "scriptIsActive"
)

// A SetError is returned in set calls for individual record changes
type SetError struct {
	// The type of SetError
	Type string `json:"type,omitempty"`

	// A description of the error to help with debugging that includes an
	// explanation of what the problem was. This is a non-localised string
	// and is not intended to be shown directly to end users.
	Description *string `json:"description,omitempty"`

	// Properties is available on InvalidProperties SetErrors and lists the
	// individual properties were
	Properties *[]string `json:"properties,omitempty"`
}

func (s *SetError) Error() string {
	if s.Description != nil {
		return s.Type + ": " + *s.Description
	}
	return s.Type
}
