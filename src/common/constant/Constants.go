package constant

type SexEnum string

func (s SexEnum) String() string {
	return string(s)
}

const (
	AuthAdmin string = "admin"

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
