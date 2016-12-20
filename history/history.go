package history

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"

	"../readline"
	"../text"
)

var Mark = "!"

var DisableMarks = "\"'"

type IHistory interface {
	Len() int
	At(int) string
	Push(string)
	Replace(string)
}

func Replace(hisObj IHistory, line string) (string, bool) {
	var mark rune
	for _, c := range Mark {
		mark = c
		break
	}

	var buffer bytes.Buffer
	isReplaced := false
	reader := strings.NewReader(line)
	history_count := hisObj.Len()

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
				line := hisObj.At(history_count - 2)
				InsertHistory(&buffer, reader, line)
				isReplaced = true
			}
			continue
		}
		if ch == mark { // !!
			if history_count >= 2 {
				line := hisObj.At(history_count - 2)
				InsertHistory(&buffer, reader, line)
				isReplaced = true
				continue
			} else {
				buffer.WriteRune(mark)
				continue
			}
		}
		if strings.IndexRune("0123456789", ch) >= 0 { // !n
			reader.UnreadRune()
			var backno int
			fmt.Fscan(reader, &backno)
			backno = backno % history_count
			if 0 <= backno && backno < history_count {
				line := hisObj.At(backno)
				InsertHistory(&buffer, reader, line)
				isReplaced = true
			}
			continue
		}
		if ch == '-' && reader.Len() > 0 { // !-n
			var number int
			if _, err := fmt.Fscan(reader, &number); err == nil {
				backno := history_count - number - 1
				for backno < 0 {
					backno += history_count
				}
				if 0 <= backno && backno < history_count {
					line := hisObj.At(backno)
					InsertHistory(&buffer, reader, line)
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
				his1 := hisObj.At(i)
				if strings.Contains(his1, seekStr) {
					buffer.WriteString(his1)
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
			his1 := hisObj.At(i)
			if strings.HasPrefix(his1, seekStr) {
				buffer.WriteString(his1)
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
	result_line := buffer.String()
	if isReplaced {
		if history_count > 0 {
			hisObj.Replace(result_line)
		} else {
			hisObj.Replace(result_line)
		}
	}
	return result_line, isReplaced
}

func InsertHistory(buffer *bytes.Buffer, reader *strings.Reader, line string) {
	ch, siz, _ := reader.ReadRune()
	if siz > 0 && ch == '^' {
		if words := text.SplitQ(line); len(words) >= 2 {
			buffer.WriteString(words[1])
		}
	} else if siz > 0 && ch == '$' {
		if words := text.SplitQ(line); len(words) >= 2 {
			buffer.WriteString(words[len(words)-1])
		}
	} else if siz > 0 && ch == '*' {
		if words := text.SplitQ(line); len(words) >= 2 {
			buffer.WriteString(strings.Join(words[1:], " "))
		}
	} else if siz > 0 && ch == ':' {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			buffer.WriteRune(':')
		} else if words := text.SplitQ(line); n < len(words) {
			buffer.WriteString(words[n])
		}
	} else {
		if siz > 0 {
			reader.UnreadRune()
		}
		buffer.WriteString(line)
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

func Load(path string, hisObj IHistory) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	sc := bufio.NewScanner(fd)
	list := make([]string, 0, 2000)
	hash := make(map[string]int)
	for sc.Scan() {
		line := sc.Text()
		if lnum, ok := hash[line]; ok {
			list[lnum] = ""
		}
		hash[line] = len(list)
		list = append(list, line)
	}
	for _, line := range list {
		if line != "" {
			hisObj.Push(line)
		}
	}
	return nil
}
