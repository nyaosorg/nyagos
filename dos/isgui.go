package dos

import (
	"github.com/hillu/go-pefile"
	"io/ioutil"
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
	if opt != nil {
		if opt.Subsystem == pefile.IMAGE_SUBSYSTEM_WINDOWS_GUI {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

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
