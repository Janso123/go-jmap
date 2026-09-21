package mailbox

import "github.com/Janso123/go-jmap"

type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	ParentID jmap.Optional[jmap.ID] `json:"parentId,omitzero"`

	Name string `json:"name,omitzero"`

	Role jmap.Optional[Role] `json:"role,omitzero"`

	HasAnyRole *bool `json:"hasAnyRole,omitzero"`

	IsSubscribed *bool `json:"isSubscribed,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }

// TopLevel matches mailboxes whose parentId is JSON null.
func TopLevel() *FilterCondition {
	return &FilterCondition{ParentID: jmap.Null[jmap.ID]()}
}

func Parent(id jmap.ID) *FilterCondition {
	return &FilterCondition{ParentID: jmap.Some(id)}
}
