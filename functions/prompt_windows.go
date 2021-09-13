package functions

import (
	"github.com/nyaosorg/go-windows-su"
)

func isElevated() bool {
	val, _ := su.IsElevated()
	return val
}
