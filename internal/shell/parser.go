package shell

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

var NoClobber = false

func isSpace(c rune) bool {
	return strings.ContainsRune(" \t\n\r\v\f", c)
}

type _Statement struct {
	Args     []string
	RawArgs  []string
	Term     string
	Redirect []func([]*os.File) (func(), error)
}

var PercentFunc = map[string]func() string{
	"CD": func() string {
		wd, err := os.Getwd()
		if err == nil {
			return wd
		}
		return ""
	},
	"ERRORLEVEL": func() string {
		return fmt.Sprintf("%d", LastErrorLevel)
	},
	"DATE": func() string {
		s, err := nodos.TimeFormatOsLayout(time.Now())
		if err != nil {
			return err.Error()
		}
		return s
	},
	"TIME": func() string {
		return time.Now().Format("15:04:05.00")
	},
}

var rxUnicode = regexp.MustCompile(`^[uU]\+?([0-9a-fA-F]+)$`)

var rxSubstitute = regexp.MustCompile(`^([^\:]+)\:([^\=]+)=(.*)$`)

var rxSubstring = regexp.MustCompile(`^([^\:]+)\:\~(\-?\d+)(?:,(\-?\d+))?$`)

var rxPercent = regexp.MustCompile(`%\w+%`)

func minusPosToAbs(s string, pos int) int {
	if pos < 0 {
		pos = len(s) + pos
	}
	if pos < 0 {
		pos = 0
	}
	if pos > len(s) {
		pos = len(s)
	}
	return pos
}

func ourGetenvSub(name string) (string, bool) {
	m := rxSubstitute.FindStringSubmatch(name)
	if m != nil {
		base, ok := OurGetEnv(m[1])
		if !ok {
			return "", false
		}
		return ReplaceIgnoreCase(base, m[2], m[3]), true
	}
	m = rxSubstring.FindStringSubmatch(name)
	if m != nil {
		base, ok := OurGetEnv(m[1])
		if !ok {
			return "", false
		}
		pos, _ := strconv.Atoi(m[2])
		pos = minusPosToAbs(base, pos)
		base = base[pos:]
		if len(m) >= 4 {
			if pos, err := strconv.Atoi(m[3]); err == nil {
				pos = minusPosToAbs(base, pos)
				base = base[:pos]
			}
		}
		return base, true
	}
	return OurGetEnv(name)
}

func rune2string(r rune) string {
	var b [utf8.UTFMax]byte
	n := utf8.EncodeRune(b[:], r)
	return string(b[:n])
}

func OurGetEnv(name string) (string, bool) {
	value := os.Getenv(name)
	if value != "" {
		return value, true
	} else if m := rxUnicode.FindStringSubmatch(name); m != nil {
		ucode, _ := strconv.ParseInt(m[1], 16, 32)
		return rune2string(rune(ucode)), true
	} else if f, ok := PercentFunc[strings.ToUpper(name)]; ok {
		return f(), true
	} else {
		return "", false
	}
}

const _NotQuoted = '\000'

var TildeExpansion = true

func writeRootRep(yenCount *int, cooked bool, buffer *strings.Builder) {
	if cooked {
		for ; *yenCount >= 2; *yenCount -= 2 {
			buffer.WriteByte('\\')
		}
	} else {
		for ; *yenCount >= 1; *yenCount-- {
			buffer.WriteByte('\\')
		}
	}
}

