package texts_test

import (
	"testing"

	"github.com/nyaosorg/nyagos/texts"
)

func TestSortedKeys(t *testing.T) {
	example := map[string]int{
		"a": 1, "b": 2, "c": 3,
	}

	result := texts.SortedKeys(example)

	if result[0] != "a" ||
		result[1] != "b" ||
		result[2] != "c" {

		t.Fail()
	}
}
