package jmap

import (
	"fmt"
	"strings"
)

// checkRefs rejects "#foo" when foo (or a typed #foo) is also set (RFC 8620 §3.7 invalidArguments).
func checkRefs(refs map[string]*ResultReference, set map[string]bool) error {
	for k := range refs {
		if !strings.HasPrefix(k, "#") {
			return fmt.Errorf("jmap: Refs key %q must start with '#'", k)
		}
		if set[k[1:]] {
			return fmt.Errorf("jmap: argument %q given both directly and as a result reference", k[1:])
		}
	}
	return nil
}

// Get fetches objects by id (RFC 8620 §5.1).
//
// T must be a non-pointer struct type. Name() calls JMAPType on the zero T;
// a pointer or interface instantiation panics (Go spec, Calls).
type Get[T Object] struct {
	Account             ID                          `json:"accountId,omitzero"`
	IDs                 Optional[[]ID]              `json:"ids,omitzero"`
	Properties          Optional[[]string]          `json:"properties,omitzero"`
	ReferenceIDs        *ResultReference            `json:"#ids,omitzero"`
	ReferenceProperties *ResultReference            `json:"#properties,omitzero"`
	Refs                map[string]*ResultReference `json:",embed"`
}

// ValidateArgs reports a local RFC 8620 §3.7 violation in the arguments.
// Request marshaling calls it; it is promoted to types that embed Get.
func (m *Get[T]) ValidateArgs() error {
	return checkRefs(m.Refs, map[string]bool{
		"ids":        !m.IDs.IsZero() || m.ReferenceIDs != nil,
		"properties": !m.Properties.IsZero() || m.ReferenceProperties != nil,
	})
}

func (m *Get[T]) Name() string {
	var zero T
	return zero.JMAPType() + "/get"
}

func (m *Get[T]) Requires() []URI {
	var zero T
	return zero.Requires()
}

// GetResponse is the result of a /get method call.
type GetResponse[T Object] struct {
	Account  ID     `json:"accountId,omitzero"`
	State    string `json:"state,omitzero"`
	List     []T    `json:"list"`
	NotFound []ID   `json:"notFound"`
}

// Changes fetches id deltas since a state (RFC 8620 §5.2).
type Changes[T Object] struct {
	Account    ID     `json:"accountId,omitzero"`
	SinceState string `json:"sinceState,omitzero"`
	// MaxChanges is UnsignedInt|null (RFC 8620 §5.2). nil omits (server default).
	// A supplied value MUST be > 0; ValidateArgs rejects 0.
	MaxChanges *UnsignedInt                `json:"maxChanges,omitzero"`
	Refs       map[string]*ResultReference `json:",embed"`
}

// ValidateArgs reports a local RFC 8620 violation in the arguments.
// Request marshaling calls it; it is promoted to types that embed Changes.
func (m *Changes[T]) ValidateArgs() error {
	if m.MaxChanges != nil && *m.MaxChanges == 0 {
		return fmt.Errorf("jmap: maxChanges must be greater than 0 (RFC 8620 §5.2)")
	}
	return checkRefs(m.Refs, nil)
}

func (m *Changes[T]) Name() string {
	var zero T
	return zero.JMAPType() + "/changes"
}

func (m *Changes[T]) Requires() []URI {
	var zero T
	return zero.Requires()
}

// ChangesResponse is the result of a /changes method call.
// Created/Updated/Destroyed are always []ID (Approach A).
type ChangesResponse struct {
	Account        ID     `json:"accountId,omitzero"`
	OldState       string `json:"oldState,omitzero"`
	NewState       string `json:"newState,omitzero"`
	HasMoreChanges bool   `json:"hasMoreChanges"`
	Created        []ID   `json:"created"`
	Updated        []ID   `json:"updated"`
	Destroyed      []ID   `json:"destroyed"`
}

// Query fetches object ids matching filter/sort (RFC 8620 §5.5).
// Type-specific Filter/Sort fields are added by thin packages via embedding.
type Query[T Object] struct {
	Account         ID                          `json:"accountId,omitzero"`
	Position        Int                         `json:"position,omitzero"`
	Anchor          Optional[ID]                `json:"anchor,omitzero"`
	ReferenceAnchor *ResultReference            `json:"#anchor,omitzero"`
	AnchorOffset    Int                         `json:"anchorOffset,omitzero"`
	Limit           *UnsignedInt                `json:"limit,omitzero"`
	CalculateTotal  bool                        `json:"calculateTotal,omitzero"`
	Refs            map[string]*ResultReference `json:",embed"`
}

// ValidateArgs reports a local RFC 8620 §3.7 violation in the arguments.
// Request marshaling calls it; it is promoted to types that embed Query.
func (m *Query[T]) ValidateArgs() error {
	return checkRefs(m.Refs, map[string]bool{
		"anchor": !m.Anchor.IsZero() || m.ReferenceAnchor != nil,
	})
}

func (m *Query[T]) Name() string {
	var zero T
	return zero.JMAPType() + "/query"
}

func (m *Query[T]) Requires() []URI {
	var zero T
	return zero.Requires()
}

// QueryResponse is the result of a /query method call.
type QueryResponse struct {
	Account             ID           `json:"accountId,omitzero"`
	QueryState          string       `json:"queryState,omitzero"`
	CanCalculateChanges bool         `json:"canCalculateChanges"`
	Position            UnsignedInt  `json:"position"`
	IDs                 []ID         `json:"ids"`
	Total               *UnsignedInt `json:"total,omitzero"`
	Limit               *UnsignedInt `json:"limit,omitzero"`
}

