//go:build !windows
// +build !windows

package nodos

func osDateLayout() (string, error) {
	return "Jan.02,2006", nil
}