func string2word(_source string, cooked bool) string {
	var buffer strings.Builder
	source := strings.NewReader(_source)

	lastchar := ' '
	quoteNow := _NotQuoted
	yenCount := 0
	for {
		ch, _, err := source.ReadRune()
		if err != nil {
			break
		}
		if TildeExpansion && ch == '~' && isSpace(lastchar) && quoteNow == _NotQuoted {
			var name strings.Builder
			var undo strings.Builder
			undo.WriteByte('~')
			for {
				ch, _, err = source.ReadRune()
				if err != nil {
					//undo.WriteRune(ch)
					break
				} else if ch == '"' {
					undo.WriteByte('"')
					if quoteNow == _NotQuoted {
						quoteNow = ch
					} else {
						quoteNow = _NotQuoted
					}
				} else if !unicode.IsLetter(ch) {
					source.UnreadRune()
					break
				} else {
					undo.WriteRune(ch)
				}
				name.WriteRune(ch)
			}
			nameStr := strings.Replace(name.String(), `"`, ``, -1)
			if len(nameStr) > 0 {
				u, err := user.Lookup(nameStr)
				if err == nil {
					if !cooked && strings.Count(undo.String(), `"`)%2 != 0 {
						buffer.WriteByte('"')
					}
					buffer.WriteString(u.HomeDir)
					lastchar = rune(u.HomeDir[len(u.HomeDir)-1])
				} else {
					if !cooked {
						buffer.WriteByte('~')
					}
					undoStr := undo.String()
					buffer.WriteString(undoStr)
					lastchar = rune(undoStr[len(undoStr)-1])
				}
				continue
			}
			if home := nodos.GetHome(); home != "" {
				if !cooked && strings.Count(undo.String(), `"`)%2 != 0 {
					buffer.WriteByte('"')
				}
				buffer.WriteString(home)
			} else {
				buffer.WriteByte('~')
			}
			lastchar = '~'
			continue
		}
		if ch == '%' && quoteNow != '\'' {
			for ; yenCount > 0; yenCount-- {
				buffer.WriteByte('\\')
			}
			var nameBuf strings.Builder
			for {
				ch, _, err = source.ReadRune()
				if err != nil {
					buffer.WriteByte('%')
					source.Seek(-int64(nameBuf.Len()), io.SeekCurrent)
					break
				}
				if ch == '%' {
					if value, ok := ourGetenvSub(nameBuf.String()); ok {
						buffer.WriteString(value)
					} else {
						buffer.WriteByte('%')
						source.Seek(-int64(nameBuf.Len()+1), io.SeekCurrent)
					}
					break
				}
				nameBuf.WriteRune(ch)
			}
			continue
		}

		if quoteNow != _NotQuoted && ch == quoteNow && yenCount%2 == 0 {
			if !cooked {
				buffer.WriteRune(ch)
			}
			// Close Quotation.
			writeRootRep(&yenCount, cooked, &buffer)

			quoteNow = _NotQuoted
		} else if (ch == '\'' || ch == '"') && quoteNow == _NotQuoted && yenCount%2 == 0 {
			if !cooked {
				buffer.WriteRune(ch)
			}
			// Open Qutation.
			writeRootRep(&yenCount, cooked, &buffer)
			quoteNow = ch
			if ch == lastchar {
				buffer.WriteRune(ch)
			}
		} else {
			if ch == '\\' {
				yenCount++
			} else if ch == '\'' || ch == '"' {
				writeRootRep(&yenCount, cooked, &buffer)
				yenCount = 0
				buffer.WriteRune(ch)
			} else {
				for ; yenCount > 0; yenCount-- {
					buffer.WriteByte('\\')
				}
				buffer.WriteRune(ch)
			}
		}
		lastchar = ch
	}
	for ; yenCount > 0; yenCount-- {
		buffer.WriteByte('\\')
	}
	return reverse.Replace(buffer.String())
}

const (
	_AndAlso   = '\uE000' + iota // &&
	_OrElse                      // ||
	_Redirect0                   // 0<
	_Redirect1                   // 1>
	_Redirect2                   // 2>
	_Append1                     // 1>>
	_Append2                     // 2>>
	_Append                      // >>
	_Force                       // >!
	_Force1                      // 1>!
	_Force2                      // 2>!
	_Force11                     // 1>|
	_Force22                     // 2>|
	_Ypipe                       // |&
	_TO2                         // >&2
	_1To2                        // 1>&2
	_2To1                        // 2>&1
	_HereDoc                     // <<
)

var replacer = strings.NewReplacer(
	"1>&2", string(_1To2),
	"2>&1", string(_2To1),
	">&2", string(_TO2),
	"1>!", string(_Force1),
	"2>!", string(_Force2),
	"1>|", string(_Force11),
	"2>|", string(_Force22),
	"1>>", string(_Append1),
	"2>>", string(_Append2),
	"0<", string(_Redirect0),
	"1>", string(_Redirect1),
	"2>", string(_Redirect2),
	"&&", string(_AndAlso),
	"||", string(_OrElse),
	">>", string(_Append),
	">!", string(_Force),
	"|&", string(_Ypipe),
	"<<", string(_HereDoc))

var reverse = strings.NewReplacer(
	string(_1To2), "1>&2",
	string(_2To1), "2>&1",
	string(_TO2), ">&2",
	string(_Force1), "1>!",
	string(_Force2), "2>!",
	string(_Force11), "1>|",
	string(_Force22), "2>|",
	string(_Append1), "1>>",
	string(_Append2), "2>>",
	string(_Redirect0), "0<",
	string(_Redirect1), "1>",
	string(_Redirect2), "2>",
	string(_AndAlso), "&&",
	string(_OrElse), "||",
	string(_Append), ">>",
	string(_Force), ">!",
	string(_Ypipe), "|&",
	string(_HereDoc), "<<")

func openSeeNoClobber(fname string) (*os.File, error) {
	if NoClobber {
		return os.OpenFile(fname, os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0666)
	}
	return os.Create(fname)
}

