package nodos

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

func ChangeConsoleMode(console Handle, ops ...ModeOp) (func(), error) {
	return changeConsoleMode(console, ops...)
}

func DisableCtrlC() (func(), error) {
	return disableCtrlC()
}
