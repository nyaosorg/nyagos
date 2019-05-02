package commands

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

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

func readMulDiv(r io.RuneScanner) (int, error) {
	value, err := readValue(r)
	if err != nil {
		return value, err
	}
	for {
		if err := skipSpace(r); err != nil {
			return value, nil
		}
		ch, _, err := r.ReadRune()
		if err != nil {
			return value, err
		}
		if ch == '*' {
			value2, err := readValue(r)
			value *= value2
			if err != nil {
				return value, err
			}
		} else if ch == '/' {
			value2, err := readValue(r)
			value /= value2
			if err != nil {
				return value, err
			}
		} else if ch == '%' {
			value2, err := readValue(r)
			value %= value2
			if err != nil {
				return value, err
			}
		} else {
			r.UnreadRune()
			return value, nil
		}
	}
}

func readAddSub(r io.RuneScanner) (int, error) {
	value, err := readMulDiv(r)
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
		if ch == '+' {
			value2, err := readMulDiv(r)
			value += value2
			if err != nil {
				return value, err
			}
		} else if ch == '-' {
			value2, err := readMulDiv(r)
			value -= value2
			if err != nil {
				return value, err
			}
		} else {
			r.UnreadRune()
			return value, nil
		}
	}
}

func readEquation(r io.RuneScanner) (int, error) {
	return readAddSub(r)
}

func evalEquation(s string) (int, error) {
	value, err := readEquation(strings.NewReader(s))
	if err == io.EOF {
		err = nil
	}
	return value, err
}
