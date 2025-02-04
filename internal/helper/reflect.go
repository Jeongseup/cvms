package helper

import "reflect"

func SetFieldByTag(obj interface{}, tag string, value string) {
	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Match struct field's tag with the given key
		if fieldType.Tag.Get("json") == tag && field.CanSet() {
			// Set the value to the field
			if field.Kind() == reflect.Slice {
				field.Set(reflect.ValueOf([]string{value}))
			} else {
				field.SetString(value)
			}
		}
	}
}
