package models

type ErrorInfo struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Time    string `json:"time"`
}
