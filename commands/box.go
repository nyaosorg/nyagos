package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/zetamatta/go-box"

	"github.com/zetamatta/nyagos/readline"
)

func cmdBox(ctx context.Context, cmd Param) (int, error) {
	data, err := ioutil.ReadAll(cmd.In())
	if err != nil {
		fmt.Fprintln(cmd.Err(), err.Error())
		return 1, err
	}
	list := strings.Split(string(data), "\n")
	if len(list) == 0 {
		return 1, nil
	}
	for i := 0; i < len(list); i++ {
		list[i] = strings.TrimSpace(list[i])
	}
	result := box.Choice(
		list,
		readline.Console)
	fmt.Fprintln(readline.Console)
	fmt.Fprintln(cmd.Out(), result)
	return 0, nil
}
