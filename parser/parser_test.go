package parser

import "testing"
import "fmt"

func TestParser(t *testing.T) {
	result := Parse("aaaa bbbb \"\" cccc <\"ddd\"\"dxx\"|ahaha \"ihihi |ufufu\" ; ohoho gegee&&hogehogeo >ihihi")
	for key, statement := range result {
		fmt.Printf("%d :\n", key)
		for i, arg := range statement.argv {
			fmt.Printf("  [%d]='%s'\n", i, arg)
		}
		for i := 0; i < 3; i++ {
			fmt.Printf("  fd=%d:'%s'\n", i, statement.redirect[i].path)
		}
		fmt.Printf("  Terminator='%s'\n", statement.term)
	}
}
