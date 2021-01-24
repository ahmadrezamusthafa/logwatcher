package valuetype

type ValueType string

const (
	Numeric      ValueType = "numeric"
	Alphanumeric ValueType = "alphanumeric"
	Date         ValueType = "date"
)

func FromString(value string) ValueType {
	return ValueType(value)
}

func (v ValueType) ToString() string {
	return string(v)
}
