package ignoreCaseSorted

import (
	"testing"
)

func TestAscend(t *testing.T) {
	dic := make_testdata()
	expect := make_expect_ascend()

	for p := dic.Front(); p != nil; p = p.Next() {
		if expect[0].key != p.Key {
			t.Fatalf("'%s' != '%s'", expect[0].key, p.Key)
			return
		}
		if expect[0].value != p.Value {
			t.Fatalf("'%d' != '%d'", expect[0].value, p.Value)
			return
		}
		expect = expect[1:]
	}
}

func TestDesend(t *testing.T) {
	dic := make_testdata()
	expect := make_expect_descend()

	for p := dic.Back(); p != nil; p = p.Prev() {
		if expect[0].key != p.Key {
			t.Fatalf("'%s' != '%s'", expect[0].key, p.Key)
			return
		}
		if expect[0].value != p.Value {
			t.Fatalf("'%d' != '%d'", expect[0].value, p.Value)
			return
		}
		expect = expect[1:]
	}
}
