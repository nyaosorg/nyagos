package interpreter

import (
	"bytes"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"../dos"
)

type StatementT struct {
	Argv     []string
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
		return ErrorLevel
	},
}

var rxUnicode = regexp.MustCompile("^[uU]\\+?([0-9a-fA-F]+)$")

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

func dequote(source *bytes.Buffer) string {
	var buffer bytes.Buffer

	lastchar := ' '
	quote := false
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
		if ch == '%' {
			var nameBuf bytes.Buffer
			for {
				ch, _, err = source.ReadRune()
				if err != nil {
					buffer.WriteRune('%')
					buffer.WriteString(nameBuf.String())
					return buffer.String()
				}
				if ch == '%' {
					break
				}
				nameBuf.WriteRune(ch)
			}
			nameStr := nameBuf.String()
			value := os.Getenv(nameStr)
			if value != "" {
				buffer.WriteString(value)
			} else if m := rxUnicode.FindStringSubmatch(nameStr); m != nil {
				ucode, _ := strconv.ParseInt(m[1], 16, 32)
				buffer.WriteRune(rune(ucode))
			} else if f, ok := PercentFunc[nameStr]; ok {
				buffer.WriteString(f())
			} else {
				buffer.WriteRune('%')
				buffer.WriteString(nameStr)
				buffer.WriteRune('%')
			}
			continue
		}
		if ch == '"' {
			quote = !quote
			if lastchar == '"' && quote {
				buffer.WriteRune('"')
				lastchar = '\000'
			}
		} else {
			buffer.WriteRune(ch)
		}
		lastchar = ch
	}
	return buffer.String()
}

func terminate(statements *[]StatementT,
	isRedirected *bool,
	redirect *[]*Redirecter,
	buffer *bytes.Buffer,
	argv *[]string,
	term string) {

	var statement1 StatementT
	if buffer.Len() > 0 {
		if *isRedirected && len(*redirect) > 0 {
			(*redirect)[len(*redirect)-1].SetPath(dequote(buffer))
			*isRedirected = false
			statement1.Argv = *argv
		} else {
			statement1.Argv = append(*argv, dequote(buffer))
		}
		buffer.Reset()
	} else if len(*argv) <= 0 {
		return
	} else {
		statement1.Argv = *argv
	}
	statement1.Redirect = *redirect
	*redirect = make([]*Redirecter, 0, 3)
	*argv = make([]string, 0)
	statement1.Term = term
	*statements = append(*statements, statement1)
}

func parse1(text string) ([]StatementT, error) {
	isQuoted := false
	statements := make([]StatementT, 0)
	argv := make([]string, 0)
	lastchar := ' '
	var buffer bytes.Buffer
	isNextRedirect := false
	redirect := make([]*Redirecter, 0, 3)

	reader := strings.NewReader(text)
	for reader.Len() > 0 {
		ch, chSize, chErr := reader.ReadRune()
		if chSize <= 0 {
			break
		}
		if chErr != nil {
			return nil, chErr
		}
		if ch == '"' {
			isQuoted = !isQuoted
		}
		if isQuoted {
			buffer.WriteRune(ch)
		} else if ch == ' ' {
			if buffer.Len() > 0 {
				if isNextRedirect && len(redirect) > 0 {
					redirect[len(redirect)-1].SetPath(dequote(&buffer))
				} else {
					argv = append(argv, dequote(&buffer))
				}
				buffer.Reset()
				isNextRedirect = false
			}
		} else if lastchar == ' ' && ch == ';' {
			terminate(&statements, &isNextRedirect, &redirect, &buffer, &argv, ";")
		} else if ch == '|' {
			if lastchar == '|' {
				statements[len(statements)-1].Term = "||"
			} else {
				terminate(&statements, &isNextRedirect, &redirect, &buffer, &argv, "|")
			}
		} else if ch == '&' {
			switch lastchar {
			case '&':
				statements[len(statements)-1].Term = "&&"
			case '|':
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
				terminate(&statements, &isNextRedirect, &redirect, &buffer, &argv, "&")
			}
		} else if ch == '>' {
			switch lastchar {
			case '1':
				// 1>
				chomp(&buffer)
				redirect = append(redirect, NewRedirecter(1))
			case '2':
				// 2>
				chomp(&buffer)
				redirect = append(redirect, NewRedirecter(2))
			case '>':
				// >>
				if len(redirect) >= 0 {
					redirect[len(redirect)-1].SetAppend()
				}
			default:
				// >
				redirect = append(redirect, NewRedirecter(1))
			}
			isNextRedirect = true
		} else if ch == '<' {
			redirect = append(redirect, NewRedirecter(0))
			isNextRedirect = true
		} else {
			buffer.WriteRune(ch)
		}
		lastchar = ch
	}
	terminate(&statements, &isNextRedirect, &redirect, &buffer, &argv, " ")
	return statements, nil
}

// Make arrays whose elements are pipelines
func parse2(statements []StatementT) [][]StatementT {
	result := make([][]StatementT, 1)
	for _, statement1 := range statements {
		result[len(result)-1] = append(result[len(result)-1], statement1)
		switch statement1.Term {
		case "|", "|&":

		default:
			result = append(result, make([]StatementT, 0))
		}
	}
	if len(result[len(result)-1]) <= 0 {
		result = result[0 : len(result)-1]
	}
	return result
}

func Parse(text string) ([][]StatementT, error) {
	result1, err := parse1(text)
	if err != nil {
		return nil, err
	}
	result2 := parse2(result1)
	return result2, nil
}
