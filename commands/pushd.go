package commands

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"../cpath"
	"../dos"
)

var dirstack = make([]string, 0, 20)

const (
	NO_DIRSTACK = 2
	GETWD_FAIL  = 3
)

func cmd_dirs(ctx context.Context, cmd *exec.Cmd) (int, error) {
	wd, err := cpath.Getwd()
	if err != nil {
		return GETWD_FAIL, err
	}
	fmt.Fprint(cmd.Stdout, wd)
	for i := len(dirstack) - 1; i >= 0; i-- {
		fmt.Fprint(cmd.Stdout, " ", dirstack[i])
	}
	fmt.Fprintln(cmd.Stdout)
	return 0, nil
}

func cmd_popd(ctx context.Context, cmd *exec.Cmd) (int, error) {
	if len(dirstack) <= 0 {
		return NO_DIRSTACK, errors.New("popd: directory stack empty.")
	}
	err := dos.Chdir(dirstack[len(dirstack)-1])
	if err != nil {
		return CHDIR_FAIL, err
	}
	dirstack = dirstack[:len(dirstack)-1]
	return cmd_dirs(ctx, cmd)
}

func cmd_pushd(ctx context.Context, cmd *exec.Cmd) (int, error) {
	wd, err := cpath.Getwd()
	if err != nil {
		return GETWD_FAIL, err
	}
	if len(cmd.Args) >= 2 {
		dirstack = append(dirstack, wd)
		err := dos.Chdir(cmd.Args[1])
		if err != nil {
			return CHDIR_FAIL, err
		}
	} else {
		if len(dirstack) <= 0 {
			return NO_DIRSTACK, errors.New("pushd: directory stack empty.")
		}
		err := dos.Chdir(dirstack[len(dirstack)-1])
		if err != nil {
			return CHDIR_FAIL, err
		}
		dirstack[len(dirstack)-1] = wd
	}
	return cmd_dirs(ctx, cmd)
}
