package alias

import (
	"context"
	"strings"

	"github.com/nyaosorg/nyagos/completion"
	"github.com/nyaosorg/nyagos/shell"
)

var dbg = false

type callableT interface {
	String() string
	Call(ctx context.Context, cmd *shell.Cmd) (int, error)
}

// Func is the type for string-type alias. It has a Call method
type Func struct {
	BaseStr string
}

// New is the constructor for Func
func New(baseStr string) *Func {
	return &Func{baseStr}
}

// String is the method to support fmt.Stringer
func (f *Func) String() string {
	return f.BaseStr
}

var LineFilter = func(ctx context.Context, s string) string { return s }

// Call is the method to support callableT and it calls the alias-function.
func (f *Func) Call(ctx context.Context, cmd *shell.Cmd) (next int, err error) {
	next, err = cmd.Interpret(ctx, LineFilter(ctx, ExpandMacro(f.BaseStr, cmd.Args(), cmd.RawArgs())))
	return
}

// Table is the ALL ALIAS table !
var Table = map[string]callableT{}

// AllNames returns all-alias names for completion
func AllNames(ctx context.Context) ([]completion.Element, error) {
	names := make([]completion.Element, 0, len(Table))
	for name1 := range Table {
		names = append(names, completion.Element1(name1))
	}
	return names, nil
}

var nextHook shell.HookT

func hook(ctx context.Context, cmd *shell.Cmd) (int, bool, error) {
	lowerName := strings.ToLower(cmd.Arg(0))
	callee, ok := Table[lowerName]
	if !ok {
		return nextHook(ctx, cmd)
	}
	// Do not refer same name as alias.
	newcmd := *cmd
	newcmd.LineHook = func(_ctx context.Context, _cmd *shell.Cmd) (int, bool, error) {
		if strings.EqualFold(_cmd.Arg(0), lowerName) {
			return nextHook(_ctx, _cmd)
		}
		return hook(_ctx, _cmd)
	}
	next, err := callee.Call(ctx, &newcmd)
	return next, true, err
}

// Init is the package initializer which inserts hook-function into shell package.
func Init() {
	nextHook = shell.SetHook(hook)
}
