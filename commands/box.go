package commands

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/zetamatta/go-box/v2"
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

	console := bufio.NewWriter(cmd.Term())
	result := box.Choice(
		list,
		console)
	fmt.Fprintln(console)
	console.Flush()
	fmt.Fprintln(cmd.Out(), result)
	return 0, nil
}
