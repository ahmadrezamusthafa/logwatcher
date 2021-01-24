package types

import "github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"

type Attribute struct {
	Name     string              `json:"name"`
	Operator string              `json:"operator"`
	Value    string              `json:"value"`
	Type     valuetype.ValueType `json:"type,omitempty"`
}

type TokenAttribute struct {
	Value     string
	HasCalled bool
}
