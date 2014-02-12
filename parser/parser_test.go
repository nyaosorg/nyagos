package parser

import "testing"
import "fmt"

func TestParser(t *testing.T) {
	text := "aaaa bbbb \"\" cccc <\"ddd\"\"ddd\"|ahaha \"ihihi |ufufu\" ; ohoho gegee&&hogehogeo >ihihi"
	fmt.Println(text)
	result1 := Parse1(text)
	result2 := Parse2(result1)
	for i,st := range result2 {
		fmt.Printf("pipeline-%d:\n",i)
		for _,stsub := range st {
			fmt.Printf("  %s\n",stsub.String())
		}
	}
}
