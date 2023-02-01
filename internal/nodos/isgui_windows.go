package nodos

import (
	"debug/pe"
)

func isGui(fname string) (bool, error) {
	file, err := pe.Open(fname)
	if err != nil {
		return false, err
	}
	defer file.Close()
	opt := file.OptionalHeader
	if opt32, ok := opt.(*pe.OptionalHeader32); ok {
		return opt32.Subsystem == pe.IMAGE_SUBSYSTEM_WINDOWS_GUI, nil
	} else if opt64, ok := opt.(*pe.OptionalHeader64); ok {
		return opt64.Subsystem == pe.IMAGE_SUBSYSTEM_WINDOWS_GUI, nil
	} else {
		return false, nil
	}
}
