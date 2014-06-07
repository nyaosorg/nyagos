package history

import "bytes"
import "strings"

import "../conio"

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
		if ch == '!' {
			if len(histories) > 0 {
				insertHisotry(&buffer, reader, histories[len(histories)-1])
				isReplaced = true
				continue
			} else {
				buffer.WriteRune('!')
				break
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

func insertHisotry(buffer *bytes.Buffer, reader *strings.Reader, base string) {
	buffer.WriteString(base)
}
