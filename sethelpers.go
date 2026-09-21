package jmap

import "fmt"

const (
	PathIDs               = "/ids"
	PathCreated           = "/created"
	PathUpdated           = "/updated"
	PathUpdatedProperties = "/updatedProperties"
)

// ListPath returns a JSON Pointer into list items for prop (e.g. "/list/*/subject").
func ListPath(prop string) string { return "/list/*/" + prop }

// CreationRef returns a client creation-id reference ("#" + createID).
func CreationRef(createID ID) ID { return "#" + createID }

// SetCreated returns the created object for id, or the corresponding SetError.
func SetCreated[T any](created map[ID]T, notCreated map[ID]*SetError, id ID) (T, error) {
	var zero T
	if se, ok := notCreated[id]; ok && se != nil {
		return zero, se
	}
	v, ok := created[id]
	if !ok {
		return zero, fmt.Errorf("id %q not in created", id)
	}
	return v, nil
}

// Ref builds a ResultReference to a prior call in this request.
func (r *Request) Ref(callID, path string) *ResultReference {
	for _, inv := range r.Calls {
		if inv.CallID == callID {
			return &ResultReference{ResultOf: callID, Name: inv.Name, Path: path}
		}
	}
	return &ResultReference{ResultOf: callID, Path: path}
}
