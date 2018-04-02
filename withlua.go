// +build !nolua

package main

import (
	"github.com/zetamatta/nyagos/mainl"
)

func switchMain() error {
	return mainl.Main()
}
