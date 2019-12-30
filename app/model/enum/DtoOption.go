package enum

type DtoOption int

const (
	DtoOptionNone DtoOption = iota
	DtoOptionOnlySelf
	DtoOptionAll
)
