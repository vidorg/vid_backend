package enum

type DtoOption int8

const (
	DtoOptionNone DtoOption = iota
	DtoOptionOnlySelf
	DtoOptionAll
)
