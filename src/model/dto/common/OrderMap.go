package common

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"reflect"
	"strings"
)

type OrderMap struct {
	m *linkedhashmap.Map
}

func NewOrderMap() *OrderMap {
	return &OrderMap{linkedhashmap.New()}
}

func (o *OrderMap) MarshalJSON() ([]byte, error) {
	if o.m == nil {
		o.m = linkedhashmap.New()
	}
	return o.m.ToJSON()
}

func (o *OrderMap) Put(key interface{}, value interface{}) {
	if o.m == nil {
		o.m = linkedhashmap.New()
	}
	o.m.Put(key, value)
}

func (o *OrderMap) Get(key interface{}) (value interface{}, found bool) {
	if o.m == nil {
		o.m = linkedhashmap.New()
	}
	return o.m.Get(key)
}

func (OrderMap) FromObject(object interface{}) *OrderMap {
	data := linkedhashmap.New()
	if object == nil {
		return nil
	}

	// check ptr and struct
	val := reflect.ValueOf(object)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return nil
	}
	relType := val.Type()

	// val, retType
	for i := 0; i < relType.NumField(); i++ {
		// !!
		tag := relType.Field(i).Tag.Get("json")
		omitempty := strings.Index(tag, "omitempty") != -1

		field := strings.Split(tag, ",")[0]
		value := val.Field(i).Interface()

		if field != "-" && (!omitempty || value != nil) {
			data.Put(field, value)
		}
	}
	return &OrderMap{m: data}
}
