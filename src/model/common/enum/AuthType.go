package enum

type AuthType string

const (
	AuthAdmin  AuthType = "admin"
	AuthNormal AuthType = "normal"
)

func (s AuthType) String() string {
	return string(s)
}
