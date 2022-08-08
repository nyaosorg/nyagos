package shell_test

import (
	"testing"

	"github.com/nyaosorg/nyagos/internal/shell"
)

func TestParserForAwk(t *testing.T) {
	source := `gawk "BEGIN{ FS=\"\v\" ; RS=\"\f\" } { printf \"%d: [%s]\n\",$2,$1 }"`
	actual, _ := shell.Parse(new(shell.NulStream), source)

	rawExpect := `"BEGIN{ FS=\"\v\" ; RS=\"\f\" } { printf \"%d: [%s]\n\",$2,$1 }"`
	if act := actual[0][0].RawArgs[1]; act != rawExpect {
		t.Fatalf("shell.Parse(`%s`) failed: expect `%s` as raw-string but `%s`",
			source, rawExpect, act)
	}
	expect := `BEGIN{ FS="\v" ; RS="\f" } { printf "%d: [%s]\n",$2,$1 }`
	if act := actual[0][0].Args[1]; act != expect {
		t.Fatalf("shell.Parse(`%s`) failed: expect `%s` as cooked-string but `%s`",
			source, expect, act)
	}
}

func TestParser(t *testing.T) {
	text := `gawk "{ print(""ahaha ihihi ufufu"") }" <"ddd""ddd"|ahaha "ihihi |ufufu" ; ohoho gegee&&hogehogeo >ihihi`
	result, _ := shell.Parse(new(shell.NulStream), text)

	if result[0][0].Args[0] != `gawk` {
		t.Fatal("Check-1")
	}
	if result[0][0].Args[1] != `{ print("ahaha ihihi ufufu") }` {
		t.Fatal("Check-2")
	}
	if result[0][1].Args[0] != `ahaha` {
		t.Fatal("Check-3")
	}
	if result[0][1].Args[1] != `ihihi |ufufu` {
		t.Fatal("Check-4")
	}
	if result[1][0].Args[0] != `ohoho` {
		t.Fatal("Check-5")
	}
	if result[1][0].Args[1] != `gegee` {
		t.Fatal("Check-6")
	}
	if result[2][0].Args[0] != `hogehogeo` {
		t.Fatal("Check-7")
	}
	result, _ = shell.Parse(new(shell.NulStream), "")
	if len(result) > 0 {
		t.Fatal("Check-8")
	}
}
