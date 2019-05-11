package frame

import (
	"github.com/zetamatta/nyagos/dos"
	"strings"
)

func tryOsOption(arg string) bool {
	const option = "--netuse="

	if !strings.HasPrefix(arg, option) {
		return false
	}
	arg = arg[len(option):]
	piece := strings.SplitN(arg, "=", 2)
	if len(piece) != 2 {
		return false
	}
	dos.NetUse(piece[0], piece[1])
	return true
}
