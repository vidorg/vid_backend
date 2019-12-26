package dto

import (
	"reflect"
)

type Converter struct {
	FieldType reflect.Type
	Converter func(interface{})
}

// data must be pointer / slice / array / struct
func (r *Result) convert(data interface{}) {
	dataType := reflect.TypeOf(data)

	if dataType.Kind() == reflect.Slice || dataType.Kind() == reflect.Array { // []
		sliceValue := reflect.ValueOf(data) // value of slice element
		for _, converter := range r.Converter {
			isTypeEq := dataType.Elem() == converter.FieldType                 // same type
			isTypePtr := reflect.PtrTo(dataType.Elem()) == converter.FieldType // convert is ptr, element type is plain

			if isTypeEq || isTypePtr {
				for i := 0; i < sliceValue.Len(); i++ {
					if isTypeEq {
						r.convert(sliceValue.Index(i).Interface()) // slice is pointer
					} else if isTypePtr {
						r.convert(sliceValue.Index(i).Addr().Interface()) // take slice element pointer
					}
				}
			}
		}
	} else { // {}
		for _, converter := range r.Converter {
			if dataType == converter.FieldType { // only support pointer
				converter.Converter(data) // <<<
			}
		}
		if dataType.Kind() == reflect.Struct { // recursion
			dataValue := reflect.ValueOf(data) // value of struct field
			for _, converter := range r.Converter {
				for i := 0; i < dataType.NumField(); i++ {
					isTypeEq := dataValue.Field(i).Type() == converter.FieldType                 // same type
					isTypePtr := reflect.PtrTo(dataValue.Field(i).Type()) == converter.FieldType // convert is ptr, field type is plain

					if isTypeEq {
						r.convert(dataValue.Field(i)) // value is pointer
					} else if isTypePtr {
						r.convert(dataValue.Field(i).Addr().Interface()) // take struct field pointer
					}
				}
			}
		}
	}
}
