package interpreter

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	text := "gawk \"{ print(\"\"ahaha ihihi ufufu\"\") }\" <\"ddd\"\"ddd\"|ahaha \"ihihi |ufufu\" ; ohoho gegee&&hogehogeo >ihihi"
	fmt.Println(text)
	result, _ := Parse(text)
	for i, st := range result {
		fmt.Printf("pipeline-%d:\n", i)
		for _, stsub := range st {
			for _, word := range stsub.Args {
				fmt.Printf("  [%s]", word)
			}
			fmt.Println()
		}
	}
	result, _ = Parse("")
	fmt.Println("<empty-line>")
	for i, st := range result {
		fmt.Printf("pipeline-%d:\n", i)
		for _, stsub := range st {
			fmt.Printf("  %s\n", stsub.String())
		}
	}
}
