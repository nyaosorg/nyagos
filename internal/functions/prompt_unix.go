//go:build !windows
// +build !windows

package functions

func isElevated() bool {
	return false
}
