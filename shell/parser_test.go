package shell_test

import (
	"fmt"
	"testing"

	"github.com/zetamatta/nyagos/shell"
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
	fmt.Println(text)
	result, _ := shell.Parse(new(shell.NulStream), text)
	for i, st := range result {
		fmt.Printf("pipeline-%d:\n", i)
		for _, stsub := range st {
			for _, word := range stsub.Args {
				fmt.Printf("  [%s]", word)
			}
			fmt.Println()
		}
	}
	result, _ = shell.Parse(new(shell.NulStream), "")
	fmt.Println("<empty-line>")
	for i, st := range result {
		fmt.Printf("pipeline-%d:\n", i)
		for _, stsub := range st {
			for _, word := range stsub.Args {
				fmt.Printf("  [%s]", word)
			}
			fmt.Println()
		}
	}
}
