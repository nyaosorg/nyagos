package history

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"../conio"
	"../interpreter"
	"../lua"
)

func atoi_(reader *strings.Reader) (int, int) {
	n := 0
	count := 0
	for reader.Len() > 0 {
		ch, _, _ := reader.ReadRune()
		index := strings.IndexRune("0123456789", ch)
		if index >= 0 {
			n = n*10 + index
			count++
		} else {
			reader.UnreadRune()
			break
		}
	}
	return n, count
}

var Mark lua.Pushable = lua.TString{"!"}

func Replace(line string) (string, bool) {
	var mark = '!'
	if tstr1, ok := Mark.(lua.TString); ok {
		for _, c := range tstr1.Value {
			mark = c
			break
		}
	} else if tstr1, ok := Mark.(lua.TRawString); ok {
		mark = rune(tstr1.Value[0])
	} else {
		return line, false
	}

	var buffer bytes.Buffer
	isReplaced := false
	reader := strings.NewReader(line)
	history_count := conio.DefaultEditor.HistoryLen()

	for reader.Len() > 0 {
		ch, _, _ := reader.ReadRune()
		if ch != mark || reader.Len() <= 0 {
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
				if strings.Contains(conio.DefaultEditor.Histories[i].Line, seekStr) {
					buffer.WriteString(conio.DefaultEditor.Histories[i].Line)
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
			if strings.HasPrefix(conio.DefaultEditor.Histories[i].Line, seekStr) {
				buffer.WriteString(conio.DefaultEditor.Histories[i].Line)
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
	result := conio.NewHistoryLine(buffer.String())
	if isReplaced {
		if history_count > 0 {
			conio.DefaultEditor.Histories[history_count-1] = result
		} else {
			conio.DefaultEditor.Histories = append(conio.DefaultEditor.Histories, result)
		}
	}
	return result.Line, isReplaced
}

func insertHistory(buffer *bytes.Buffer, reader *strings.Reader, historyNo int) {
	history1 := conio.DefaultEditor.Histories[historyNo]
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

func CmdHistory(cmd *interpreter.Interpreter) (interpreter.ErrorLevel, error) {
	var num int
	if len(cmd.Args) >= 2 {
		num64, err := strconv.ParseInt(cmd.Args[1], 0, 32)
		if err != nil {
			switch err.(type) {
			case *strconv.NumError:
				return interpreter.NOERROR, fmt.Errorf(
					"history: %s not a number", cmd.Args[1])
			default:
				return interpreter.NOERROR, err
			}
		}
		num = int(num64)
		if num < 0 {
			num = -num
		}
	} else {
		num = 10
	}
	var start int
	if conio.DefaultEditor.HistoryLen() > num {
		start = conio.DefaultEditor.HistoryLen() - num
	} else {
		start = 0
	}
	for i, s := range conio.DefaultEditor.Histories[start:] {
		fmt.Fprintf(cmd.Stdout, "%3d : %-s\n", start+i, s.Line)
	}
	return interpreter.NOERROR, nil
}

const max_histories = 2000

func Save(path string) error {
	conio.DefaultEditor.ShrinkHistory()
	start := 0
	if conio.DefaultEditor.HistoryLen() > max_histories {
		start = conio.DefaultEditor.HistoryLen() - max_histories
	}
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	for i := start; i < len(conio.DefaultEditor.Histories); i++ {
		fmt.Fprintln(fd, conio.DefaultEditor.Histories[i].Line)
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
		conio.DefaultEditor.Histories = append(conio.DefaultEditor.Histories, conio.NewHistoryLine(sc.Text()))
	}
	conio.DefaultEditor.ShrinkHistory()
	return nil
}
