//go:build !windows
// +build !windows

package nodos

func isGui(fname string) (bool, error) {
	return false, nil
}
