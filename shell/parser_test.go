package shell_test

import (
	"fmt"
	"testing"

	"github.com/zetamatta/nyagos/shell"
)

func TestParser(t *testing.T) {
	text := "gawk \"{ print(\"\"ahaha ihihi ufufu\"\") }\" <\"ddd\"\"ddd\"|ahaha \"ihihi |ufufu\" ; ohoho gegee&&hogehogeo >ihihi"
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
