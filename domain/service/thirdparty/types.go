package thirdparty

type QueryInput struct {
	MessageQuery string
	ContextQuery string
}

type LogAttribute struct {
	Attribute string `json:"attribute"`
	Name      string `json:"name"`
}
