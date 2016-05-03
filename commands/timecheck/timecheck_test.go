package timecheck

import "testing"

type testlist_t struct {
	r bool
	d []int
}

func TestIsOk(t *testing.T) {
	testlist := []testlist_t{
		testlist_t{true, []int{2016, 5, 13, 17, 20, 0}},
		testlist_t{true, []int{2016, 2, 29, 17, 20, 0}},
		testlist_t{false, []int{2015, 2, 29, 17, 20, 0}},
		testlist_t{false, []int{2016, 14, 20, 17, 20, 0}},
		testlist_t{false, []int{2016, 12, 32, 17, 20, 0}},
		testlist_t{false, []int{2016, 12, 31, 24, 20, 0}},
		testlist_t{false, []int{2016, 12, 31, 23, 70, 0}},
	}
	for _, p := range testlist {
		d := p.d
		if p.r != IsOk(d[0], d[1], d[2], d[3], d[4], d[5]) {
			t.Fatalf("[NG] %d/%d/%d %d:%d:%d\n", d[0], d[1], d[2], d[3], d[4], d[5])
			t.Fail()
		}
	}
}
