package history

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mattn/go-isatty"
)

type Dumper interface {
	Len() int
	DumpAt(int) string
}

type Param interface {
	Arg(int) string
	Args() []string
	Out() io.Writer
	Err() io.Writer
}

func CmdHistory(ctx context.Context, cmd Param, historyObj Dumper) (int, error) {
	if ctx == nil {
		fmt.Fprintln(cmd.Err(), "history not found (case1)")
		return 1, nil
	}
	var num int
	if len(cmd.Args()) >= 2 {
		num64, err := strconv.ParseInt(cmd.Arg(1), 0, 32)
		if err != nil {
			switch err.(type) {
			case *strconv.NumError:
				return 0, fmt.Errorf(
					"history: %s not a number", cmd.Arg(1))
			default:
				return 0, err
			}
		}
		num = int(num64)
		if num < 0 {
			num = -num
		}
	} else {
		num = 10
	}
	start := 0

	if f, ok := cmd.Out().(*os.File); ok && isatty.IsTerminal(f.Fd()) && historyObj.Len() > num {
		start = historyObj.Len() - num
	}
	for i := start; i < historyObj.Len(); i++ {
		fmt.Fprintf(cmd.Out(), "%4d  %s\n", i, historyObj.DumpAt(i))
	}
	return 0, nil
}
