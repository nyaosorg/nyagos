package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

var dirstack = make([]string, 0, 20)

const (
	noDirStack = 2
	getwdFail  = 3
)

func cmdDirs(ctx context.Context, cmd Param) (int, error) {
	wd, err := os.Getwd()
	if err != nil {
		return getwdFail, err
	}
	io.WriteString(cmd.Out(), wd)
	for i := len(dirstack) - 1; i >= 0; i-- {
		fmt.Fprint(cmd.Out(), " ", dirstack[i])
	}
	fmt.Fprintln(cmd.Out())
	return 0, nil
}

func cmdPopd(ctx context.Context, cmd Param) (int, error) {
	if len(dirstack) <= 0 {
		return noDirStack, errors.New("popd: directory stack empty")
	}
	err := nodos.Chdir(dirstack[len(dirstack)-1])
	if err != nil {
		return errnoChdirFail, err
	}
	dirstack = dirstack[:len(dirstack)-1]
	return cmdDirs(ctx, cmd)
}

func cmdPushd(ctx context.Context, cmd Param) (int, error) {
	wd, err := os.Getwd()
	if err != nil {
		return getwdFail, err
	}
	if len(cmd.Args()) >= 2 {
		dirstack = append(dirstack, wd)
		_, err := cmdCdSub(cmd.Arg(1))
		if err != nil {
			return errnoChdirFail, err
		}
	} else {
		if len(dirstack) <= 0 {
			return noDirStack, errors.New("pushd: directory stack empty")
		}
		err := nodos.Chdir(dirstack[len(dirstack)-1])
		if err != nil {
			return errnoChdirFail, err
		}
		dirstack[len(dirstack)-1] = wd
	}
	return cmdDirs(ctx, cmd)
}
