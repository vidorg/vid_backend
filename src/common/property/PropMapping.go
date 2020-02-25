package property

type PropMappingValue struct {
	DestProps []string
	Revert    bool
}

func NewPropMappingValue(destProps []string, revert bool) *PropMappingValue {
	if destProps == nil {
		destProps = make([]string, 0)
	}
	return &PropMappingValue{
		DestProps: destProps,
		Revert:    revert,
	}
}

type PropMapping struct {
	DtoModel interface{}
	PoModel  interface{}
	Dict     map[string]*PropMappingValue // dto -> po
}

func NewPropMapping(dtoModel interface{}, poModel interface{}, dict map[string]*PropMappingValue) *PropMapping {
	return &PropMapping{
		DtoModel: dtoModel,
		PoModel:  poModel,
		Dict:     dict,
	}
}
