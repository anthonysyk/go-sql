package copy

import (
	"reflect"
)

// GetValues TODO: take columns as input
func GetValues(item interface{}) []interface{} {
	itemType := reflect.TypeOf(item)
	var values []interface{}
	for i := 0; i < itemType.NumField(); i++ {
		value := reflect.ValueOf(item).Field(i).Interface()
		values = append(values, value)
	}
	return values
}
