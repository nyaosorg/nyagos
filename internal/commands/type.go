package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nyaosorg/go-windows-mbcs"
)

func cat(ctx context.Context, r io.Reader, w io.Writer) error {
	scanner := mbcs.NewFilter(r, mbcs.ConsoleCP())
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return err
		}
		text := scanner.Text()
		text = strings.Replace(text, "\xEF\xBB\xBF", "", 1)
		fmt.Fprintln(w, text)
	}
	if err := scanner.Err(); err != io.EOF {
		return err
	}
	return nil
}

func cmdType(ctx context.Context, cmd Param) (int, error) {
	args := cmd.Args()
	if len(args) <= 1 {
		return 1, errors.New("The syntax of the command is incorrect.")
	}
	for _, arg1 := range args[1:] {
		r, err := os.Open(arg1)
		if err != nil {
			return 1, err
		}
		stat1, err := r.Stat()
		if err != nil {
			r.Close()
			return 2, err
		}
		if stat1.IsDir() {
			r.Close()
			return 3, fmt.Errorf("%s: Permission denied", arg1)
		}
		err = cat(ctx, r, cmd.Out())
		r.Close()
		if err != nil {
			return 0, err
		}
	}
	return 0, nil
}
