package dos

import (
	"fmt"
	"testing"
)

func testJoin_(a, b string) {
	fmt.Printf("'%s'+'%s' -> '%s'\n", a, b, Join(a, b))
}

func TestJoin(t *testing.T) {
	testJoin_("a", "b")
	testJoin_("a", "\\b")
	testJoin_("a:", "b")
	testJoin_("a/", "b")
	testJoin_("a\\", "b")
	testJoin_("a", "c:b")
	testJoin_("a", "c:\\b")
}
