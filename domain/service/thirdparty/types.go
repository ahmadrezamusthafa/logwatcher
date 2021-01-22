package thirdparty

type QueryInput struct {
	MessageQuery string
	ContextQuery string
	Limit        int
}

type QueryOutput struct {
	Timestamp *string                `json:"timestamp"`
	Hostname  *string                `json:"hostname"`
	FlowID    *string                `json:"flowid"`
	Type      *string                `json:"type"`
	Part      *string                `json:"part"`
	Message   *string                `json:"message"`
	Context   map[string]interface{} `json:"context"`
}

type LogAttribute struct {
	Attribute string `json:"attribute"`
	Name      string `json:"name"`
}
