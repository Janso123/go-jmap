package jmap

// Filter is a JMAP filter condition or operator. Types that implement Filter
// must embed FilterBase (or otherwise provide jmapFilter).
type Filter interface {
	jmapFilter()
}

// FilterBase is an empty embeddable marker so external packages and tests can
// satisfy Filter without accessing the unexported jmapFilter method.
type FilterBase struct{}

func (FilterBase) jmapFilter() {}

// FilterOperator combines Filter conditions with a logical Operator.
type FilterOperator struct {
	Operator   Operator `json:"operator"`
	Conditions []Filter `json:"conditions"`
}

func (FilterOperator) jmapFilter() {}

// And returns a FilterOperator that matches when all conditions match.
func And(conds ...Filter) *FilterOperator {
	return &FilterOperator{Operator: OperatorAND, Conditions: conds}
}

// Or returns a FilterOperator that matches when any condition matches.
func Or(conds ...Filter) *FilterOperator {
	return &FilterOperator{Operator: OperatorOR, Conditions: conds}
}

// Not returns a FilterOperator that matches when none of the conditions match.
func Not(conds ...Filter) *FilterOperator {
	return &FilterOperator{Operator: OperatorNOT, Conditions: conds}
}

// Comparator describes a sort key for Query methods.
// IsAscending has no omitzero: RFC default is true, so Desc must emit false.
type Comparator struct {
	Property    string        `json:"property"`
	IsAscending bool          `json:"isAscending"`
	Collation   CollationAlgo `json:"collation,omitzero"`
	Keyword     string        `json:"keyword,omitzero"`
}

// Asc returns a Comparator sorted ascending on prop.
func Asc(prop string) *Comparator {
	return &Comparator{Property: prop, IsAscending: true}
}

// Desc returns a Comparator sorted descending on prop.
func Desc(prop string) *Comparator {
	return &Comparator{Property: prop, IsAscending: false}
}
