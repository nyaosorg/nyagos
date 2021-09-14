package nodos

import (
	"fmt"
	"golang.org/x/sys/windows"
	// "github.com/nyaosorg/go-windows-dbg"
)

type Handle = windows.Handle

func changeConsoleMode(console Handle, ops ...ModeOp) (func(), error) {
	var mode uint32
	err := windows.GetConsoleMode(console, &mode)
	if err != nil {
		return func() {}, fmt.Errorf("windows.GetConsoleMode: %w", err)
	}
	restore := func() { windows.SetConsoleMode(console, mode) }

	if len(ops) > 0 {
		newMode := mode
		for _, op1 := range ops {
			newMode = op1.Op(newMode)
		}
		// dbg.Printf("windows.SetConsoleMode(%v,0x%X): %w", console, newMode, err)
		err = windows.SetConsoleMode(console, newMode)
		if err != nil {
			err = fmt.Errorf("windows.SetConsoleMode(%v,0x%X): %w", console, newMode, err)
		}
	}
	return restore, err
}

func enableProcessInput() (func(), error) {
	// If ENABLE_ECHO_INPUT and ENABLE_PROCESSED_INPUT were set,
	// but ENABLE_LINE_INPUT was reset, SetConsoleMode would fail.
	return changeConsoleMode(windows.Stdin,
		ModeSet(windows.ENABLE_ECHO_INPUT),
		ModeSet(windows.ENABLE_LINE_INPUT),
		ModeSet(windows.ENABLE_PROCESSED_INPUT))
}
