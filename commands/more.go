package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/zetamatta/go-box"
	"github.com/zetamatta/go-getch"
	"github.com/zetamatta/go-mbcs"
	"github.com/zetamatta/nyagos/shell"
)

var ansiStrip = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")

var bold = false
var screenWidth int
var screenHeight int

func more(r io.Reader, cmd *shell.Cmd) bool {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
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
		width := runewidth.StringWidth(ansiStrip.ReplaceAllString(text, ""))
		lines := (width + screenWidth) / screenWidth
		for count+lines >= screenHeight {
			fmt.Fprint(cmd.Stderr, "more>")
			ch := getch.Rune()
			fmt.Fprint(cmd.Stderr, "\r     \b\b\b\b\b")
			if ch == 'q' {
				return false
			} else if ch == '\r' {
				count--
			} else {
				count = 0
			}
		}
		if bold {
			fmt.Fprint(cmd.Stdout, "\x1B[1m")
		}
		fmt.Fprintln(cmd.Stdout, text)
		count += lines
	}
	return true
}

func cmd_more(ctx context.Context, cmd *shell.Cmd) (int, error) {
	count := 0
	screenWidth, screenHeight = box.GetScreenBufferInfo().ViewSize()
	for _, arg1 := range os.Args[1:] {
		if arg1 == "-b" {
			bold = true
			continue
		} else if arg1 == "-h" {
			return 1, errors.New("more : Color-Unicoded more")
		}
		r, err := os.Open(arg1)
		if err != nil {
			return 1, err
		}
		if !more(r, cmd) {
			r.Close()
			return 0, nil
		}
		r.Close()
		count++
	}
	if count <= 0 {
		more(cmd.Stdin, cmd)
	}
	return 0, nil
}
