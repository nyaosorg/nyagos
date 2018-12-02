package functions

import (
	"github.com/zetamatta/nyagos/dos"
)

func isElevated() bool {
	val, _ := dos.IsElevated()
	return val
}
