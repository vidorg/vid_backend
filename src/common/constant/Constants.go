package constant

type AuthEnum string

func (s AuthEnum) String() string {
	return string(s)
}

type SexEnum string

func (s SexEnum) String() string {
	return string(s)
}

const (
	AuthAdmin  AuthEnum = "admin"
	AuthNormal AuthEnum = "normal"

	SexUnknown SexEnum = "unknown"
	SexMale    SexEnum = "male"
	SexFemale  SexEnum = "female"
)

func ParseSexEnum(sexString string) SexEnum {
	switch sexString {
	case SexMale.String():
		return SexMale
	case SexFemale.String():
		return SexFemale
	default:
		return SexUnknown
	}
}
