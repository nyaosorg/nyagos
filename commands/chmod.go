package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var rxOOO = regexp.MustCompile("^[0-7][0-7][0-7]$")

func getOct(s string) (int, error) {
	val := 0
	for _, r := range s {
		n := strings.IndexRune("01234567", r)
		if n < 0 {
			return 0, fmt.Errorf("%s: invalid permission str", s)
		}
		val = val*8 + n
	}
	return val, nil
}

var rxEqu = regexp.MustCompile(`^([aogu]+)([\-\+\=])([rwx]+)$`)

func cmd_chmod_(args []string) error {
	if len(args) < 2 {
		return errors.New("Usage: chmod ooo (files...)")
	}
	var f func(string) error
	if rxOOO.MatchString(args[0]) {
		val, err := getOct(args[0])
		if err != nil {
			return fmt.Errorf("%s: invalid permission str", args[0])
		}
		f = func(fname string) error {
			return os.Chmod(fname, os.FileMode(val))
		}
	} else if m := rxEqu.FindStringSubmatch(args[0]); m != nil {
		var basebits os.FileMode
		for _, r := range m[3] {
			switch r {
			case 'r':
				basebits |= 4
			case 'w':
				basebits |= 2
			case 'x':
				basebits |= 1
			}
		}
		var bits os.FileMode
		for _, r := range m[1] {
			switch r {
			case 'u':
				bits |= (basebits << 6)
			case 'g':
				bits |= (basebits << 3)
			case 'o':
				bits |= basebits
			case 'a':
				bits = (basebits << 6) | (basebits << 3) | basebits
			}
		}
		switch m[2] {
		case "+":
			f = func(fname string) error {
				stat, err := os.Stat(fname)
				if err != nil {
					return err
				}
				mod := stat.Mode()
				return os.Chmod(fname, mod|bits)
			}
		case "-":
			f = func(fname string) error {
				stat, err := os.Stat(fname)
				if err != nil {
					return err
				}
				mod := stat.Mode()
				return os.Chmod(fname, mod&^bits)
			}
		case "=":
			f = func(fname string) error {
				return os.Chmod(fname, bits)
			}
		}
	} else {
		return fmt.Errorf("%s: invalid permission str", args[0])
	}
	for _, fname := range args[1:] {
		if err := f(fname); err != nil {
			return fmt.Errorf("%s: %s", fname, err.Error())
		}
	}
	return nil
}

func cmdChmod(_ context.Context, cmd Param) (int, error) {
	if err := cmd_chmod_(cmd.Args()[1:]); err != nil {
		return 1, err
	} else {
		return 0, nil
	}

}
