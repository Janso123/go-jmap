package jmap

// Get fetches objects by id (RFC 8620 §5.1).
type Get[T Object] struct {
	Account             ID               `json:"accountId,omitzero"`
	IDs                 []ID             `json:"ids,omitzero"`
	Properties          []string         `json:"properties,omitzero"`
	ReferenceIDs        *ResultReference `json:"#ids,omitzero"`
	ReferenceProperties *ResultReference `json:"#properties,omitzero"`
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
	List     []T    `json:"list,omitzero"`
	NotFound []ID   `json:"notFound,omitzero"`
}

// Changes fetches id deltas since a state (RFC 8620 §5.2).
type Changes[T Object] struct {
	Account    ID     `json:"accountId,omitzero"`
	SinceState string `json:"sinceState,omitzero"`
	MaxChanges uint64 `json:"maxChanges,omitzero"`
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
	HasMoreChanges bool   `json:"hasMoreChanges,omitzero"`
	Created        []ID   `json:"created,omitzero"`
	Updated        []ID   `json:"updated,omitzero"`
	Destroyed      []ID   `json:"destroyed,omitzero"`
}

// Query fetches object ids matching filter/sort (RFC 8620 §5.5).
// Type-specific Filter/Sort fields are added by thin packages via embedding.
type Query[T Object] struct {
	Account         ID               `json:"accountId,omitzero"`
	Position        int64            `json:"position,omitzero"`
	Anchor          ID               `json:"anchor,omitzero"`
	ReferenceAnchor *ResultReference `json:"#anchor,omitzero"`
	AnchorOffset    int64            `json:"anchorOffset,omitzero"`
	Limit           *uint64          `json:"limit,omitzero"`
	CalculateTotal  bool             `json:"calculateTotal,omitzero"`
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
	Account             ID     `json:"accountId,omitzero"`
	QueryState          string `json:"queryState,omitzero"`
	CanCalculateChanges bool   `json:"canCalculateChanges,omitzero"`
	Position            uint64 `json:"position,omitzero"`
	IDs                 []ID   `json:"ids,omitzero"`
	Total               uint64 `json:"total,omitzero"`
	Limit               uint64 `json:"limit,omitzero"`
}

// QueryChanges fetches query result deltas since a query state (RFC 8620 §5.6).
// Type-specific Filter/Sort fields are added by thin packages via embedding.
type QueryChanges[T Object] struct {
	Account         ID               `json:"accountId,omitzero"`
	SinceQueryState string           `json:"sinceQueryState,omitzero"`
	MaxChanges      uint64           `json:"maxChanges,omitzero"`
	UpToID          ID               `json:"upToId,omitzero"`
	ReferenceUpToID *ResultReference `json:"#upToId,omitzero"`
	CalculateTotal  bool             `json:"calculateTotal,omitzero"`
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
	Account       ID          `json:"accountId,omitzero"`
	OldQueryState string      `json:"oldQueryState,omitzero"`
	NewQueryState string      `json:"newQueryState,omitzero"`
	Removed       []ID        `json:"removed,omitzero"`
	Added         []AddedItem `json:"added,omitzero"`
	Total         uint64      `json:"total,omitzero"`
}

// Set creates, updates, or destroys objects (RFC 8620 §5.3).
type Set[T Object] struct {
	Account          ID               `json:"accountId,omitzero"`
	IfInState        string           `json:"ifInState,omitzero"`
	Create           map[ID]*T        `json:"create,omitzero"`
	Update           map[ID]Patch     `json:"update,omitzero"`
	Destroy          []ID             `json:"destroy,omitzero"`
	ReferenceDestroy *ResultReference `json:"#destroy,omitzero"`
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
	Account      ID               `json:"accountId,omitzero"`
	OldState     string           `json:"oldState,omitzero"`
	NewState     string           `json:"newState,omitzero"`
	Created      map[ID]*T        `json:"created,omitzero"`
	Updated      map[ID]*T        `json:"updated,omitzero"`
	Destroyed    []ID             `json:"destroyed,omitzero"`
	NotCreated   map[ID]*SetError `json:"notCreated,omitzero"`
	NotUpdated   map[ID]*SetError `json:"notUpdated,omitzero"`
	NotDestroyed map[ID]*SetError `json:"notDestroyed,omitzero"`
}

// Copy copies objects between accounts (RFC 8620 §5.4).
type Copy[T Object] struct {
	FromAccount              ID        `json:"fromAccountId,omitzero"`
	IfFromInState            string    `json:"ifFromInState,omitzero"`
	Account                  ID        `json:"accountId,omitzero"`
	IfInState                string    `json:"ifInState,omitzero"`
	Create                   map[ID]*T `json:"create,omitzero"`
	OnSuccessDestroyOriginal bool      `json:"onSuccessDestroyOriginal,omitzero"`
	DestroyFromIfInState     string    `json:"destroyFromIfInState,omitzero"`
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
	FromAccount ID               `json:"fromAccountId,omitzero"`
	Account     ID               `json:"accountId,omitzero"`
	OldState    string           `json:"oldState,omitzero"`
	NewState    string           `json:"newState,omitzero"`
	Created     map[ID]*T        `json:"created,omitzero"`
	NotCreated  map[ID]*SetError `json:"notCreated,omitzero"`
}
