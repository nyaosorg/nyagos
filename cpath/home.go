package cpath

import (
	"os/user"
	"path/filepath"
	"strings"
)

// C:\users\name\foo\bar -> ~\foo\bar
func ReplaceHomeToTilde(wd string) string {
	my, err := user.Current()
	if err != nil {
		return ""
	}
	home := my.HomeDir
	homeLen := len(home)
	if len(wd) >= homeLen && strings.EqualFold(home, wd[0:homeLen]) {
		wd = "~" + wd[homeLen:]
	}
	return wd
}

// C:\users\name\foo\bar -> ~/foo/bar
func ReplaceHomeToTildeSlash(wd string) string {
	return filepath.ToSlash(ReplaceHomeToTilde(wd))
}
