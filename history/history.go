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
	"time"

	"github.com/mattn/go-isatty"

	"../text"
)

var Mark = "!"

var DisableMarks = "\"'"

func (hisObj *Container) Replace(line string) (string, bool) {
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
			if history_count >= 1 {
				line := hisObj.At(history_count - 1)
				ExpandMacro(&buffer, reader, line)
				isReplaced = true
			}
			continue
		}
		if ch == mark { // !!
			if history_count >= 1 {
				line := hisObj.At(history_count - 1)
				ExpandMacro(&buffer, reader, line)
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
				ExpandMacro(&buffer, reader, line)
				isReplaced = true
			}
			continue
		}
		if ch == '-' && reader.Len() > 0 { // !-n
			var number int
			if _, err := fmt.Fscan(reader, &number); err == nil {
				backno := history_count - number
				for backno < 0 {
					backno += history_count
				}
				if 0 <= backno && backno < history_count {
					line := hisObj.At(backno)
					ExpandMacro(&buffer, reader, line)
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
			for i := history_count - 1; i >= 0; i-- {
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
		for i := history_count - 1; i >= 0; i-- {
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
	return buffer.String(), isReplaced
}

func ExpandMacro(buffer *bytes.Buffer, reader *strings.Reader, line string) {
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
	if ctx == nil {
		fmt.Fprintln(cmd.Stderr, "history not found (case1)\n")
		return 1, nil
	}
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

	historyObj_ := ctx.Value("history")
	if historyObj, ok := historyObj_.(*Container); ok {
		if f, ok := cmd.Stdout.(*os.File); (!ok || isatty.IsTerminal(f.Fd())) &&
			historyObj.Len() > num {

			start = historyObj.Len() - num
		}
		home := os.Getenv("USERPROFILE")
		for i := start; i < historyObj.Len(); i++ {
			row := historyObj.rows[i]
			dir := row.Dir
			if strings.HasPrefix(strings.ToUpper(dir), strings.ToUpper(home)) {
				dir = "~" + dir[len(home):]
			}
			dir = strings.Replace(dir, "\\", "/", -1)
			fmt.Fprintf(cmd.Stdout, "%s %-s (%s)\n",
				row.Stamp.Format("Jan _2 15:04:05"),
				row.Text,
				dir)
		}
	} else {
		fmt.Fprintln(cmd.Stderr, "history not found (case 2)")
	}
	return 0, nil
}

const max_histories = 2000

func (row *Line) String() string {
	return fmt.Sprintf("%s\t%s\t%s",
		row.Text,
		row.Dir,
		row.Stamp.Format("2006-01-02 15:04:05"))
}

func (hisObj *Container) WriteTo(w io.Writer) {
	i := 0
	if len(hisObj.rows) > max_histories {
		i = len(hisObj.rows) - max_histories
	}
	bw := bufio.NewWriter(w)
	for ; i < len(hisObj.rows); i++ {
		fmt.Fprintln(bw, hisObj.rows[i].String())
	}
	bw.Flush()
}

func (hisObj *Container) Save(path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	hisObj.WriteTo(fd)
	fd.Close()
	return nil
}

func (hisObj *Container) ReadFrom(reader io.Reader) {
	sc := bufio.NewScanner(reader)
	list := make([][]string, 0, 2000)
	hash := make(map[string]int)
	for sc.Scan() {
		line := sc.Text()
		if lnum, ok := hash[line]; ok {
			// delete duplicated record (marking)
			list[lnum] = nil
		}
		hash[line] = len(list)

		p := strings.Split(line, "\t")
		list = append(list, p)
	}
	for _, p := range list {
		// push only not duplicated record.
		if p != nil {
			var stamp time.Time
			var dir string
			if len(p) >= 3 {
				dir = p[1]
				stamp, _ = time.Parse("2006-01-02 15:04:05", p[2])
			}
			hisObj.PushLine(Line{
				Text:  p[0],
				Dir:   dir,
				Stamp: stamp})
		}
	}
}

func (hisObj *Container) Load(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	hisObj.ReadFrom(fd)
	fd.Close()
	return nil
}
