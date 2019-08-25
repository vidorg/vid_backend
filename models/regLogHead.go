package models

var MinLen_Username int
var MaxLen_Username int
var MinLen_Password int
var MaxLen_Password int

type RegLogHead struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @override
func (r *RegLogHead) CheckValid() bool {
	return r.Username != "" && r.Password != ""
}

// @override
func (r *RegLogHead) CheckFormat() bool {
	return len(r.Username) >= MinLen_Username &&
		len(r.Username) <= MaxLen_Username &&
		len(r.Password) >= MinLen_Password &&
		len(r.Password) <= MaxLen_Password
}
