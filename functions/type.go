package functions

import "io"

type Param struct {
	Args []interface{}
	In   io.Reader
	Out  io.Writer
	Err  io.Writer
	Term io.Writer
}
