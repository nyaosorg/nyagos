package nodos

import (
	"debug/pe"
)

const imageSubsystemWindowsGui = 2

func isGui(fname string) (bool, error) {
	file, err := pe.Open(fname)
	if err != nil {
		return false, err
	}
	defer file.Close()
	opt := file.OptionalHeader
	if opt32, ok := opt.(*pe.OptionalHeader32); ok {
		return opt32.Subsystem == imageSubsystemWindowsGui, nil
	} else if opt64, ok := opt.(*pe.OptionalHeader64); ok {
		return opt64.Subsystem == imageSubsystemWindowsGui, nil
	} else {
		return false, nil
	}
}
