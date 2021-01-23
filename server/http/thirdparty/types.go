package thirdparty

type QueryParam struct {
	Service      string `json:"service"`
	Source       string `json:"source"`
	MessageQuery string `json:"message_query"`
	ContextQuery string `json:"context_query"`
	Limit        int    `json:"limit"`
}
