// +build !windows

package nodos

func setConsoleExeIcon() (func(bool), error) {
	return func(bool) {}, nil
}
