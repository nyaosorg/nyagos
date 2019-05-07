// +build !windows

package commands

func findBatch(name string) (string, bool) {
	return name, true
}
