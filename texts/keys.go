package texts

import (
	"reflect"
	"sort"
)

func SortedKeys(mapInt interface{}) []string {
	values := reflect.ValueOf(mapInt).MapKeys()
	result := make([]string, len(values))
	for i, value1 := range values {
		result[i] = value1.String()
	}
	sort.Strings(result)
	return result
}
