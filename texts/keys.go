package texts

import (
	"reflect"
	"sort"
)

// SortedKeys makes sorted strings' array from keys of the given map whose key's type is string.
func SortedKeys(mapInt interface{}) []string {
	values := reflect.ValueOf(mapInt).MapKeys()
	result := make([]string, len(values))
	for i, value1 := range values {
		result[i] = value1.String()
	}
	sort.Strings(result)
	return result
}
