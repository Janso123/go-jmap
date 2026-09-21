package quota

import "github.com/Janso123/go-jmap"

// Filter criteria for quota queries.
type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	Name string `json:"name,omitzero"`

	Scope string `json:"scope,omitzero"`

	ResourceType string `json:"resourceType,omitzero"`

	Type string `json:"type,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }
