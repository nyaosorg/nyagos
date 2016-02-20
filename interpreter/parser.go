package interpreter

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"../dos"
)

type StatementT struct {
	Args     []string
	Redirect []*Redirecter
	Term     string
}

var prefix []string = []string{" 0<", " 1>", " 2>"}

var PercentFunc = map[string]func() string{
	"CD": func() string {
		wd, err := os.Getwd()
		if err == nil {
			return wd
		} else {
			return ""
		}
	},
	"ERRORLEVEL": func() string {
		return ErrorLevelStr
	},
}

var rxUnicode = regexp.MustCompile("^[uU]\\+?([0-9a-fA-F]+)$")

func OurGetEnv(name string) (string, bool) {
	value := os.Getenv(name)
	if value != "" {
		return value, true
	} else if m := rxUnicode.FindStringSubmatch(name); m != nil {
		ucode, _ := strconv.ParseInt(m[1], 16, 32)
		return fmt.Sprintf("%c", rune(ucode)), true
	} else if f, ok := PercentFunc[strings.ToUpper(name)]; ok {
		return f(), true
	} else {
		return "", false
	}
}

func chomp(buffer *bytes.Buffer) {
	original := buffer.String()
	buffer.Reset()
	var lastchar rune
	for i, ch := range original {
		if i > 0 {
			buffer.WriteRune(lastchar)
		}
		lastchar = ch
	}
}

const NOTQUOTED = '\000'

const EMPTY_COMMAND_FOUND = "Empty command found"

func dequote(source *bytes.Buffer) string {
	var buffer bytes.Buffer

	lastchar := ' '
	quoteNow := NOTQUOTED
	yenCount := 0
	for {
		ch, _, err := source.ReadRune()
		if err != nil {
			break
		}
		if ch == '~' && unicode.IsSpace(lastchar) {
			if home := dos.GetHome(); home != "" {
				buffer.WriteString(home)
			} else {
				buffer.WriteRune('~')
			}
			lastchar = '~'
			continue
		}
		if ch == '%' && quoteNow != '\'' && yenCount%2 == 0 {
			var nameBuf bytes.Buffer
			for {
				ch, _, err = source.ReadRune()
				if err != nil {
					buffer.WriteRune('%')
					nameBuf.WriteTo(&buffer)
					return buffer.String()
				}
				if ch == '%' {
					if value, ok := OurGetEnv(nameBuf.String()); ok {
						buffer.WriteString(value)
					} else {
						buffer.WriteRune('%')
						nameBuf.WriteTo(&buffer)
						buffer.WriteRune('%')
					}
					break
				}
				if ch == '=' {
					source.UnreadRune()
					buffer.WriteRune('%')
					nameBuf.WriteTo(&buffer)
					break
				}
				nameBuf.WriteRune(ch)
			}
			continue
		}

		if quoteNow != NOTQUOTED && ch == quoteNow && yenCount%2 == 0 {
			// Close Quotation.
			for ; yenCount >= 2; yenCount -= 2 {
				buffer.WriteRune('\\')
			}
			quoteNow = NOTQUOTED
		} else if (ch == '\'' || ch == '"') && quoteNow == NOTQUOTED && yenCount%2 == 0 {
			// Open Qutation.
			for ; yenCount >= 2; yenCount -= 2 {
				buffer.WriteRune('\\')
			}
			quoteNow = ch
			if ch == lastchar {
				buffer.WriteRune(ch)
			}
		} else {
			if ch == '\\' {
				yenCount++
			} else if ch == '\'' || ch == '"' {
				for ; yenCount >= 2; yenCount -= 2 {
					buffer.WriteRune('\\')
				}
				yenCount = 0
				buffer.WriteRune(ch)
			} else {
				for ; yenCount > 0; yenCount-- {
					buffer.WriteRune('\\')
				}
				buffer.WriteRune(ch)
			}
		}
		lastchar = ch
	}
	for ; yenCount > 0; yenCount-- {
		buffer.WriteRune('\\')
	}
	return buffer.String()
}

func parse1(text string) ([]*StatementT, error) {
	quoteNow := NOTQUOTED
	yenCount := 0
	statements := make([]*StatementT, 0)
	args := make([]string, 0)
	lastchar := ' '
	var buffer bytes.Buffer
	isNextRedirect := false
	redirect := make([]*Redirecter, 0, 3)

	terminate := func(term string) {
		statement1 := new(StatementT)
		if buffer.Len() > 0 {
			if isNextRedirect && len(redirect) > 0 {
				redirect[len(redirect)-1].SetPath(dequote(&buffer))
				isNextRedirect = false
				statement1.Args = args
			} else {
				statement1.Args = append(args, dequote(&buffer))
			}
			buffer.Reset()
		} else if len(args) <= 0 {
			return
		} else {
			statement1.Args = args
		}
		statement1.Redirect = redirect
		redirect = make([]*Redirecter, 0, 3)
		args = make([]string, 0)
		statement1.Term = term
		statements = append(statements, statement1)
	}

	TermWord := func() {
		if isNextRedirect && len(redirect) > 0 {
			redirect[len(redirect)-1].SetPath(dequote(&buffer))
		} else {
			if buffer.Len() > 0 {
				args = append(args, dequote(&buffer))
			}
		}
		buffer.Reset()
	}

	reader := strings.NewReader(text)
	for reader.Len() > 0 {
		ch, chSize, chErr := reader.ReadRune()
		if chSize <= 0 {
			break
		}
		if chErr != nil {
			return nil, chErr
		}
		if quoteNow == NOTQUOTED {
			if yenCount%2 == 0 && (ch == '"' || ch == '\'') {
				quoteNow = ch
			}
		} else if yenCount%2 == 0 && ch == quoteNow {
			quoteNow = NOTQUOTED
		}
		if quoteNow != NOTQUOTED {
			buffer.WriteRune(ch)
		} else if ch == ' ' {
			if buffer.Len() > 0 {
				TermWord()
				isNextRedirect = false
			}
		} else if lastchar == ' ' && ch == ';' {
			terminate(";")
		} else if ch == '|' {
			if lastchar == '|' {
				if len(statements) <= 0 {
					return nil, errors.New(EMPTY_COMMAND_FOUND)
				}
				statements[len(statements)-1].Term = "||"
			} else {
				terminate("|")
			}
		} else if ch == '&' {
			switch lastchar {
			case '&':
				if len(statements) <= 0 {
					return nil, errors.New(EMPTY_COMMAND_FOUND)
				}
				statements[len(statements)-1].Term = "&&"
			case '|':
				if len(statements) <= 0 {
					return nil, errors.New(EMPTY_COMMAND_FOUND)
				}
				statements[len(statements)-1].Term = "|&"
			case '>':
				// >&[n]
				ch2, ch2siz, ch2err := reader.ReadRune()
				if ch2err != nil {
					return nil, ch2err
				}
				if ch2siz <= 0 {
					return nil, errors.New("Too Near EOF for >&")
				}
				red := redirect[len(redirect)-1]
				switch ch2 {
				case '1':
					red.DupFrom(1)
				case '2':
					red.DupFrom(2)
				default:
					return nil, errors.New("Syntax error after >&")
				}
				isNextRedirect = false
			default:
				terminate("&")
			}
		} else if ch == '>' {
			switch lastchar {
			case '1':
				// 1>
				chomp(&buffer)
				TermWord()
				redirect = append(redirect, NewRedirecter(1))
			case '2':
				// 2>
				chomp(&buffer)
				TermWord()
				redirect = append(redirect, NewRedirecter(2))
			case '>':
				// >>
				TermWord()
				if len(redirect) >= 0 {
					redirect[len(redirect)-1].SetAppend()
				}
			default:
				// >
				TermWord()
				redirect = append(redirect, NewRedirecter(1))
			}
			isNextRedirect = true
		} else if ch == '<' {
			TermWord()
			redirect = append(redirect, NewRedirecter(0))
			isNextRedirect = true
		} else {
			buffer.WriteRune(ch)
		}
		if ch == '\\' {
			yenCount++
		} else {
			yenCount = 0
		}
		lastchar = ch
	}
	terminate(" ")
	return statements, nil
}

// Make arrays whose elements are pipelines
func parse2(statements []*StatementT) [][]*StatementT {
	result := make([][]*StatementT, 1)
	for _, statement1 := range statements {
		result[len(result)-1] = append(result[len(result)-1], statement1)
		switch statement1.Term {
		case "|", "|&":

		default:
			result = append(result, make([]*StatementT, 0))
		}
	}
	if len(result[len(result)-1]) <= 0 {
		result = result[0 : len(result)-1]
	}
	return result
}

func Parse(text string) ([][]*StatementT, error) {
	result1, err := parse1(text)
	if err != nil {
		return nil, err
	}
	result2 := parse2(result1)
	return result2, nil
}
