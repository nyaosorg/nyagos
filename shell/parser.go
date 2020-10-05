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

	"github.com/zetamatta/nyagos/nodos"
	"github.com/zetamatta/nyagos/texts"
)

var NoClobber = false

func isSpace(c rune) bool {
	return strings.IndexRune(" \t\n\r\v\f", c) >= 0
}

type StatementT struct {
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
		} else {
			return ""
		}
	},
	"ERRORLEVEL": func() string {
		return fmt.Sprintf("%d", LastErrorLevel)
	},
	"DATE": func() string {
		layout, err := nodos.OsDateLayout()
		if err != nil {
			return err.Error()
		}
		return time.Now().Format(layout)
	},
	"TIME": func() string {
		return time.Now().Format("15:04:05.00")
	},
}

var rxUnicode = regexp.MustCompile("^[uU]\\+?([0-9a-fA-F]+)$")

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
		return texts.ReplaceIgnoreCase(base, m[2], m[3]), true
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

func chomp(buffer *bytes.Buffer) {
	buffer.Truncate(buffer.Len() - 1)
}

const NOTQUOTED = '\000'

var TildeExpansion = true

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
		if TildeExpansion && ch == '~' && isSpace(lastchar) && quoteNow == NOTQUOTED {
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
					if quoteNow == NOTQUOTED {
						quoteNow = ch
					} else {
						quoteNow = NOTQUOTED
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
					if !removeQuote && strings.Count(undo.String(), `"`)%2 != 0 {
						buffer.WriteByte('"')
					}
					buffer.WriteString(u.HomeDir)
					lastchar = rune(u.HomeDir[len(u.HomeDir)-1])
				} else {
					if !removeQuote {
						buffer.WriteByte('~')
					}
					undoStr := undo.String()
					buffer.WriteString(undoStr)
					lastchar = rune(undoStr[len(undoStr)-1])
				}
				continue
			}
			if home := nodos.GetHome(); home != "" {
				if !removeQuote && strings.Count(undo.String(), `"`)%2 != 0 {
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

		if quoteNow != NOTQUOTED && ch == quoteNow && yenCount%2 == 0 {
			if !removeQuote {
				buffer.WriteRune(ch)
			}
			// Close Quotation.
			for ; yenCount >= 2; yenCount -= 2 {
				buffer.WriteByte('\\')
			}
			quoteNow = NOTQUOTED
		} else if (ch == '\'' || ch == '"') && quoteNow == NOTQUOTED && yenCount%2 == 0 {
			if !removeQuote {
				buffer.WriteRune(ch)
			}
			// Open Qutation.
			for ; yenCount >= 2; yenCount -= 2 {
				buffer.WriteByte('\\')
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
					buffer.WriteByte('\\')
				}
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
	_ANDALSO   = '\uE000' + iota // &&
	_ORELSE                      // ||
	_REDIRECT0                   // 0<
	_REDIRECT1                   // 1>
	_REDIRECT2                   // 2>
	_APPEND1                     // 1>>
	_APPEND2                     // 2>>
	_APPEND                      // >>
	_FORCE                       // >!
	_FORCE1                      // 1>!
	_FORCE2                      // 2>!
	_FORCE11                     // 1>|
	_FORCE22                     // 2>|
	_YPIPE                       // |&
	_TO2                         // >&2
	_1TO2                        // 1>&2
	_2TO1                        // 2>&1
	_HEREDOC                     // <<
)

var replacer = strings.NewReplacer(
	"1>&2", string(_1TO2),
	"2>&1", string(_2TO1),
	">&2", string(_TO2),
	"1>!", string(_FORCE1),
	"2>!", string(_FORCE2),
	"1>|", string(_FORCE11),
	"2>|", string(_FORCE22),
	"1>>", string(_APPEND1),
	"2>>", string(_APPEND2),
	"0<", string(_REDIRECT0),
	"1>", string(_REDIRECT1),
	"2>", string(_REDIRECT2),
	"&&", string(_ANDALSO),
	"||", string(_ORELSE),
	">>", string(_APPEND),
	">!", string(_FORCE),
	"|&", string(_YPIPE),
	"<<", string(_HEREDOC))

var reverse = strings.NewReplacer(
	string(_1TO2), "1>&2",
	string(_2TO1), "2>&1",
	string(_TO2), ">&2",
	string(_FORCE1), "1>!",
	string(_FORCE2), "2>!",
	string(_FORCE11), "1>|",
	string(_FORCE22), "2>|",
	string(_APPEND1), "1>>",
	string(_APPEND2), "2>>",
	string(_REDIRECT0), "0<",
	string(_REDIRECT1), "1>",
	string(_REDIRECT2), "2>",
	string(_ANDALSO), "&&",
	string(_ORELSE), "||",
	string(_APPEND), ">>",
	string(_FORCE), ">!",
	string(_YPIPE), "|&",
	string(_HEREDOC), "<<")

func openSeeNoClobber(fname string) (*os.File, error) {
	if NoClobber {
		return os.OpenFile(fname, os.O_EXCL|os.O_CREATE, 0666)
	} else {
		return os.Create(fname)
	}
}

func parse1(stream Stream, text string) ([]*StatementT, error) {
	text = replacer.Replace(text)
	quoteNow := NOTQUOTED
	yenCount := 0
	statements := make([]*StatementT, 0)
	args := make([]string, 0)
	rawArgs := make([]string, 0)
	lastchar := ' '
	var buffer bytes.Buffer
	var todo_nextword func(string)

	todo_redirect := make([]func([]*os.File) (func(), error), 0, 3)

	term_line := func(term string) {
		statement1 := new(StatementT)
		if buffer.Len() > 0 {
			if todo_nextword != nil {
				todo_nextword(buffer.String())
				todo_nextword = nil
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
		statement1.Redirect = todo_redirect
		statement1.Term = term
		statements = append(statements, statement1)

		todo_redirect = make([]func([]*os.File) (func(), error), 0, 3)
		rawArgs = make([]string, 0)
		args = make([]string, 0)
	}

	term_word := func() {
		if todo_nextword != nil {
			todo_nextword(buffer.String())
			todo_nextword = nil
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
		} else if isSpace(ch) {
			if buffer.Len() > 0 {
				term_word()
			}
		} else if isSpace(lastchar) && ch == '#' {
			break
		} else if isSpace(lastchar) && ch == ';' {
			term_line(";")
		} else if ch == _ORELSE {
			term_line("||")
		} else if ch == '|' {
			term_line("|")
		} else if ch == '&' {
			term_line("&")
		} else if ch == _ANDALSO {
			term_line("&&")
		} else if ch == _YPIPE {
			term_line("|&")
		} else if ch == _2TO1 {
			term_word()
			todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
				fds[2] = fds[1]
				return func() {}, nil
			})
		} else if ch == _1TO2 || ch == _TO2 {
			term_word()
			todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
				fds[1] = fds[2]
				return func() {}, nil
			})
		} else if ch == _HEREDOC {
			term_word()

			todo_nextword = func(word string) {
				dont_expand_env := (word[0] == '"')
				word = string2word(word, true)
				todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
					lines := make([]string, 0, 20)
					prompt := os.Getenv("PROMPT")
					if dont_expand_env {
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
						if !dont_expand_env {
							line = rxPercent.ReplaceAllStringFunc(
								line,
								func(s string) string {
									name := s[1 : len(s)-1]
									if val, ok := OurGetEnv(name); ok {
										return val
									} else {
										return s
									}
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
		} else if ch == '<' || ch == _REDIRECT0 {
			term_word()
			todo_nextword = func(word string) {
				word = string2word(word, true)
				todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
					fd, err := os.Open(word)
					if err != nil {
						return func() {}, err
					}
					fds[0] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == '>' || ch == _REDIRECT1 {
			term_word()
			todo_nextword = func(word string) {
				word = string2word(word, true)
				todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
					fd, err := openSeeNoClobber(word)
					if err != nil {
						return func() {}, err
					}
					fds[1] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _FORCE || ch == _FORCE1 || ch == _FORCE11 {
			term_word()
			todo_nextword = func(word string) {
				word = string2word(word, true)
				todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
					fd, err := os.Create(word)
					if err != nil {
						return func() {}, err
					}
					fds[1] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _REDIRECT2 {
			term_word()
			todo_nextword = func(word string) {
				word = string2word(word, true)
				todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
					fd, err := openSeeNoClobber(word)
					if err != nil {
						return func() {}, err
					}
					fds[2] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _FORCE2 || ch == _FORCE22 {
			term_word()
			todo_nextword = func(word string) {
				word = string2word(word, true)
				todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
					fd, err := os.Create(word)
					if err != nil {
						return func() {}, err
					}
					fds[2] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _APPEND || ch == _APPEND1 {
			term_word()
			todo_nextword = func(word string) {
				word = string2word(word, true)
				todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
					fd, err := os.OpenFile(word, os.O_APPEND|os.O_CREATE, 0666)
					if err != nil {
						return func() {}, err
					}
					fds[1] = fd
					return func() { fd.Close() }, nil
				})
			}
		} else if ch == _APPEND2 {
			term_word()
			todo_nextword = func(word string) {
				word = string2word(word, true)
				todo_redirect = append(todo_redirect, func(fds []*os.File) (func(), error) {
					fd, err := os.OpenFile(word, os.O_APPEND|os.O_CREATE, 0666)
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

func parse(stream Stream, text string) ([][]*StatementT, error) {
	result1, err := parse1(stream, text)
	if err != nil {
		return nil, err
	}
	result2 := parse2(result1)
	return result2, nil
}
