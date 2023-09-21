package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/mattn/go-tty"

	"github.com/nyaosorg/go-windows-mbcs"
	"github.com/nyaosorg/nyagos/internal/nodos"
	"github.com/nyaosorg/nyagos/internal/textwidth"
)

var bold = false
var screenWidth int
var screenHeight int

func isTerminalIn(in io.Reader) bool {
	f, ok := in.(*os.File)
	if !ok {
		return false
	}
	return f.Fd() == os.Stdin.Fd()
}

func getkey() (rune, error) {
	tty1, err := tty.Open()
	if err != nil {
		return 0, err
	}
	defer tty1.Close()
	for {
		ch, err := tty1.ReadRune()
		if err != nil {
			return 0, err
		}
		if ch != 0 {
			return ch, nil
		}
	}
}

func splitLinesWithWidth(text string, screenWidth int) (lines []string) {
	var buffer strings.Builder
	w := 0
	ansiStrip := false
	for _, c := range text {
		if ansiStrip {
			if ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') {
				ansiStrip = false
			}
		} else if c == '\x1B' {
			ansiStrip = true
		} else {
			w1 := textwidth.RuneWidth(c)
			if w+w1 >= screenWidth {
				lines = append(lines, buffer.String())
				buffer.Reset()
				w = 0
			}
			w += w1
		}
		buffer.WriteRune(c)
	}
	if buffer.Len() > 0 {
		lines = append(lines, buffer.String())
	}
	if len(lines) <= 0 {
		lines = []string{""}
	}
	return
}

func more(r io.Reader, cmd Param) error {
	scanner := mbcs.NewFilter(r, mbcs.ConsoleCP())
	count := 0

	if f, ok := cmd.Out().(*os.File); !ok || !isatty.IsTerminal(f.Fd()) {
		screenHeight = math.MaxInt32
	}
	for scanner.Scan() {
		text := scanner.Text()
		lines := splitLinesWithWidth(text, screenWidth)
		for _, line := range lines {
			if count+1 >= screenHeight {
				io.WriteString(cmd.Err(), "more>")
				ch, err := getkey()
				if err != nil {
					return err
				}
				if ch == '\x03' {
					fmt.Fprintln(cmd.Err(), "^C")
					return io.EOF
				}
				io.WriteString(cmd.Err(), "\r     \b\b\b\b\b")
				if ch == 'q' {
					return io.EOF
				} else if ch == '\r' {
					count--
				} else {
					count = 0
				}
				if bold {
					io.WriteString(cmd.Out(), "\x1B[1m")
				}
			}
			fmt.Fprintln(cmd.Out(), line)
			count++
		}
	}
	return scanner.Err()
}

func cmdMore(ctx context.Context, cmd Param) (int, error) {
	count := 0

	tty1, err := tty.Open()
	if err != nil {
		return 1, err
	}
	screenWidth, screenHeight, err = tty1.Size()
	tty1.Close()
	if err != nil {
		return 1, err
	}
	for _, arg1 := range cmd.Args()[1:] {
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
		if err := more(r, cmd); err != nil {
			r.Close()
			if err == io.EOF {
				return 0, nil
			}
			return 1, err
		}
		r.Close()
		count++
	}
	if count <= 0 {
		if isTerminalIn(cmd.In()) {
			c, err := nodos.EnableProcessInput()
			if err != nil {
				return 1, err
			}
			defer c()
		}
		err := more(cmd.In(), cmd)
		if err != nil && err != io.EOF {
			return 1, err
		}
	}
	return 0, nil
}
