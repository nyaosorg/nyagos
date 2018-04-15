package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/zetamatta/go-mbcs"
)

func cat(ctx context.Context, r io.Reader, w io.Writer) bool {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if done := ctx.Done(); done != nil {
			select {
			case <-done:
				return false
			default:

			}
		}
		line := scanner.Bytes()
		var text string
		if utf8.Valid(line) {
			text = string(line)
		} else {
			var err error
			text, err = mbcs.AtoU(line)
			if err != nil {
				text = err.Error()
			}
		}
		text = strings.Replace(text, "\xEF\xBB\xBF", "", 1)
		fmt.Fprintln(w, text)
	}
	return true
}

func cmdType(ctx context.Context, cmd Param) (int, error) {
	if len(cmd.Args()) <= 1 {
		cat(ctx, cmd.In(), cmd.Out())
	} else {
		for _, arg1 := range cmd.Args()[1:] {
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
			cont := cat(ctx, r, cmd.Out())
			r.Close()
			if !cont {
				return 0, nil
			}
		}
	}
	return 0, nil
}
