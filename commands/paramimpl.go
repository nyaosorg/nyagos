package commands

import (
	"context"

	"github.com/zetamatta/nyagos/shell"
)

type paramimp_t struct{ *shell.Cmd }

func (this *paramimp_t) Args() []string    { return this.Cmd.Args }
func (this *paramimp_t) RawArgs() []string { return this.Cmd.RawArgs }
func (this *paramimp_t) Spawn(ctx context.Context, args, rawargs []string) (int, error) {
	subCmd, err := this.Clone()
	if err != nil {
		return -1, err
	}
	subCmd.Args = args
	subCmd.RawArgs = rawargs
	return subCmd.SpawnvpContext(ctx)
}

func Exec(ctx context.Context, cmd *shell.Cmd) (int, bool, error) {
	return exec(ctx, &paramimp_t{cmd})
}
