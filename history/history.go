package history

import "bufio"
import "bytes"
import "fmt"
import "os"
import "os/exec"
import "strings"

import "../conio"
import "../interpreter"

var histories = make([]string, 0)
var pointor = 0

func KeyFuncHistoryUp(this *conio.ReadLineBuffer) conio.KeyFuncResult {
	if pointor <= 0 {
		pointor = len(histories)
	}
	pointor -= 1
	conio.KeyFuncClear(this)
	if pointor >= 0 {
		this.InsertString(0, histories[pointor])
		this.ViewStart = 0
		this.Cursor = 0
		conio.KeyFuncTail(this)
	}
	return conio.CONTINUE
}

func KeyFuncHistoryDown(this *conio.ReadLineBuffer) conio.KeyFuncResult {
	pointor += 1
	if pointor >= len(histories) {
		pointor = 0
	}
	conio.KeyFuncClear(this)
	if pointor < len(histories) {
		this.InsertString(0, histories[pointor])
		this.ViewStart = 0
		this.Cursor = 0
		conio.KeyFuncTail(this)
	}
	return conio.CONTINUE
}

func Push(input string) {
	histories = append(histories, input)
	pointor = len(histories)
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
		if ch == '!' {
			if len(histories) > 0 {
				insertHisotry(&buffer, reader, histories[len(histories)-1])
				isReplaced = true
				continue
			} else {
				buffer.WriteRune('!')
				continue
			}
		}
		if n := strings.IndexRune("0123456789", ch); n >= 0 {
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
		if ch == '-' && reader.Len() > 0 {
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
			} else {
				buffer.WriteString("!-")
				buffer.WriteRune(ch)
			}
		} else {
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

func CmdHistory(cmd *exec.Cmd) interpreter.WhatToDoAfterCmd {
	for i, s := range histories {
		fmt.Fprintf(cmd.Stdout, "%3d : %-s\n", i, s)
	}
	return interpreter.CONTINUE
}

const max_histories = 256

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
		fd.WriteString(s)
		fd.WriteString("\n")
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
