package history

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"

	"../readline"
)

func atoi_(reader io.Reader) (int, int) {
	n := 0
	count, err := fmt.Fscanf(reader, "%d", &n)
	if err == nil {
		return n, count
	} else {
		return 0, 0
	}
}

var Mark = "!"

var DisableMarks = "\"'"

func Replace(line string) (string, bool) {
	var mark rune
	for _, c := range Mark {
		mark = c
		break
	}

	var buffer bytes.Buffer
	isReplaced := false
	reader := strings.NewReader(line)
	history_count := readline.DefaultEditor.HistoryLen()

	quotedChar := '\000'

	for reader.Len() > 0 {
		ch, _, _ := reader.ReadRune()
		if quotedChar == '\000' && strings.IndexRune(DisableMarks, ch) >= 0 {
			quotedChar = ch
			buffer.WriteRune(ch)
			continue
		} else if ch == quotedChar {
			quotedChar = '\000'
			buffer.WriteRune(ch)
			continue
		}
		if ch != mark || reader.Len() <= 0 || quotedChar != '\000' {
			buffer.WriteRune(ch)
			continue
		}
		ch, _, _ = reader.ReadRune()
		if n := strings.IndexRune("^$:*", ch); n >= 0 {
			reader.UnreadRune()
			if history_count >= 2 {
				insertHistory(&buffer, reader, history_count-2)
				isReplaced = true
			}
			continue
		}
		if ch == mark { // !!
			if history_count >= 2 {
				insertHistory(&buffer, reader, history_count-2)
				isReplaced = true
				continue
			} else {
				buffer.WriteRune(mark)
				continue
			}
		}
		if strings.IndexRune("0123456789", ch) >= 0 { // !n
			reader.UnreadRune()
			backno, _ := atoi_(reader)
			backno = backno % history_count
			if 0 <= backno && backno < history_count {
				insertHistory(&buffer, reader, backno)
				isReplaced = true
			}
			continue
		}
		if ch == '-' && reader.Len() > 0 { // !-n
			if number, count := atoi_(reader); count > 0 {
				backno := history_count - number - 1
				for backno < 0 {
					backno += history_count
				}
				if 0 <= backno && backno < history_count {
					insertHistory(&buffer, reader, backno)
					isReplaced = true
				} else {
					buffer.WriteRune(mark)
					buffer.WriteString("-0")
				}
				continue
			} else {
				reader.UnreadRune() // next char of '-'
			}
		}
		if ch == '?' { // !?str?
			var seekStrBuf bytes.Buffer
			lastCharIsQuestionMark := false
			for reader.Len() > 0 {
				ch, _, _ := reader.ReadRune()
				if ch == '?' {
					lastCharIsQuestionMark = true
					break
				}
				seekStrBuf.WriteRune(ch)
			}
			seekStr := seekStrBuf.String()
			found := false
			for i := history_count - 2; i >= 0; i-- {
				if strings.Contains(readline.DefaultEditor.Histories[i].Line, seekStr) {
					buffer.WriteString(readline.DefaultEditor.Histories[i].Line)
					isReplaced = true
					found = true
					break
				}
			}
			if !found {
				buffer.WriteRune('?')
				buffer.WriteString(seekStr)
				if lastCharIsQuestionMark {
					buffer.WriteRune('?')
				}
			}
			continue
		}
		// !str
		var seekStrBuf bytes.Buffer
		seekStrBuf.WriteRune(ch)
		for reader.Len() > 0 {
			ch, _, _ := reader.ReadRune()
			if ch == ' ' {
				reader.UnreadRune()
				break
			}
			seekStrBuf.WriteRune(ch)
		}
		seekStr := seekStrBuf.String()
		found := false
		for i := history_count - 2; i >= 0; i-- {
			if strings.HasPrefix(readline.DefaultEditor.Histories[i].Line, seekStr) {
				buffer.WriteString(readline.DefaultEditor.Histories[i].Line)
				isReplaced = true
				found = true
				break
			}
		}
		if !found {
			buffer.WriteRune(mark)
			buffer.WriteRune(ch)
		}
	}
	result := readline.NewHistoryLine(buffer.String())
	if isReplaced {
		if history_count > 0 {
			readline.DefaultEditor.Histories[history_count-1] = result
		} else {
			readline.DefaultEditor.Histories = append(readline.DefaultEditor.Histories, result)
		}
	}
	return result.Line, isReplaced
}

func insertHistory(buffer *bytes.Buffer, reader *strings.Reader, historyNo int) {
	history1 := readline.DefaultEditor.Histories[historyNo]
	ch, siz, _ := reader.ReadRune()
	if siz > 0 && ch == '^' {
		if len(history1.Word) >= 2 {
			buffer.WriteString(history1.At(1))
		}
	} else if siz > 0 && ch == '$' {
		if len(history1.Word) >= 2 {
			buffer.WriteString(history1.At(-1))
		}
	} else if siz > 0 && ch == '*' {
		if len(history1.Word) >= 2 {
			buffer.WriteString(strings.Join(history1.Word[1:], " "))
		}
	} else if siz > 0 && ch == ':' {
		n, count := atoi_(reader)
		if count <= 0 {
			buffer.WriteRune(':')
		} else if n < len(history1.Word) {
			buffer.WriteString(history1.Word[n])
		}
	} else {
		if siz > 0 {
			reader.UnreadRune()
		}
		buffer.WriteString(history1.Line)
	}
}

func CmdHistory(ctx context.Context, cmd *exec.Cmd) (int, error) {
	var num int
	if len(cmd.Args) >= 2 {
		num64, err := strconv.ParseInt(cmd.Args[1], 0, 32)
		if err != nil {
			switch err.(type) {
			case *strconv.NumError:
				return 0, fmt.Errorf(
					"history: %s not a number", cmd.Args[1])
			default:
				return 0, err
			}
		}
		num = int(num64)
		if num < 0 {
			num = -num
		}
	} else {
		num = 10
	}
	start := 0
	if f, ok := cmd.Stdout.(*os.File); (!ok || isatty.IsTerminal(f.Fd())) &&
		readline.DefaultEditor.HistoryLen() > num {

		start = readline.DefaultEditor.HistoryLen() - num
	}
	for i, s := range readline.DefaultEditor.Histories[start:] {
		fmt.Fprintf(cmd.Stdout, "%3d : %-s\n", start+i, s.Line)
	}
	return 0, nil
}

const max_histories = 2000

func Save(path string) error {
	readline.DefaultEditor.ShrinkHistory()
	start := 0
	if readline.DefaultEditor.HistoryLen() > max_histories {
		start = readline.DefaultEditor.HistoryLen() - max_histories
	}
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	for i := start; i < len(readline.DefaultEditor.Histories); i++ {
		fmt.Fprintln(fd, readline.DefaultEditor.Histories[i].Line)
	}
	return nil
}

func Load(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		readline.DefaultEditor.Histories = append(readline.DefaultEditor.Histories, readline.NewHistoryLine(sc.Text()))
	}
	readline.DefaultEditor.ShrinkHistory()
	return nil
}
