package models

type Message struct {
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}
