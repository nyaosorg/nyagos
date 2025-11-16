package functions

import (
	"io"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/nyagos/internal/shell"
)

// Param is the parameter type for nyagos.xxxxxx which uses stdin/stdout/stderr/colored-console.
type Param struct {
	Args []interface{}
	In   io.Reader
	Out  io.Writer
	Err  io.Writer
	Term io.Writer
	shell.History
	shell.Stream
	*readline.Editor
}
