package commands

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	lshift = '\uE000' + iota
	rshift
)

var replacer = strings.NewReplacer(
	"<<", string(lshift),
	">>", string(rshift))

func skipSpace(r io.RuneScanner) error {
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		if strings.IndexRune(" \t\v\r\n", ch) < 0 {
			r.UnreadRune()
			return nil
		}
	}
}

func readValue(r io.RuneScanner) (int, error) {
	if err := skipSpace(r); err != nil {
		return 1, err
	}
	ch, _, err := r.ReadRune()
	if err != nil {
		return 1, err
	}
	if ch == '-' {
		value, err := readValue(r)
		return -value, err
	}
	if ch == '+' {
		return readValue(r)
	}
	if ch == '~' {
		value, err := readValue(r)
		return ^value, err
	}
	if ch == '!' {
		value, err := readValue(r)
		if value != 0 {
			return 0, err
		} else {
			return 1, err
		}
	}
	if n := strings.IndexRune("0123456789", ch); n >= 0 {
		value := n
		for {
			ch, _, err := r.ReadRune()
			if err != nil {
				return value, err
			}
			m := strings.IndexRune("0123456789", ch)
			if m < 0 {
				r.UnreadRune()
				return value, nil
			}
			value = value*10 + m
		}
	}
	if unicode.IsLetter(ch) {
		var name strings.Builder
		for {
			name.WriteRune(ch)
			ch, _, err = r.ReadRune()
			if err != nil {
				break
			}
			if !unicode.IsLetter(ch) {
				r.UnreadRune()
				break
			}
		}
		envValue := os.Getenv(name.String())
		return strconv.Atoi(envValue)
	}
	if ch == '(' {
		value, err := readEquation(r)
		if err == nil {
			ch, _, err := r.ReadRune()
			if err != nil || ch != ')' {
				return value, errors.New("() pair is not closed")
			}
			return value, nil
		}
		return value, errors.New("() pair is not closed")
	}
	return 1, errors.New("syntax error")
}

type operation struct {
	Operator map[rune]func(v1, v2 int) int
	Sub      func(io.RuneScanner) (int, error)
}

func (this *operation) Eval(r io.RuneScanner) (int, error) {
	value, err := this.Sub(r)
	if err != nil {
		return value, err
	}
	for {
		if err := skipSpace(r); err != nil {
			return value, err
		}
		ch, _, err := r.ReadRune()
		if err != nil {
			return value, err
		}
		if f := this.Operator[ch]; f != nil {
			value2, err := this.Sub(r)
			value = f(value, value2)
			if err != nil {
				return value, err
			}
		} else {
			r.UnreadRune()
			return value, nil
		}
	}
}

var opComma *operation

func readEquation(r io.RuneScanner) (int, error) {
	if opComma == nil {
		opMulDiv := &operation{
			Operator: map[rune]func(v1, v2 int) int{
				'*': func(v1, v2 int) int { return v1 * v2 },
				'/': func(v1, v2 int) int { return v1 / v2 },
				'%': func(v1, v2 int) int { return v1 % v2 },
			},
			Sub: func(r io.RuneScanner) (int, error) { return readValue(r) },
		}
		opAddSub := &operation{
			Operator: map[rune]func(v1, v2 int) int{
				'+': func(v1, v2 int) int { return v1 + v2 },
				'-': func(v1, v2 int) int { return v1 - v2 },
			},
			Sub: func(r io.RuneScanner) (int, error) { return opMulDiv.Eval(r) },
		}
		opShift := &operation{
			Operator: map[rune]func(v1, v2 int) int{
				lshift: func(v1, v2 int) int { return v1 << uint(v2) },
				rshift: func(v1, v2 int) int { return v1 >> uint(v2) },
			},
			Sub: func(r io.RuneScanner) (int, error) { return opAddSub.Eval(r) },
		}
		opBitAnd := &operation{
			Operator: map[rune]func(v1, v2 int) int{
				'&': func(v1, v2 int) int { return v1 & v2 },
			},
			Sub: func(r io.RuneScanner) (int, error) { return opShift.Eval(r) },
		}
		opBitXor := &operation{
			Operator: map[rune]func(v1, v2 int) int{
				'^': func(v1, v2 int) int { return v1 ^ v2 },
			},
			Sub: func(r io.RuneScanner) (int, error) { return opBitAnd.Eval(r) },
		}
		opBitOr := &operation{
			Operator: map[rune]func(v1, v2 int) int{
				'|': func(v1, v2 int) int { return v1 | v2 },
			},
			Sub: func(r io.RuneScanner) (int, error) { return opBitXor.Eval(r) },
		}
		opComma = &operation{
			Operator: map[rune]func(v1, v2 int) int{
				',': func(v1, v2 int) int { return v2 },
			},
			Sub: func(r io.RuneScanner) (int, error) { return opBitOr.Eval(r) },
		}
	}
	return opComma.Eval(r)
}

func evalEquation(s string) (int, error) {
	value, err := readEquation(strings.NewReader(replacer.Replace(s)))
	if err == io.EOF {
		err = nil
	}
	return value, err
}
