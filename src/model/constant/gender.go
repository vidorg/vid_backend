package constant

type Gender int8

const (
	Secret Gender = iota
	Male
	Female
)

func (g Gender) Data() int8 {
	return int8(g)
}

func (g Gender) String() string {
	switch g {
	case Secret:
		return "secret"
	case Male:
		return "male"
	case Female:
		return "female"
	default:
		return "<unknown gender>"
	}
}

func ParseGender(data int8) Gender {
	switch data {
	case Secret.Data():
		return Secret
	case Male.Data():
		return Male
	case Female.Data():
		return Female
	default:
		return Secret
	}
}
