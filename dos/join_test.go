package dos

import (
	"fmt"
	"testing"
)

func testJoin_(a, b string) {
	fmt.Printf("'%s'+'%s' -> '%s'\n", a, b, Join(a, b))
}

func TestJoin(t *testing.T) {
	testJoin_(`foo`, `bar`)
	testJoin_(`foo`, `\bar`)
	testJoin_(`c:`, `bar`)
	testJoin_(`foo/`, `bar`)
	testJoin_(`foo\`, `bar`)
	testJoin_(`foo`, `c:bar`)
	testJoin_(`foo`, `c:\bar`)
	testJoin_(`c:foo`, `\bar`)
	testJoin_(`c:foo`, `\\host\path\to`)
	testJoin_(`\\host\path\to`, `c:foo`)
}
