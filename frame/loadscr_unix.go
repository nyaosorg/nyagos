//go:build !windows
// +build !windows

package frame

import (
	"os"
	"path/filepath"
)

var appdatapath_ string

func appDataDir() string {
	if appdatapath_ == "" {
		appdatapath_ = filepath.Join(os.Getenv("HOME"), ".nyaos_org")
		os.Mkdir(appdatapath_, 0700)
	}
	return appdatapath_
}
