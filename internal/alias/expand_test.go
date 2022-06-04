package alias_test

import (
	"testing"

	"github.com/nyaosorg/nyagos/internal/alias"
)

func TestExpandMacro(t *testing.T) {
	args := []string{"<0>", "<1>", "<2>", "<3>"}
	rawargs := []string{`"<0>"`, `"<1>"`, `"<2>"`, `"<3>"`}
	result := alias.ExpandMacro("foo $0 $1 $2 $3", args, rawargs)
	if result != `foo "<0>" "<1>" "<2>" "<3>"` {
		t.Fatalf("$0...$3 error: %s", result)
	}
	result = alias.ExpandMacro("foo", args, rawargs)
	if result != `foo "<1>" "<2>" "<3>"` {
		t.Fatalf("no $n error: %s", result)
	}
	result = alias.ExpandMacro("foo $~0 $~1 $~2 $~3", args, rawargs)
	if result != `foo <0> <1> <2> <3>` {
		t.Fatalf("$~0...$~3 error: %s", result)
	}
}
