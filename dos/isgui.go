package dos

import (
	"debug/pe"
)

const _IMAGE_SUBSYSTEM_WINDOWS_GUI = 2

func isGui(fname string) (bool, error) {
	file, err := pe.Open(fname)
	if err != nil {
		return false, err
	}
	defer file.Close()
	opt := file.OptionalHeader
	if opt32, ok := opt.(*pe.OptionalHeader32); ok {
		return opt32.Subsystem == _IMAGE_SUBSYSTEM_WINDOWS_GUI, nil
	} else if opt64, ok := opt.(*pe.OptionalHeader64); ok {
		return opt64.Subsystem == _IMAGE_SUBSYSTEM_WINDOWS_GUI, nil
	} else {
		return false, nil
	}
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
