package dos

import "os"
import "strings"

// Get %HOME% || %USERPROFILE% || ""
func GetHome() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}

func ReplaceHomeToTilde(wd string) string {
	home := GetHome()
	homeLen := len(home)
	if len(wd) >= homeLen && strings.EqualFold(home, wd[0:homeLen]) {
		wd = "~" + wd[homeLen:]
	}
	return wd
}

func ReplaceHomeToTildeSlash(wd string) string {
	return strings.Replace(ReplaceHomeToTilde(wd), "\\", "/", -1)
}
