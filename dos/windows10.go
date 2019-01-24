package dos

import (
	"golang.org/x/sys/windows"
)

type ModeOp interface {
	Op(mode uint32) uint32
}

type ModeReset uint32

func (this ModeReset) Op(mode uint32) uint32 {
	return mode &^ uint32(this)
}

type ModeSet uint32

func (this ModeSet) Op(mode uint32) uint32 {
	return mode | uint32(this)
}

func ChangeConsoleMode(console windows.Handle, ops ...ModeOp) (func(), error) {
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

const enableVirtualTerminalProcessing uint32 = 0x0004

// EnableStdoutVirtualTerminalProcessing enables Windows10's native ESCAPE SEQUENCE support on STDOUT
func EnableStdoutVirtualTerminalProcessing() (func(), error) {
	return ChangeConsoleMode(windows.Stdout, ModeSet(enableVirtualTerminalProcessing))
}

// EnableStderrVirtualTerminalProcessing enables Windows10's native ESCAPE SEQUENCE support on STDERR
func EnableStderrVirtualTerminalProcessing() (func(), error) {
	return ChangeConsoleMode(windows.Stderr, ModeSet(enableVirtualTerminalProcessing))
}