func parse1(stream Stream, text string) ([]*_Statement, error) {
	text = replacer.Replace(text)
	quoteNow := _NotQuoted
	yenCount := 0
	statements := make([]*_Statement, 0)
	args := make([]string, 0)
	rawArgs := make([]string, 0)
	lastchar := ' '
	var buffer bytes.Buffer
	var toDoNextWord func(string)

	toDoRedirect := make([]func([]*os.File) (func(), error), 0, 3)

	termLine := func(term string) {
		statement1 := new(_Statement)
		if buffer.Len() > 0 {
			if toDoNextWord != nil {
				toDoNextWord(buffer.String())
				toDoNextWord = nil
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
		statement1.Redirect = toDoRedirect
		statement1.Term = term
		statements = append(statements, statement1)

		toDoRedirect = make([]func([]*os.File) (func(), error), 0, 3)
		rawArgs = make([]string, 0)
		args = make([]string, 0)
	}

	termWord := func() {
		if toDoNextWord != nil {
			toDoNextWord(buffer.String())
			toDoNextWord = nil
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
		if quoteNow == _NotQuoted {
			if yenCount%2 == 0 && (ch == '"' || ch == '\'') {
				quoteNow = ch
			}
		} else if yenCount%2 == 0 && ch == quoteNow {
			quoteNow = _NotQuoted
		}
		if quoteNow != _NotQuoted {
			buffer.WriteRune(ch)
		} else if isSpace(ch) {
			if buffer.Len() > 0 {
				termWord()
			}
		} else if isSpace(lastchar) && ch == '#' {
			break
		} else if isSpace(lastchar) && ch == ';' {
			termLine(";")
		} else if ch == _OrElse {
			termLine("||")
		} else if ch == '|' {
			termLine("|")
		} else if ch == '&' {
			termLine("&")
		} else if ch == _AndAlso {
			termLine("&&")
		} else if ch == _Ypipe {
			termLine("|&")
		} else if ch == _2To1 {
			termWord()
			toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
				fds[2] = fds[1]
				return func() {}, nil
			})
		} else if ch == _1To2 || ch == _TO2 {
			termWord()
			toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
				fds[1] = fds[2]
				return func() {}, nil
			})
		} else if ch == _HereDoc {
			termWord()

			toDoNextWord = func(word string) {
				dontExpandEnv := (word[0] == '"')
				word = string2word(word, true)
				toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
					lines := make([]string, 0, 20)
					prompt := os.Getenv("PROMPT")
					if dontExpandEnv {
						os.Setenv("PROMPT", fmt.Sprintf("\"%s\">", word))
					} else {
						os.Setenv("PROMPT", fmt.Sprintf("%s>", word))
					}
					defer os.Setenv("PROMPT", prompt)
					ctx := context.Background()
					backup := stream.DisableHistory(true)
					defer stream.DisableHistory(backup)
					for {
						_, line, err := stream.ReadLine(ctx)
						if err != nil {
							if err != io.EOF {
								return func() {}, err
							}
							break
						}
						if strings.HasPrefix(line, word) {
							break
						}
						if !dontExpandEnv {
							line = rxPercent.ReplaceAllStringFunc(
								line,
								func(s string) string {
									name := s[1 : len(s)-1]
									if val, ok := OurGetEnv(name); ok {
										return val
									}
									return s
								})
						}
						lines = append(lines, line)
					}

					r, w, err := os.Pipe()
					if err != nil {
						return func() {}, err
					}
					fds[0] = r

					go func() {
						for _, line := range lines {
							fmt.Fprintln(w, line)
						}
						w.Close()
					}()
					return func() { r.Close() }, nil
				})
			}
		} else if ch == '<' || ch == _Redirect0 {
			termWord()
			toDoNextWord = func(word string) {
				word = string2word(word, true)
				toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
					fd, err := os.Open(word)
					if err != nil {
						return func() {}, err
					}
					fds[0] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == '>' || ch == _Redirect1 {
			termWord()
			toDoNextWord = func(word string) {
				word = string2word(word, true)
				toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
					fd, err := openSeeNoClobber(word)
					if err != nil {
						return func() {}, err
					}
					fds[1] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _Force || ch == _Force1 || ch == _Force11 {
			termWord()
			toDoNextWord = func(word string) {
				word = string2word(word, true)
				toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
					fd, err := os.Create(word)
					if err != nil {
						return func() {}, err
					}
					fds[1] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _Redirect2 {
			termWord()
			toDoNextWord = func(word string) {
				word = string2word(word, true)
				toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
					fd, err := openSeeNoClobber(word)
					if err != nil {
						return func() {}, err
					}
					fds[2] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _Force2 || ch == _Force22 {
			termWord()
			toDoNextWord = func(word string) {
				word = string2word(word, true)
				toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
					fd, err := os.Create(word)
					if err != nil {
						return func() {}, err
					}
					fds[2] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _Append || ch == _Append1 {
			termWord()
			toDoNextWord = func(word string) {
				word = string2word(word, true)
				toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
					fd, err := os.OpenFile(word, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
					if err != nil {
						return func() {}, err
					}
					fds[1] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _Append2 {
			termWord()
			toDoNextWord = func(word string) {
				word = string2word(word, true)
				toDoRedirect = append(toDoRedirect, func(fds []*os.File) (func(), error) {
					fd, err := os.OpenFile(word, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
					if err != nil {
						return func() {}, err
					}
					fds[2] = fd
					return func() { fd.Close() }, nil
				})
			}
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
	termLine(" ")
	return statements, nil
}

// Make arrays whose elements are pipelines
func parse2(statements []*_Statement) [][]*_Statement {
	result := make([][]*_Statement, 1)
	for _, statement1 := range statements {
		result[len(result)-1] = append(result[len(result)-1], statement1)
		switch statement1.Term {
		case "|", "|&":

		default:
			result = append(result, make([]*_Statement, 0))
		}
	}
	if len(result[len(result)-1]) <= 0 {
		result = result[0 : len(result)-1]
	}
	return result
}

// Parse parses the string and make Statement objects.
func Parse(stream Stream, text string) ([][]*_Statement, error) {
	result1, err := parse1(stream, text)
	if err != nil {
		return nil, err
	}
	result2 := parse2(result1)
	return result2, nil
}
