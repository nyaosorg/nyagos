package frame

import (
	"strings"

	"github.com/nyaosorg/go-windows-subst"
	"github.com/zetamatta/go-windows-netresource"
)

func optionNetUse(arg string) {
	piece := strings.SplitN(arg, "=", 2)
	if len(piece) >= 2 {
		netresource.NetUse(piece[0], piece[1])
	}
}

func optionSubst(arg string) {
	piece := strings.SplitN(arg, "=", 2)
	if len(piece) >= 2 {
		subst.Define(piece[0], piece[1])
	}
}
