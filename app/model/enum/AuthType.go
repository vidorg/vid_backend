package enum

type AuthType string

const (
	AuthAdmin  AuthType = "admin"
	AuthNormal AuthType = "normal"
)

func (s AuthType) String() string {
	return string(s)
}

func (AuthType) FromString(authString string) AuthType {
	if authString == string(AuthAdmin) {
		return AuthAdmin
	} else {
		return AuthNormal
	}
}
