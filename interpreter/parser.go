package interpreter

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"unicode"

	"../dos"
)

type RedirectT struct {
	Path     string
	IsAppend bool
}

type StatementT struct {
	Argv     []string
	Redirect [3]RedirectT
	IsAppend [3]bool
	Term     string
}

var prefix []string = []string{" 0<", " 1>", " 2>"}

var percentFunc = map[string]func() string{
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

func (this StatementT) String() string {
	var buffer bytes.Buffer
	for _, arg := range this.Argv {
		buffer.WriteRune('[')
		buffer.WriteString(arg)
		buffer.WriteRune(']')
	}
	for i := 0; i < len(prefix); i++ {
		if len(this.Redirect[i].Path) > 0 {
			buffer.WriteString(prefix[i])
			buffer.WriteString("[")
			buffer.WriteString(this.Redirect[i].Path)
			buffer.WriteString("]")
		}
	}
	buffer.WriteString(" ")
	buffer.WriteString(this.Term)
	return buffer.String()
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
			} else if f, ok := percentFunc[nameStr]; ok {
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
	nextword *int,
	redirect *[3]RedirectT,
	buffer *bytes.Buffer,
	argv *[]string,
	term string) {
	var statement1 StatementT
	if buffer.Len() > 0 {
		if *nextword == WORD_ARGV {
			statement1.Argv = append(*argv, dequote(buffer))
		} else {
			statement1.Argv = *argv
			(*redirect)[*nextword].Path = dequote(buffer)
			*nextword = WORD_ARGV
		}
		buffer.Reset()
	} else if len(*argv) <= 0 {
		return
	} else {
		statement1.Argv = *argv
	}
	statement1.Redirect[0] = redirect[0]
	statement1.Redirect[1] = redirect[1]
	statement1.Redirect[2] = redirect[2]
	redirect[0].Path = ""
	redirect[0].IsAppend = false
	redirect[1].Path = ""
	redirect[1].IsAppend = false
	redirect[2].Path = ""
	redirect[2].IsAppend = false
	*argv = make([]string, 0)
	statement1.Term = term
	*statements = append(*statements, statement1)
}

const (
	WORD_ARGV   = -1
	WORD_STDIN  = 0
	WORD_STDOUT = 1
	WORD_STDERR = 2
)

func Parse1(text string) []StatementT {
	isQuoted := false
	statements := make([]StatementT, 0)
	argv := make([]string, 0)
	lastchar := ' '
	lastredirected := -1
	var buffer bytes.Buffer
	nextword := WORD_ARGV
	var redirect [3]RedirectT
	for _, ch := range text {
		if ch == '"' {
			isQuoted = !isQuoted
		}
		if isQuoted {
			buffer.WriteRune(ch)
		} else {
			if ch == ' ' {
				if buffer.Len() > 0 {
					if nextword == WORD_ARGV {
						argv = append(argv, dequote(&buffer))
					} else {
						redirect[nextword].Path = dequote(&buffer)
					}
					buffer.Reset()
					nextword = WORD_ARGV
				}
			} else if lastchar == ' ' && ch == ';' {
				terminate(&statements, &nextword, &redirect, &buffer, &argv, ";")
			} else if ch == '|' {
				if lastchar == '|' {
					statements[len(statements)-1].Term = "||"
				} else {
					terminate(&statements, &nextword, &redirect, &buffer, &argv, "|")
				}
			} else if ch == '&' {
				if lastchar == '&' {
					statements[len(statements)-1].Term = "&&"
				} else {
					terminate(&statements, &nextword, &redirect, &buffer, &argv, "&")
				}
			} else if ch == '>' {
				if lastchar == '1' {
					chomp(&buffer)
					nextword = WORD_STDOUT
					redirect[1].IsAppend = false
					lastredirected = 1
				} else if lastchar == '2' {
					chomp(&buffer)
					nextword = WORD_STDERR
					redirect[2].IsAppend = false
					lastredirected = 2
				} else if lastchar == '>' && lastredirected >= 0 {
					redirect[lastredirected].IsAppend = true
				} else {
					nextword = WORD_STDOUT
					lastredirected = 1
				}
			} else if ch == '<' {
				nextword = WORD_STDIN
				redirect[0].IsAppend = false
				lastredirected = 0
			} else {
				buffer.WriteRune(ch)
			}
		}
		lastchar = ch
	}
	terminate(&statements, &nextword, &redirect, &buffer, &argv, " ")
	return statements
}

func Parse2(statements []StatementT) [][]StatementT {
	result := make([][]StatementT, 1)
	for _, statement1 := range statements {
		result[len(result)-1] = append(result[len(result)-1], statement1)
		if statement1.Term != "|" {
			result = append(result, make([]StatementT, 0))
		}
	}
	if len(result[len(result)-1]) <= 0 {
		result = result[0 : len(result)-1]
	}
	return result
}

func Parse(text string) [][]StatementT {
	result1 := Parse1(text)
	result2 := Parse2(result1)
	return result2
}
