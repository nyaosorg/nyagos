package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zetamatta/go-box"

	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

func cmd_box(ctx context.Context, cmd *shell.Cmd) (int, error) {
	data, err := ioutil.ReadAll(cmd.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1, err
	}
	list := strings.Split(string(data), "\n")
	for i := 0; i < len(list); i++ {
		list[i] = strings.TrimSpace(list[i])
	}
	result := box.Choice(
		list,
		readline.Console)
	fmt.Fprintln(cmd.Stdout, result)
	return 0, nil
}
