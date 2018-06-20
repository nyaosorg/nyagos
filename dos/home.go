package dos

import (
	"os"
	"path/filepath"
	"strings"
)

// GetHome get %HOME% or %USERPROFILE%
func GetHome() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}

// ReplaceHomeToTilde replaces path like C:\users\name\foo\bar -> ~\foo\bar
func ReplaceHomeToTilde(wd string) string {
	home := GetHome()
	homeLen := len(home)
	if len(wd) >= homeLen && strings.EqualFold(home, wd[0:homeLen]) {
		wd = "~" + wd[homeLen:]
	}
	return wd
}

// ReplaceHomeToTildeSlash replaces path like C:\users\name\foo\bar -> ~/foo/bar
func ReplaceHomeToTildeSlash(wd string) string {
	return filepath.ToSlash(ReplaceHomeToTilde(wd))
}
