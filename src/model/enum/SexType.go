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
	if sexString == string(SexMale) {
		return SexMale
	} else if sexString == string(SexFemale) {
		return SexFemale
	} else {
		return SexUnknown
	}
}
