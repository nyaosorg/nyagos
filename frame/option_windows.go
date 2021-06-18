package frame

import (
	"fmt"
	"os"
	"strings"

	"github.com/zetamatta/go-windows-netresource"
	"github.com/zetamatta/go-windows-subst"
)

func optionNetUse(arg string) {
	piece := strings.SplitN(arg, "=", 2)
	if len(piece) >= 2 {
		_, err := netresource.NetUse(piece[0], piece[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "--netuse: %s: %s\n", arg, err.Error())
		}
	}
}

func optionSubst(arg string) {
	piece := strings.SplitN(arg, "=", 2)
	if len(piece) >= 2 {
		subst.Define(piece[0], piece[1])
	}
}