// QueryChanges fetches query result deltas since a query state (RFC 8620 §5.6).
// Type-specific Filter/Sort fields are added by thin packages via embedding.
type QueryChanges[T Object] struct {
	Account         ID     `json:"accountId,omitzero"`
	SinceQueryState string `json:"sinceQueryState,omitzero"`
	// MaxChanges is UnsignedInt|null (RFC 8620 §5.6). nil omits (server default).
	// A supplied value MUST be > 0; ValidateArgs rejects 0.
	MaxChanges      *UnsignedInt                `json:"maxChanges,omitzero"`
	UpToID          Optional[ID]                `json:"upToId,omitzero"`
	ReferenceUpToID *ResultReference            `json:"#upToId,omitzero"`
	CalculateTotal  bool                        `json:"calculateTotal,omitzero"`
	Refs            map[string]*ResultReference `json:",embed"`
}

// ValidateArgs reports a local RFC 8620 violation in the arguments.
// Request marshaling calls it; it is promoted to types that embed QueryChanges.
func (m *QueryChanges[T]) ValidateArgs() error {
	if m.MaxChanges != nil && *m.MaxChanges == 0 {
		return fmt.Errorf("jmap: maxChanges must be greater than 0 (RFC 8620 §5.6)")
	}
	return checkRefs(m.Refs, map[string]bool{
		"upToId": !m.UpToID.IsZero() || m.ReferenceUpToID != nil,
	})
}

func (m *QueryChanges[T]) Name() string {
	var zero T
	return zero.JMAPType() + "/queryChanges"
}

func (m *QueryChanges[T]) Requires() []URI {
	var zero T
	return zero.Requires()
}

// QueryChangesResponse is the result of a /queryChanges method call.
type QueryChangesResponse struct {
	Account       ID           `json:"accountId,omitzero"`
	OldQueryState string       `json:"oldQueryState,omitzero"`
	NewQueryState string       `json:"newQueryState,omitzero"`
	Removed       []ID         `json:"removed"`
	Added         []AddedItem  `json:"added"`
	Total         *UnsignedInt `json:"total,omitzero"`
}

// Set creates, updates, or destroys objects (RFC 8620 §5.3).
// T must be Creatable. Destroy-only types use a typed Set without create/update.
type Set[T Creatable] struct {
	Account          ID                          `json:"accountId,omitzero"`
	IfInState        Optional[string]            `json:"ifInState,omitzero"`
	Create           Optional[map[ID]*T]         `json:"create,omitzero"`
	Update           Optional[map[ID]Patch]      `json:"update,omitzero"`
	Destroy          Optional[[]ID]              `json:"destroy,omitzero"`
	ReferenceDestroy *ResultReference            `json:"#destroy,omitzero"`
	Refs             map[string]*ResultReference `json:",embed"`
}

// ValidateArgs reports a local RFC 8620 §3.7 violation in the arguments.
// Request marshaling calls it; it is promoted to types that embed Set.
func (m *Set[T]) ValidateArgs() error {
	return checkRefs(m.Refs, map[string]bool{
		"destroy": !m.Destroy.IsZero() || m.ReferenceDestroy != nil,
	})
}

func (m *Set[T]) Name() string {
	var zero T
	return zero.JMAPType() + "/set"
}

func (m *Set[T]) Requires() []URI {
	var zero T
	return zero.Requires()
}

// SetResponse is the result of a /set method call.
type SetResponse[T Object] struct {
	Account      ID                         `json:"accountId,omitzero"`
	OldState     Optional[string]           `json:"oldState,omitzero"`
	NewState     string                     `json:"newState,omitzero"`
	Created      map[ID]T                   `json:"created,omitzero"` // values are Foo, never null
	Updated      map[ID]*T                  `json:"updated,omitzero"` // values are Foo|null
	Destroyed    Optional[[]ID]             `json:"destroyed,omitzero"`
	NotCreated   Optional[map[ID]*SetError] `json:"notCreated,omitzero"`
	NotUpdated   Optional[map[ID]*SetError] `json:"notUpdated,omitzero"`
	NotDestroyed Optional[map[ID]*SetError] `json:"notDestroyed,omitzero"`
}

// Copy copies objects between accounts (RFC 8620 §5.4).
type Copy[T Creatable] struct {
	FromAccount              ID                          `json:"fromAccountId,omitzero"`
	IfFromInState            Optional[string]            `json:"ifFromInState,omitzero"`
	Account                  ID                          `json:"accountId,omitzero"`
	IfInState                Optional[string]            `json:"ifInState,omitzero"`
	Create                   Optional[map[ID]*T]         `json:"create,omitzero"`
	OnSuccessDestroyOriginal bool                        `json:"onSuccessDestroyOriginal,omitzero"`
	DestroyFromIfInState     Optional[string]            `json:"destroyFromIfInState,omitzero"`
	Refs                     map[string]*ResultReference `json:",embed"`
}

// ValidateArgs reports a local RFC 8620 §3.7 violation in the arguments.
// Request marshaling calls it; it is promoted to types that embed Copy.
func (m *Copy[T]) ValidateArgs() error {
	return checkRefs(m.Refs, nil)
}

func (m *Copy[T]) Name() string {
	var zero T
	return zero.JMAPType() + "/copy"
}

func (m *Copy[T]) Requires() []URI {
	var zero T
	return zero.Requires()
}

// CopyResponse is the result of a /copy method call.
type CopyResponse[T Object] struct {
	FromAccount ID                         `json:"fromAccountId,omitzero"`
	Account     ID                         `json:"accountId,omitzero"`
	OldState    Optional[string]           `json:"oldState,omitzero"`
	NewState    string                     `json:"newState,omitzero"`
	Created     map[ID]T                   `json:"created,omitzero"`
	NotCreated  Optional[map[ID]*SetError] `json:"notCreated,omitzero"`
}
