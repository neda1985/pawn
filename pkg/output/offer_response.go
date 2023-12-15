package output
type PawnResponse struct {
	Code  string `json:"code"`
	Message string `json:"message,omitempty"`
	Value int    `json:"value,omitempty"`
}
