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
)

func Replace(line string) (string, bool) {
	var buffer bytes.Buffer
	isReplaced := false
	reader := strings.NewReader(line)
	history_count := len(conio.Histories)

	for reader.Len() > 0 {
		ch, _, _ := reader.ReadRune()
		if ch != '!' || reader.Len() <= 0 {
			buffer.WriteRune(ch)
			continue
		}
		ch, _, _ = reader.ReadRune()
		if n := strings.IndexRune("^$:*", ch); n >= 0 {
			reader.UnreadRune()
			if history_count >= 2 {
				insertHistory(&buffer, reader, conio.Histories[history_count-2])
				isReplaced = true
			}
			continue
		}
		if ch == '!' { // !!
			if history_count >= 2 {
				insertHistory(&buffer, reader, conio.Histories[history_count-2])
				isReplaced = true
				continue
			} else {
				buffer.WriteRune('!')
				continue
			}
		}
		if n := strings.IndexRune("0123456789", ch); n >= 0 { // !n
			backno := n
			for reader.Len() > 0 {
				ch, _, _ = reader.ReadRune()
				if n = strings.IndexRune("0123456789", ch); n >= 0 {
					backno = backno*10 + n
				} else {
					reader.UnreadRune()
					break
				}
			}
			backno = backno % history_count
			if 0 <= backno && backno < history_count {
				insertHistory(&buffer, reader, conio.Histories[backno])
				isReplaced = true
			}
			continue
		}
		if ch == '-' && reader.Len() > 0 { // !-n
			ch, _, _ := reader.ReadRune()
			n := strings.IndexRune("0123456789", ch)
			if n >= 0 {
				number := n
				for reader.Len() > 0 {
					ch, _, _ = reader.ReadRune()
					n = strings.IndexRune("0123456789", ch)
					if n < 0 {
						reader.UnreadRune()
						break
					}
					number = number*10 + n
				}
				backno := history_count - number - 1
				for backno < 0 {
					backno += history_count
				}
				if 0 <= backno && backno < history_count {
					insertHistory(&buffer, reader, conio.Histories[backno])
					isReplaced = true
				} else {
					buffer.WriteString("!-0")
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
				if strings.Contains(conio.Histories[i], seekStr) {
					buffer.WriteString(conio.Histories[i])
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
			if strings.HasPrefix(conio.Histories[i], seekStr) {
				buffer.WriteString(conio.Histories[i])
				isReplaced = true
				found = true
				break
			}
		}
		if !found {
			buffer.WriteRune('!')
			buffer.WriteRune(ch)
		}
	}
	result := buffer.String()
	if isReplaced {
		if history_count > 0 {
			conio.Histories[history_count-1] = result
		} else {
			conio.Histories = append(conio.Histories, result)
		}
	}
	return result, isReplaced
}

func splitQ(s string) []string {
	args := []string{}
	reader := strings.NewReader(s)
	for reader.Len() > 0 {
		var buffer bytes.Buffer
		for {
			if reader.Len() <= 0 {
				return args
			}
			ch, _, _ := reader.ReadRune()
			if ch != ' ' {
				reader.UnreadRune()
				break
			}
		}
		quote := false
		for reader.Len() > 0 {
			ch, _, _ := reader.ReadRune()
			if ch == '"' {
				quote = !quote
			}
			if ch == ' ' && !quote {
				break
			}
			buffer.WriteRune(ch)
		}
		if buffer.Len() > 0 {
			args = append(args, buffer.String())
		}
	}
	return args
}

func insertHistory(buffer *bytes.Buffer, reader *strings.Reader, history1 string) {
	ch, siz, _ := reader.ReadRune()
	if siz > 0 && ch == '^' {
		args := splitQ(history1)
		if len(args) >= 2 {
			buffer.WriteString(args[1])
		}
	} else if siz > 0 && ch == '$' {
		args := splitQ(history1)
		if len(args) >= 2 {
			buffer.WriteString(args[len(args)-1])
		}
	} else if siz > 0 && ch == '*' {
		args := splitQ(history1)
		if len(args) >= 2 {
			buffer.WriteString(strings.Join(args[1:], " "))
		}
	} else if siz > 0 && ch == ':' {
		args := splitQ(history1)
		n := 0
		count := 0
		for reader.Len() > 0 {
			ch, _, _ = reader.ReadRune()
			index := strings.IndexRune("0123456789", ch)
			if index >= 0 {
				n = n*10 + index
				count++
			} else {
				reader.UnreadRune()
				break
			}
		}
		if count <= 0 {
			buffer.WriteRune(':')
		} else if n < len(args) {
			buffer.WriteString(args[n])
		}
	} else {
		if siz > 0 {
			reader.UnreadRune()
		}
		buffer.WriteString(history1)
	}
}

func CmdHistory(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	var num int
	if len(cmd.Args) >= 2 {
		num64, err := strconv.ParseInt(cmd.Args[1], 0, 32)
		if err != nil {
			return interpreter.CONTINUE, err
		}
		num = int(num64)
	} else {
		num = 10
	}
	var start int
	if len(conio.Histories) > num {
		start = len(conio.Histories) - num
	} else {
		start = 0
	}
	for i, s := range conio.Histories[start:] {
		fmt.Fprintf(cmd.Stdout, "%3d : %-s\n", start+i, s)
	}
	return interpreter.CONTINUE, nil
}

const max_histories = 2000

func Save(path string) error {
	var hist_ []string
	if len(conio.Histories) > max_histories {
		hist_ = conio.Histories[(len(conio.Histories) - max_histories):]
	} else {
		hist_ = conio.Histories
	}
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	for _, s := range hist_ {
		fmt.Fprintln(fd, s)
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
		conio.Histories = append(conio.Histories, sc.Text())
	}
	return nil
}
