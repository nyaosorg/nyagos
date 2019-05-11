package frame

import (
	"strings"

	"github.com/zetamatta/nyagos/dos"
)

func optionNetUse(arg string) {
	piece := strings.SplitN(arg, "=", 2)
	if len(piece) >= 2 {
		dos.NetUse(piece[0], piece[1])
	}
}
