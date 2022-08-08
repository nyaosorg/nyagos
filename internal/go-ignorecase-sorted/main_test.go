package ignoreCaseSorted

import (
	"testing"
)

func make_testdata() *Dictionary[int] {
	var dic Dictionary[int]

	dic.Set("x", 7)
	dic.Set("y", 8)
	dic.Set("Z", 9)

	return &dic
}

type ExpectT struct {
	key   string
	value int
}

func make_expect_ascend() []ExpectT {
	return []ExpectT{
		{key: "x", value: 7},
		{key: "y", value: 8},
		{key: "Z", value: 9},
	}
}

func make_expect_descend() []ExpectT {
	return []ExpectT{
		{key: "Z", value: 9},
		{key: "y", value: 8},
		{key: "x", value: 7},
	}
}

func TestRange(t *testing.T) {
	dic := make_testdata()
	expect := make_expect_ascend()

	dic.Range(func(key string, value int) bool {
		if expect[0].key != key {
			t.Fatalf("'%s' != '%s'", expect[0].key, key)
			return false
		}
		if expect[0].value != value {
			t.Fatalf("'%d' != '%d'", expect[0].value, value)
			return false
		}
		expect = expect[1:]
		return true
	})
}
