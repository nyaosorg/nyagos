package alias

import (
	"context"
	"strings"

	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/shell"
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

// Call is the method to support callableT and it calls the alias-function.
func (f *Func) Call(ctx context.Context, cmd *shell.Cmd) (next int, err error) {
	next, err = cmd.Interpret(ctx, ExpandMacro(f.BaseStr, cmd.Args(), cmd.RawArgs()))
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

type noAliasT string

func hook(ctx context.Context, cmd *shell.Cmd) (int, bool, error) {
	lowerName := strings.ToLower(cmd.Arg(0))
	ctxKey := noAliasT(lowerName)
	if ctx.Value(ctxKey) != nil {
		return nextHook(ctx, cmd)
	}
	callee, ok := Table[lowerName]
	if !ok {
		return nextHook(ctx, cmd)
	}
	next, err := callee.Call(context.WithValue(ctx, ctxKey, true), cmd)
	return next, true, err
}

// Init is the package initializer which inserts hook-function into shell package.
func Init() {
	nextHook = shell.SetHook(hook)
}
