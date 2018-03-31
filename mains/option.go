package mains

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/zetamatta/nyagos/lua"
	"github.com/zetamatta/nyagos/shell"
)

var optionNorc = false

func optionParse(sh *shell.Shell, L lua.Lua) (func() error, error) {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg1 := args[i]
		if arg1 == "-k" {
			i++
			if i >= len(args) {
				return nil, errors.New("-k: requires parameters")
			}
			return func() error {
				sh.Interpret(args[i])
				return nil
			}, nil
		} else if arg1 == "-c" {
			i++
			if i >= len(args) {
				return nil, errors.New("-c: requires parameters")
			}
			return func() error {
				sh.Interpret(args[i])
				return io.EOF
			}, nil
		} else if arg1 == "-b" {
			i++
			if i >= len(args) {
				return nil, errors.New("-b: requires parameters")
			}
			data, err := base64.StdEncoding.DecodeString(args[i])
			if err != nil {
				return nil, err
			}
			text := string(data)
			return func() error {
				sh.Interpret(text)
				return io.EOF
			}, nil
		} else if arg1 == "-f" {
			i++
			if i >= len(args) {
				return nil, errors.New("-f: requires parameters")
			}
			if strings.HasSuffix(strings.ToLower(args[i]), ".lua") {
				// lua script
				return func() error {
					setLuaArg(L, args[i:])
					_, err := runLua(sh, L, args[i])
					if err != nil {
						return err
					} else {
						return io.EOF
					}
				}, nil
			} else {
				return func() error {
					// command script
					fd, fd_err := os.Open(args[i])
					if fd_err != nil {
						return fmt.Errorf("%s: %s\n", args[i], fd_err.Error())
					}
					scanner := bufio.NewScanner(fd)
					for scanner.Scan() {
						sh.Interpret(scanner.Text())
					}
					fd.Close()
					return io.EOF
				}, nil
			}
		} else if arg1 == "-e" {
			i++
			if i >= len(args) {
				return nil, errors.New("-e: requires parameters")
			}
			return func() error {
				err := L.LoadString(args[i])
				if err != nil {
					return err
				}
				setLuaArg(L, args[i:])
				L.Call(0, 0)
				return io.EOF
			}, nil
		} else if arg1 == "--norc" {
			optionNorc = true
		} else if arg1 == "--lua-file" {
			i++
			if i >= len(args) {
				return nil, errors.New("--lua-file: requires parameters")
			}
			return func() error {
				setLuaArg(L, args[i:])
				_, err := runLua(sh, L, args[i])
				if err != nil {
					return err
				} else {
					return io.EOF
				}
			}, nil
		}
	}
	return nil, nil
}
