package enum

type AuthType string

const (
	AuthAdmin  AuthType = "admin"
	AuthNormal AuthType = "normal"
)

type SexType string

const (
	SexUnknown SexType = "unknown"
	SexMale    SexType = "male"
	SexFemale  SexType = "female"
)

func StringToSex(sexString string) SexType {
	if sexString == string(SexMale) {
		return SexMale
	} else if sexString == string(SexFemale) {
		return SexFemale
	} else {
		return SexUnknown
	}
}