package history

import "bufio"
import "bytes"
import "fmt"
import "os"
import "strconv"
import "strings"

import "../conio/readline"
import "../interpreter"

var histories = make([]string, 0)
var pointer = 0

func Get(n int) string {
	if n < 0 {
		n = len(histories) + n
	}
	if n >= len(histories) {
		return ""
	} else {
		return histories[n]
	}
}

func Len() int {
	return len(histories)
}

func LastHistory() string {
	if len(histories) <= 0 {
		return ""
	} else {
		return histories[len(histories)-1]
	}
}

func KeyFuncHistoryUp(this *readline.Buffer) readline.Result {
	if pointer <= 0 {
		pointer = len(histories)
	}
	pointer -= 1
	readline.KeyFuncClear(this)
	if pointer >= 0 {
		this.InsertString(0, histories[pointer])
		this.ViewStart = 0
		this.Cursor = 0
		readline.KeyFuncTail(this)
	}
	return readline.CONTINUE
}

func KeyFuncHistoryDown(this *readline.Buffer) readline.Result {
	pointer += 1
	if pointer >= len(histories) {
		pointer = 0
	}
	readline.KeyFuncClear(this)
	if pointer < len(histories) {
		this.InsertString(0, histories[pointer])
		this.ViewStart = 0
		this.Cursor = 0
		readline.KeyFuncTail(this)
	}
	return readline.CONTINUE
}

func Push(input string) {
	histories = append(histories, input)
	ResetPointer()
}

func Replace(line string) (string, bool) {
	var buffer bytes.Buffer
	var isReplaced = false
	reader := strings.NewReader(line)

	for reader.Len() > 0 {
		ch, _, _ := reader.ReadRune()
		if ch != '!' || reader.Len() <= 0 {
			buffer.WriteRune(ch)
			continue
		}
		ch, _, _ = reader.ReadRune()
		if n := strings.IndexRune("^$:*", ch); n >= 0 {
			reader.UnreadRune()
			if len(histories) > 0 {
				insertHisotry(&buffer, reader, histories[len(histories)-1])
				isReplaced = true
			}
			continue
		}
		if ch == '!' { // !!
			if len(histories) > 0 {
				insertHisotry(&buffer, reader, histories[len(histories)-1])
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
			backno = backno % len(histories)
			if 0 <= backno && backno < len(histories) {
				insertHisotry(&buffer, reader, histories[backno])
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
				backno := len(histories) - number
				for backno < 0 {
					backno += len(histories)
				}
				if 0 <= backno && backno < len(histories) {
					insertHisotry(&buffer, reader, histories[backno])
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
			for i := len(histories) - 1; i >= 0; i-- {
				if strings.Contains(histories[i], seekStr) {
					buffer.WriteString(histories[i])
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
		for i := len(histories) - 1; i >= 0; i-- {
			if strings.HasPrefix(histories[i], seekStr) {
				buffer.WriteString(histories[i])
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
	return buffer.String(), isReplaced
}

func splitQ(s string) []string {
	args := make([]string, 0)
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
		s := buffer.String()
		if s != "" {
			args = append(args, s)
		}
	}
	return args
}

func insertHisotry(buffer *bytes.Buffer, reader *strings.Reader, history1 string) {
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
	if len(histories) > num {
		start = len(histories) - num
	} else {
		start = 0
	}
	for i, s := range histories[start:] {
		fmt.Fprintf(cmd.Stdout, "%3d : %-s\n", start+i, s)
	}
	return interpreter.CONTINUE, nil
}

const max_histories = 2000

func Save(path string) error {
	var hist_ []string
	if len(histories) > max_histories {
		hist_ = histories[(len(histories) - max_histories):]
	} else {
		hist_ = histories
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
		histories = append(histories, sc.Text())
	}
	return nil
}

func ResetPointer() {
	pointer = len(histories)
}
