package resp

type Message struct {
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}
