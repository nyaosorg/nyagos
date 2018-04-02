// +build nolua

package main

import (
	"github.com/zetamatta/nyagos/mains"
)

func switchMain() error {
	return mains.Main()
}
