package jmap

// Object is a JMAP data type that can use the generic method kit.
type Object interface {
	JMAPType() string
	Requires() []URI
}

// Creatable is an Object whose /set method permits create and update
// (RFC 8620 §5.3). Destroy-only types must not implement JMAPCreatable.
type Creatable interface {
	Object
	JMAPCreatable()
}

// MethodFlags selects which standard methods to register for an Object.
type MethodFlags uint64

const (
	MethodGet MethodFlags = 1 << iota
	MethodChanges
	MethodQuery
	MethodQueryChanges
	MethodSet
	MethodCopy
)

// RegisterObject registers MethodResponse factories for the selected standard
// methods of T (e.g. "Email/get", "Email/changes").
func RegisterObject[T Object](flags MethodFlags) {
	var zero T
	base := zero.JMAPType()
	if flags&MethodGet != 0 {
		RegisterMethod(base+"/get", func() MethodResponse { return &GetResponse[T]{} })
	}
	if flags&MethodChanges != 0 {
		RegisterMethod(base+"/changes", func() MethodResponse { return &ChangesResponse{} })
	}
	if flags&MethodQuery != 0 {
		RegisterMethod(base+"/query", func() MethodResponse { return &QueryResponse{} })
	}
	if flags&MethodQueryChanges != 0 {
		RegisterMethod(base+"/queryChanges", func() MethodResponse { return &QueryChangesResponse{} })
	}
	if flags&MethodSet != 0 {
		RegisterMethod(base+"/set", func() MethodResponse { return &SetResponse[T]{} })
	}
	if flags&MethodCopy != 0 {
		RegisterMethod(base+"/copy", func() MethodResponse { return &CopyResponse[T]{} })
	}
}
