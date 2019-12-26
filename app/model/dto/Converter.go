package dto

import (
	"github.com/shomali11/util/xconditions"
	"reflect"
	"strings"
)

type Converter struct {
	FieldType reflect.Type
	Converter func(interface{})
}

// !!
// data must be pointer / slice / array / struct
func (r *Result) convert(data interface{}) {
	dataType := reflect.TypeOf(data)

	if dataType.Kind() == reflect.Slice || dataType.Kind() == reflect.Array {
		elsVal := reflect.ValueOf(data)
		for i := 0; i < elsVal.Len(); i++ {
			val := elsVal.Index(i)
			r.convert(xconditions.IfThenElse(val.Type().Kind() == reflect.Ptr, val.Interface(), val.Addr().Interface()))
		}
	} else {
		for _, converter := range r.Converter {
			if dataType == converter.FieldType {
				converter.Converter(data)
			}
		}

		dataValue := reflect.ValueOf(data)
		if dataType.Kind() == reflect.Ptr {
			dataValue = dataValue.Elem()
		}
		dataType = dataValue.Type()

		if dataType.Kind() == reflect.Struct {
			for i := 0; i < dataType.NumField(); i++ {
				// fmt.Println(dataType.Field(i).Type, dataType.Field(i).Name)
				jsonTag := dataType.Field(i).Tag.Get("json")
				// !!
				if jsonTag == "" || strings.Split(jsonTag, ",")[0] == "-" {
					continue
				}
				val := dataValue.Field(i)
				r.convert(xconditions.IfThenElse(val.Type().Kind() == reflect.Ptr, val.Interface(), val.Addr().Interface()))
			}
		}
	}
}
