package nodos

import (
	"golang.org/x/sys/windows"
)

type Handle = windows.Handle

func changeConsoleMode(console Handle, ops ...ModeOp) (func(), error) {
	var mode uint32
	err := windows.GetConsoleMode(console, &mode)
	if err != nil {
		return func() {}, err
	}
	restore := func() { windows.SetConsoleMode(console, mode) }

	if len(ops) > 0 {
		newMode := mode
		for _, op1 := range ops {
			newMode = op1.Op(newMode)
		}
		err = windows.SetConsoleMode(console, newMode)
	}
	return restore, err
}

func disableCtrlC() (func(), error) {
	return changeConsoleMode(windows.Stdin,
		ModeReset(windows.ENABLE_PROCESSED_INPUT))
}
