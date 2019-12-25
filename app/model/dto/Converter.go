package dto

import (
	"reflect"
)

type Converter struct {
	FieldType reflect.Type
	Converter func(interface{}) interface{}
}

func (r *Result) convert(data interface{}) interface{} {
	dataType := reflect.TypeOf(data)

	if dataType.Kind() != reflect.Slice && dataType.Kind() != reflect.Array {
		for _, converter := range r.Converter {
			if dataType == converter.FieldType || dataType == reflect.PtrTo(converter.FieldType) {
				data = converter.Converter(data)
			}
		}
	} else {
		sliceValue := reflect.ValueOf(data)
		sliceValueType := dataType.Elem()

		for _, converter := range r.Converter {
			isValuePlain := sliceValueType == converter.FieldType
			isValuePtr := sliceValueType == reflect.PtrTo(converter.FieldType)

			if isValuePlain || isValuePtr {
				for i := 0; i < sliceValue.Len(); i++ {
					sliceValue.Index(i).Set(reflect.ValueOf(converter.Converter(sliceValue.Index(i).Interface())))
				}
			}
		}
	}
	return data
}
