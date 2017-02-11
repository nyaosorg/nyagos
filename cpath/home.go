package cpath

import (
	"os"
	"path/filepath"
	"strings"
)

// Get %HOME% || %USERPROFILE%
func GetHome() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}

// C:\users\name\foo\bar -> ~\foo\bar
func ReplaceHomeToTilde(wd string) string {
	home := GetHome()
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
