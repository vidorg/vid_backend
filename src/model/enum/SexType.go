package enum

type SexType string

const (
	SexUnknown SexType = "unknown"
	SexMale    SexType = "male"
	SexFemale  SexType = "female"
)

func (s SexType) String() string {
	return string(s)
}

func StringToSexType(sexString string) SexType {
	switch sexString {
	case SexMale.String():
		return SexMale
	case SexFemale.String():
		return SexFemale
	default:
		return SexUnknown
	}
}
