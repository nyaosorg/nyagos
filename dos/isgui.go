package dos

import (
	"io/ioutil"

	"github.com/H5eye/go-pefile"
)

func isGui(fname string) (bool, error) {
	bin, err := ioutil.ReadFile(fname)
	if err != nil {
		return false, err
	}
	pe, err := pefile.Parse(bin)
	if err != nil {
		return false, err
	}
	opt := pe.OptionalHeader
	return (opt != nil && opt.Subsystem == pefile.IMAGE_SUBSYSTEM_WINDOWS_GUI), nil
}

// IsGui returns true if fname is Windows GUI Application
func IsGui(fname string) bool {
	if fname == "" {
		return false
	}
	result, err := isGui(fname)
	if err != nil {
		return false
	}
	return result
}
