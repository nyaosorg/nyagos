package history

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/nyaosorg/nyagos/texts"
)

var Mark = "!"

var DisableMarks = "\"'"

func (row *Line) dump() string {
	home := os.Getenv("USERPROFILE")
	dir := row.Dir
	if strings.HasPrefix(strings.ToUpper(dir), strings.ToUpper(home)) {
		dir = "~" + dir[len(home):]
	}
	dir = filepath.ToSlash(dir)
	return fmt.Sprintf("%s [%d] %-s (%s)",
		row.Stamp.Format("Jan _2 15:04:05"),
		row.Pid,
		row.Text,
		dir)
}

func (hisObj *Container) Replace(line string) (string, bool, error) {
	var mark rune
	for _, c := range Mark {
		mark = c
		break
	}

	var buffer strings.Builder
	isReplaced := false
	reader := strings.NewReader(line)
	historyCount := hisObj.Len()

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
		if unicode.IsSpace(ch) {
			buffer.WriteRune('!')
			buffer.WriteRune(ch)
			continue
		}
		if n := strings.IndexRune("^$:*@", ch); n >= 0 {
			reader.UnreadRune()
			if historyCount >= 1 {
				line := hisObj.GetAt(historyCount - 1)
				ExpandMacro(&buffer, reader, line)
				isReplaced = true
			}
			continue
		}
		if ch == mark { // !!
			if historyCount >= 1 {
				line := hisObj.GetAt(historyCount - 1)
				ExpandMacro(&buffer, reader, line)
				isReplaced = true
				continue
			} else {
				return "", false, errors.New("!!: event not found")
			}
		}
		if unicode.IsDigit(ch) { // !n
			reader.UnreadRune()
			var backno int
			fmt.Fscan(reader, &backno)
			if 0 <= backno && backno < historyCount {
				line := hisObj.GetAt(backno)
				ExpandMacro(&buffer, reader, line)
				isReplaced = true
			} else {
				return "", false, fmt.Errorf("!%d: event not found", backno)
			}
			continue
		}
		if ch == '-' && reader.Len() > 0 { // !-n
			var number int
			if _, err := fmt.Fscan(reader, &number); err == nil {
				backno := historyCount - number
				if 0 <= backno && backno < historyCount {
					line := hisObj.GetAt(backno)
					ExpandMacro(&buffer, reader, line)
					isReplaced = true
				} else {
					return "", false, fmt.Errorf("!-%d: event not found", number)
				}
				continue
			} else {
				reader.UnreadRune() // next char of '-'
			}
		}
		if ch == '?' { // !?str?
			var seekStrBuf strings.Builder
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
			for i := historyCount - 1; i >= 0; i-- {
				his1 := hisObj.GetAt(i)
				if strings.Contains(his1.Text, seekStr) {
					ExpandMacro(&buffer, reader, his1)
					isReplaced = true
					found = true
					break
				}
			}
			if !found {
				if lastCharIsQuestionMark {
					return "", false, fmt.Errorf("?%s?: event not found", seekStr)
				}
				return "", false, fmt.Errorf("?%s: event not found", seekStr)
			}
			continue
		}
		// !str
		var seekStrBuf strings.Builder
		seekStrBuf.WriteRune(ch)
		for reader.Len() > 0 {
			ch, _, _ := reader.ReadRune()
			if unicode.IsSpace(ch) || ch == ':' {
				reader.UnreadRune()
				break
			}
			seekStrBuf.WriteRune(ch)
		}
		seekStr := seekStrBuf.String()
		found := false
		for i := historyCount - 1; i >= 0; i-- {
			his1 := hisObj.GetAt(i)
			if strings.HasPrefix(his1.Text, seekStr) {
				ExpandMacro(&buffer, reader, his1)
				isReplaced = true
				found = true
				break
			}
		}
		if !found {
			return "", false, fmt.Errorf("%c%s: event not found", mark, seekStr)
		}
	}
	return buffer.String(), isReplaced, nil
}

func ExpandMacro(buffer *strings.Builder, reader *strings.Reader, row *Line) {
	line := row.Text
	dir := row.Dir

	ch, siz, _ := reader.ReadRune()
	if siz > 0 && ch == ':' {
		ch, siz, _ = reader.ReadRune()
	}
	if siz > 0 && ch == '^' {
		if words := texts.SplitLikeShellString(line); len(words) >= 2 {
			buffer.WriteString(words[1])
		}
	} else if siz > 0 && ch == '$' {
		if words := texts.SplitLikeShellString(line); len(words) >= 2 {
			buffer.WriteString(words[len(words)-1])
		}
	} else if siz > 0 && ch == '*' {
		if words := texts.SplitLikeShellString(line); len(words) >= 2 {
			buffer.WriteString(strings.Join(words[1:], " "))
		}
	} else if siz > 0 && ch == '@' {
		buffer.WriteString(dir)
	} else if siz > 0 && unicode.IsDigit(ch) {
		var n int
		reader.UnreadRune()
		fmt.Fscan(reader, &n)
		words := texts.SplitLikeShellString(line)
		if n < len(words) {
			buffer.WriteString(words[n])
		}
	} else {
		if siz > 0 {
			reader.UnreadRune()
		}
		buffer.WriteString(line)
	}
}

var MaxSaveHistory = 1000

func (hisObj *Container) SaveViaWriter(w io.Writer) {
	i := 0
	if len(hisObj.rows) > MaxSaveHistory {
		i = len(hisObj.rows) - MaxSaveHistory
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
	hisObj.SaveViaWriter(fd)
	return fd.Close()
}

func (hisObj *Container) LoadViaReader(reader io.Reader) {
	sc := bufio.NewScanner(reader)
	for sc.Scan() {
		line := sc.Text()
		p := strings.Split(line, "\t")

		var stamp time.Time
		var dir string
		pid := 0
		if len(p) >= 3 {
			dir = p[1]
			stamp, _ = time.Parse("2006-01-02 15:04:05", p[2])
			if len(p) >= 4 {
				pid, _ = strconv.Atoi(p[3])
			}
		}
		hisObj.PushLine(Line{
			Text:  p[0],
			Dir:   dir,
			Stamp: stamp,
			Pid:   pid})
	}
	sort.Slice(hisObj.rows, func(i, j int) bool {
		p := hisObj.rows[i]
		q := hisObj.rows[j]
		if p.Stamp != q.Stamp {
			return p.Stamp.Before(q.Stamp)
		}
		return p.Text < q.Text
	})
}

func (hisObj *Container) Load(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	hisObj.LoadViaReader(fd)
	return fd.Close()
}
