package ignoreCaseSorted

import (
	"testing"
)

func TestAscend2(t *testing.T) {
	dic := make_testdata()
	expect := make_expect_ascend()

	for p := dic.Ascend(); p.Range(); {
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

func TestDesend2(t *testing.T) {
	dic := make_testdata()
	expect := make_expect_descend()

	for p := dic.Descend(); p.Range(); {
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
