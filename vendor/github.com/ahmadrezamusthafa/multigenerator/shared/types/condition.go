package types

type BaseCondition struct {
	Fields     []string     `json:"fields,omitempty"`
	Distinct   bool         `json:"distinct"`
	Conditions []*Condition `json:"conditions,omitempty"`
	Footer     Footer       `json:"footer"`
}

type Condition struct {
	Operator   string       `json:"operator,omitempty"`
	Attribute  *Attribute   `json:"attribute,omitempty"`
	Conditions []*Condition `json:"conditions,omitempty"`
}
