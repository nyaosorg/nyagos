package frame

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/atotto/clipboard"
	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-windows-consoleicon"

	"github.com/nyaosorg/nyagos/internal/config"
	"github.com/nyaosorg/nyagos/internal/history"
	"github.com/nyaosorg/nyagos/internal/shell"

	"github.com/nyaosorg/go-readline-skk"
)

type CmdStreamConsole struct {
	shell.CmdSeeker
	DoPrompt func(io.Writer) (int, error)
	History  *history.Container
	Editor   *readline.Editor
	HistPath string
}

type OSClipboard struct{}

func (OSClipboard) Read() (string, error) {
	return clipboard.ReadAll()
}

func (OSClipboard) Write(s string) error {
	return clipboard.WriteAll(s)
}

func NewCmdStreamConsole(doPrompt func(io.Writer) (int, error)) *CmdStreamConsole {
	history1 := &history.Container{}
	stream := &CmdStreamConsole{
		History: history1,
		Editor: &readline.Editor{
			History:        history1,
			PromptWriter:   doPrompt,
			Writer:         colorable.NewColorableStdout(),
			HistoryCycling: true,
		},
		HistPath: filepath.Join(appDataDir(), "nyagos.history"),
		CmdSeeker: shell.CmdSeeker{
			PlainHistory: []string{},
			Pointer:      -1,
		},
	}
	history1.Load(stream.HistPath)
	history1.Save(stream.HistPath)
	return stream
}

func (stream *CmdStreamConsole) LazySetup() {
	if config.AccessClipboard {
		stream.Editor.Clipboard = OSClipboard{}
	}
	if _, ok := os.LookupEnv("NO_COLOR"); !ok {
		stream.Editor.Highlight = []readline.Highlight{
			// Options -> Dark Yellow
			{Pattern: regexp.MustCompile(` \-[-\w]+`), Sequence: "\x1B[0;33m"},
			// Backquotation -> Red
			{Pattern: regexp.MustCompile("`[^`]*`|`[^`]*$"), Sequence: "\x1B[0;31;1m"},
			// & | < > ; -> Green
			{Pattern: regexp.MustCompile(`[&\|<>]| ;`), Sequence: "\x1B[0;32;1m"},
			// Double quotation -> Magenta
			{Pattern: regexp.MustCompile(`"([^"]*\\")*[^"]*$|"([^"]*\\")*[^"]*"`), Sequence: "\x1B[0;35;1m"},
			// Enviroment variable -> Cyan
			{Pattern: regexp.MustCompile(`%[A-Za-z][A-Za-z0-9]*$|%[A-Za-z][A-Za-z0-9]*%`), Sequence: "\x1B[0;36;1m"},
			// Control characters -> Blue
			{Pattern: regexp.MustCompile("[\x00-\x1F]+"), Sequence: "\x1B[0;34;1m"},
			// Wide space -> Background Red
			{Pattern: regexp.MustCompile("\u3000"), Sequence: "\x1B[0;41m"},

			skk.WhiteMarkerHighlight,
			skk.BlackMarkerHighlight,
		}
		stream.Editor.ResetColor = "\x1B[0m"
		stream.Editor.DefaultColor = "\x1B[0;1m"

		if config.OptionPredictColor {
			stream.Editor.PredictColor = config.PredictColor
		}
	}
}

func (stream *CmdStreamConsole) DisableHistory(value bool) bool {
	return stream.History.IgnorePush(value)
}

// endsWithSep returns
//
//	false when line does not end with `^`
//	true when line ends with `^`
//	false when line ends with `^^`
//	true when line ends with `^^^`
func endsWithSep(line []byte, contMark byte) bool {
	markCount := 0
	for len(line) > 0 && line[len(line)-1] == contMark {
		markCount++
		line = line[:len(line)-1]
	}
	return markCount%2 != 0
}

func (stream *CmdStreamConsole) readLineContinued(ctx context.Context) (string, error) {
	continued := false
	originalPrompt := os.Getenv("PROMPT")
	buffer := make([]byte, 0, 256)
	for {
		line, err := stream.Editor.ReadLine(ctx)
		stream.Editor.Default = ""
		buffer = append(buffer, line...)
		if err != nil || !endsWithSep(buffer, '^') {
			if continued {
				os.Setenv("PROMPT", originalPrompt)
				// stream.Editor.Coloring.(*_Coloring).defaultBits &^= quotedBit
			}
			return string(buffer), err
		}
		buffer = buffer[:len(buffer)-1]
		buffer = append(buffer, '\r', '\n')
		continued = true
		os.Setenv("PROMPT", "> ")
		// if bytes.Count(buffer, []byte{'"'})%2 != 0 {
		// stream.Editor.Coloring.(*_Coloring).defaultBits |= quotedBit
		// }
	}
}

func (stream *CmdStreamConsole) ReadLine(ctx context.Context) (string, error) {
	if stream.Pointer >= 0 {
		if stream.Pointer < len(stream.PlainHistory) {
			stream.Pointer++
			return stream.PlainHistory[stream.Pointer-1], nil
		}
		stream.Pointer = -1
	}
	var line string
	var err error
	for {
		disabler := colorable.EnableColorsStdout(nil)
		clean, err2 := consoleicon.SetFromExe()
		line, err = stream.readLineContinued(ctx)
		if err2 == nil {
			clean(false)
		}
		disabler()
		if err != nil {
			return line, err
		}
		var isReplaced bool
		line, isReplaced, err = stream.History.Replace(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line != "" {
			break
		}
	}
	row := history.NewHistoryLine(line)
	stream.History.PushLine(row)
	fd, err := os.OpenFile(stream.HistPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err == nil {
		fmt.Fprintln(fd, row.String())
		fd.Close()
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	stream.PlainHistory = append(stream.PlainHistory, line)
	return line, err
}
