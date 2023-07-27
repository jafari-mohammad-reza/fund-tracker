package utils

import (
	"reflect"
	"sort"
)

type Item struct {
	Base   string
	Target string
}

func FindIndex(items []Item, base string, target string) int {
	for i, v := range items {
		if v.Base == base && v.Target == target {
			return i
		}
	}
	return -1
}

func SortResponseDataItems[T any](responseData []T, fieldName string) {
	sort.Slice(responseData, func(i, j int) bool {
		// Use reflection to get the field value based on the fieldName

		fieldI := reflect.ValueOf(responseData[i]).FieldByName(fieldName)
		fieldJ := reflect.ValueOf(responseData[j]).FieldByName(fieldName)

		// Compare the field values based on their types
		switch fieldI.Kind() {
		case reflect.String:
			return fieldI.String() < fieldJ.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fieldI.Int() < fieldJ.Int()
		// Add more cases for other types if needed
		default:
			// If the field type is not supported for comparison, you might want to handle it accordingly.
			// For example, return false to maintain the original order.
			return false
		}
	})
}
