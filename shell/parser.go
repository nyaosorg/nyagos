package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/texts"
)

type StatementT struct {
	Args     []string
	RawArgs  []string
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
		return fmt.Sprintf("%d", LastErrorLevel)
	},
}

var rxUnicode = regexp.MustCompile("^[uU]\\+?([0-9a-fA-F]+)$")

var rxSubstitute = regexp.MustCompile(`^([^\:]+)\:([^\=]+)=(.*)$`)

func ourGetenvSub(name string) (string, bool) {
	m := rxSubstitute.FindStringSubmatch(name)
	if m != nil {
		base, ok := OurGetEnv(m[1])
		if ok {
			return texts.ReplaceIgnoreCase(base, m[2], m[3]), true
		} else {
			return "", false
		}
	} else {
		return OurGetEnv(name)
	}
}

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

func chomp(buffer *strings.Builder) {
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

func string2word(source_ string, removeQuote bool) string {
	var buffer strings.Builder
	source := strings.NewReader(source_)

	lastchar := ' '
	quoteNow := NOTQUOTED
	yenCount := 0
	for {
		ch, _, err := source.ReadRune()
		if err != nil {
			break
		}
		if ch == '~' && unicode.IsSpace(lastchar) && quoteNow == NOTQUOTED {
			if home := dos.GetHome(); home != "" {
				buffer.WriteString(home)
			} else {
				buffer.WriteRune('~')
			}
			lastchar = '~'
			continue
		}
		if ch == '%' && quoteNow != '\'' {
			for ; yenCount > 0; yenCount-- {
				buffer.WriteRune('\\')
			}
			var nameBuf strings.Builder
			for {
				ch, _, err = source.ReadRune()
				if err != nil {
					buffer.WriteRune('%')
					source.Seek(-int64(nameBuf.Len()), io.SeekCurrent)
					break
				}
				if ch == '%' {
					if value, ok := ourGetenvSub(nameBuf.String()); ok {
						buffer.WriteString(value)
					} else {
						buffer.WriteRune('%')
						source.Seek(-int64(nameBuf.Len()+1), io.SeekCurrent)
					}
					break
				}
				nameBuf.WriteRune(ch)
			}
			continue
		}

		if quoteNow != NOTQUOTED && ch == quoteNow && yenCount%2 == 0 {
			if !removeQuote {
				buffer.WriteRune(ch)
			}
			// Close Quotation.
			for ; yenCount >= 2; yenCount -= 2 {
				buffer.WriteRune('\\')
			}
			quoteNow = NOTQUOTED
		} else if (ch == '\'' || ch == '"') && quoteNow == NOTQUOTED && yenCount%2 == 0 {
			if !removeQuote {
				buffer.WriteRune(ch)
			}
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
	rawArgs := make([]string, 0)
	lastchar := ' '
	var buffer strings.Builder
	isNextRedirect := false
	redirect := make([]*Redirecter, 0, 3)

	term_line := func(term string) {
		statement1 := new(StatementT)
		if buffer.Len() > 0 {
			if isNextRedirect && len(redirect) > 0 {
				redirect[len(redirect)-1].SetPath(string2word(buffer.String(), true))
				isNextRedirect = false
				statement1.RawArgs = rawArgs
				statement1.Args = args
			} else {
				statement1.RawArgs = append(rawArgs, string2word(buffer.String(), false))
				statement1.Args = append(args, string2word(buffer.String(), true))
			}
			buffer.Reset()
		} else if len(args) <= 0 {
			return
		} else {
			statement1.RawArgs = rawArgs
			statement1.Args = args
		}
		statement1.Redirect = redirect
		redirect = make([]*Redirecter, 0, 3)
		rawArgs = make([]string, 0)
		args = make([]string, 0)
		statement1.Term = term
		statements = append(statements, statement1)
	}

	term_word := func() {
		if isNextRedirect && len(redirect) > 0 {
			redirect[len(redirect)-1].SetPath(string2word(buffer.String(), true))
		} else {
			if buffer.Len() > 0 {
				rawArgs = append(rawArgs, string2word(buffer.String(), false))
				args = append(args, string2word(buffer.String(), true))
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
		} else if unicode.IsSpace(ch) {
			if buffer.Len() > 0 {
				term_word()
				isNextRedirect = false
			}
		} else if unicode.IsSpace(lastchar) && ch == '#' {
			break
		} else if unicode.IsSpace(lastchar) && ch == ';' {
			term_line(";")
		} else if ch == '!' && lastchar == '>' && isNextRedirect && len(redirect) > 0 {
			redirect[len(redirect)-1].force = true
		} else if ch == '|' {
			if lastchar == '>' && isNextRedirect && len(redirect) > 0 {
				redirect[len(redirect)-1].force = true
			} else if lastchar == '|' {
				if len(statements) <= 0 {
					return nil, errors.New(EMPTY_COMMAND_FOUND)
				}
				statements[len(statements)-1].Term = "||"
			} else {
				term_line("|")
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
				term_line("&")
			}
		} else if ch == '>' {
			switch lastchar {
			case '1':
				// 1>
				chomp(&buffer)
				term_word()
				redirect = append(redirect, NewRedirecter(1))
			case '2':
				// 2>
				chomp(&buffer)
				term_word()
				redirect = append(redirect, NewRedirecter(2))
			case '>':
				// >>
				term_word()
				if len(redirect) >= 0 {
					redirect[len(redirect)-1].SetAppend()
				}
			default:
				// >
				term_word()
				redirect = append(redirect, NewRedirecter(1))
			}
			isNextRedirect = true
		} else if ch == '<' {
			term_word()
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
	term_line(" ")
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
