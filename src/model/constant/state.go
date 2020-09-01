package constant

type UserState int8

const (
	Inactive UserState = iota
	Active
	Suspend
)

func (u UserState) Data() int8 {
	return int8(u)
}

func (u UserState) String() string {
	switch u {
	case Inactive:
		return "inactive"
	case Active:
		return "active"
	case Suspend:
		return "suspend"
	default:
		return "<unknown user state>"
	}
}
