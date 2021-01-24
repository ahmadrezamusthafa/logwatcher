package types

type Footer struct {
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
	Sort  map[string]string `json:"sort"`
}
