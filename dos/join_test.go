package dos

import (
	"fmt"
	"testing"
)

func _testJoin(a, b string) {
	fmt.Printf("'%s'+'%s' -> '%s'\n", a, b, Join(a, b))
}

func TestJoin(t *testing.T) {
	_testJoin(`foo`, `bar`)
	_testJoin(`foo`, `\bar`)
	_testJoin(`c:`, `bar`)
	_testJoin(`foo/`, `bar`)
	_testJoin(`foo\`, `bar`)
	_testJoin(`foo`, `c:bar`)
	_testJoin(`foo`, `c:\bar`)
	_testJoin(`c:foo`, `\bar`)
	_testJoin(`c:foo`, `\\host\path\to`)
	_testJoin(`\\host\path\to`, `c:foo`)
}
