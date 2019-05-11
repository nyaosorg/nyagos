package frame

import (
	"os"
	"strings"

	"github.com/zetamatta/nyagos/dos"
)

func tryOsOption(arg string) bool {
	const netuse_opt = "--netuse="
	const chdir_opt = "--chdir="

	if strings.HasPrefix(arg, netuse_opt) {
		arg = arg[len(netuse_opt):]
		piece := strings.SplitN(arg, "=", 2)
		if len(piece) != 2 {
			return false
		}
		dos.NetUse(piece[0], piece[1])
		return true
	} else if strings.HasPrefix(arg, chdir_opt) {
		os.Chdir(arg[len(chdir_opt):])
		return true
	} else {
		return false
	}
}
