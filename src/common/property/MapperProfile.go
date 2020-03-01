package property

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
)

func CreatePropertyMappers() *xproperty.PropertyMappers {
	mappers := xproperty.NewPropertyMappers()

	addPropMappers(mappers)

	return mappers
}
