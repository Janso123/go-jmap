package jmap

import (
	"fmt"

	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

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

// UnmarshalFilter decodes a JMAP filter tree (RFC 8620 §5.5). Leaves decode
// into *C, which must implement Filter (embed FilterBase). JSON null returns nil.
func UnmarshalFilter[C any](raw jsontext.Value) (Filter, error) {
	if raw.Kind() == 'n' || len(raw) == 0 {
		return nil, nil
	}
	var probe struct {
		Operator   Operator         `json:"operator"`
		Conditions []jsontext.Value `json:"conditions"`
	}
	if err := jsonv2.Unmarshal(raw, &probe); err != nil {
		return nil, err
	}
	if probe.Operator == "" {
		c := new(C)
		if err := jsonv2.Unmarshal(raw, c); err != nil {
			return nil, err
		}
		f, ok := any(c).(Filter)
		if !ok {
			return nil, fmt.Errorf("jmap: %T does not implement Filter", c)
		}
		return f, nil
	}
	op := &FilterOperator{Operator: probe.Operator, Conditions: make([]Filter, 0, len(probe.Conditions))}
	for _, rc := range probe.Conditions {
		sub, err := UnmarshalFilter[C](rc)
		if err != nil {
			return nil, err
		}
		op.Conditions = append(op.Conditions, sub)
	}
	return op, nil
}

// Comparator describes a sort key for Query methods.
// IsAscending nil omits the field so the RFC 8620 §5.5 default (true) applies;
// Desc emits false.
type Comparator struct {
	Property    string        `json:"property"`
	IsAscending *bool         `json:"isAscending,omitzero"`
	Collation   CollationAlgo `json:"collation,omitzero"`
	Keyword     string        `json:"keyword,omitzero"`
}

// Asc returns a Comparator sorted ascending on prop (RFC default).
func Asc(prop string) *Comparator {
	return &Comparator{Property: prop}
}

// Desc returns a Comparator sorted descending on prop.
func Desc(prop string) *Comparator {
	return &Comparator{Property: prop, IsAscending: new(false)}
}
