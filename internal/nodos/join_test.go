package nodos_test

import (
	"testing"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

func _testJoin(t *testing.T, a, b, expect string) {
	result := nodos.Join(a, b)
	if result != expect {
		t.Fatalf("'%s'+'%s' should be '%s',but '%s'\n", a, b, expect, result)
	}
}

func TestJoin(t *testing.T) {
	_testJoin(t, `foo`, `bar`, `foo\bar`)
	_testJoin(t, `foo`, `\bar`, `\bar`)
	_testJoin(t, `c:`, `bar`, `c:bar`)
	_testJoin(t, `foo/`, `bar`, `foo/bar`)
	_testJoin(t, `foo\`, `bar`, `foo\bar`)
	_testJoin(t, `foo`, `c:bar`, `c:bar`)
	_testJoin(t, `foo`, `c:\bar`, `c:\bar`)
	_testJoin(t, `c:foo`, `\bar`, `c:\bar`)
	_testJoin(t, `c:foo`, `\\host\path\to`, `\\host\path\to`)
	_testJoin(t, `\\host\path\to`, `c:foo`, `c:foo`)
}
