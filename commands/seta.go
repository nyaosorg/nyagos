package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	lshift = '\uE000' + iota
	rshift
	assign_add
	assign_sub
	assign_mul
	assign_div
	assign_mod
	assign_and
	assign_or
	assign_xor
	assign_lshift
	assign_rshift
)

var replacer = strings.NewReplacer(
	"<<=", string(assign_lshift),
	">>=", string(assign_rshift),
	"<<", string(lshift),
	">>", string(rshift),
	"+=", string(assign_add),
	"-=", string(assign_sub),
	"*=", string(assign_mul),
	"/=", string(assign_div),
	"%=", string(assign_mod),
	"&=", string(assign_and),
	"|=", string(assign_or),
	"^=", string(assign_xor))

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

type value_t interface {
	Get() int
}

type rvalue_t int

func (n rvalue_t) Get() int {
	return int(n)
}

type lvalue_t string

func (name lvalue_t) Get() int {
	n, err := strconv.Atoi(os.Getenv(string(name)))
	if err != nil {
		return 0
	} else {
		return n
	}
}

func (name lvalue_t) Set(value int) {
	os.Setenv(string(name), fmt.Sprintf("%d", value))
}

func readValue(r io.RuneScanner) (value_t, error) {
	if err := skipSpace(r); err != nil {
		return rvalue_t(1), err
	}
	ch, _, err := r.ReadRune()
	if err != nil {
		return rvalue_t(1), err
	}
	if ch == '-' {
		value, err := readValue(r)
		return rvalue_t(-value.Get()), err
	}
	if ch == '+' {
		return readValue(r)
	}
	if ch == '~' {
		value, err := readValue(r)
		return rvalue_t(^value.Get()), err
	}
	if ch == '!' {
		value, err := readValue(r)
		if value.Get() != 0 {
			return rvalue_t(0), err
		} else {
			return rvalue_t(1), err
		}
	}
	if ch == '0' {
		value := 0
		ch, _, err := r.ReadRune()
		if err != nil {
			return rvalue_t(0), err
		}
		if ch == 'x' || ch == 'X' {
			for {
				ch, _, err := r.ReadRune()
				if err != nil {
					return rvalue_t(value), err
				}
				m := strings.IndexRune("0123456789ABCDEFF", unicode.ToUpper(ch))
				if m < 0 {
					r.UnreadRune()
					return rvalue_t(value), nil
				}
				value = value*16 + m
			}
		} else {
			for {
				m := strings.IndexRune("01234567", ch)
				if m < 0 {
					r.UnreadRune()
					return rvalue_t(value), nil
				}
				value = value*8 + m
				ch, _, err = r.ReadRune()
				if err != nil {
					return rvalue_t(value), err
				}
			}
		}
	}
	if n := strings.IndexRune("0123456789", ch); n >= 0 {
		value := n
		for {
			ch, _, err := r.ReadRune()
			if err != nil {
				return rvalue_t(value), err
			}
			m := strings.IndexRune("0123456789", ch)
			if m < 0 {
				r.UnreadRune()
				return rvalue_t(value), nil
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
		return lvalue_t(name.String()), nil
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
	return rvalue_t(1), errors.New("syntax error")
}

type operation struct {
	Operator map[rune]func(v1, v2 value_t) value_t
	Sub      func(io.RuneScanner) (value_t, error)
}

func (this *operation) Eval(r io.RuneScanner) (value_t, error) {
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

func readEquation(r io.RuneScanner) (value_t, error) {
	if opComma == nil {
		opMulDiv := &operation{
			Operator: map[rune]func(v1, v2 value_t) value_t{
				'*': func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() * v2.Get())
				},
				'/': func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() / v2.Get())
				},
				'%': func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() % v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (value_t, error) {
				return readValue(r)
			},
		}
		opAddSub := &operation{
			Operator: map[rune]func(v1, v2 value_t) value_t{
				'+': func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() + v2.Get())
				},
				'-': func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() - v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (value_t, error) {
				return opMulDiv.Eval(r)
			},
		}
		opShift := &operation{
			Operator: map[rune]func(v1, v2 value_t) value_t{
				lshift: func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() << uint(v2.Get()))
				},
				rshift: func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() >> uint(v2.Get()))
				},
			},
			Sub: func(r io.RuneScanner) (value_t, error) {
				return opAddSub.Eval(r)
			},
		}
		opBitAnd := &operation{
			Operator: map[rune]func(v1, v2 value_t) value_t{
				'&': func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() & v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (value_t, error) {
				return opShift.Eval(r)
			},
		}
		opBitXor := &operation{
			Operator: map[rune]func(v1, v2 value_t) value_t{
				'^': func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() ^ v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (value_t, error) {
				return opBitAnd.Eval(r)
			},
		}
		opBitOr := &operation{
			Operator: map[rune]func(v1, v2 value_t) value_t{
				'|': func(v1, v2 value_t) value_t {
					return rvalue_t(v1.Get() | v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (value_t, error) {
				return opBitXor.Eval(r)
			},
		}
		opAssign := &operation{
			Operator: map[rune]func(v1, v2 value_t) value_t{
				'=': func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v2.Get())
						return V1
					}
					return v2
				},
				assign_add: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() + v2.Get())
						return V1
					}
					return v2
				},
				assign_sub: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() - v2.Get())
						return V1
					}
					return v2
				},
				assign_mul: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() * v2.Get())
						return V1
					}
					return v2
				},
				assign_div: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() / v2.Get())
						return V1
					}
					return v2
				},
				assign_mod: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() % v2.Get())
						return V1
					}
					return v2
				},
				assign_and: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() & v2.Get())
						return V1
					}
					return v2
				},
				assign_or: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() | v2.Get())
						return V1
					}
					return v2
				},
				assign_xor: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() ^ v2.Get())
						return V1
					}
					return v2
				},
				assign_lshift: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() << uint(v2.Get()))
						return V1
					}
					return v2
				},
				assign_rshift: func(v1, v2 value_t) value_t {
					if V1, ok := v1.(lvalue_t); ok {
						V1.Set(v1.Get() >> uint(v2.Get()))
						return V1
					}
					return v2
				},
			},
			Sub: func(r io.RuneScanner) (value_t, error) {
				return opBitOr.Eval(r)
			},
		}
		opComma = &operation{
			Operator: map[rune]func(v1, v2 value_t) value_t{
				',': func(v1, v2 value_t) value_t {
					return v2
				},
			},
			Sub: func(r io.RuneScanner) (value_t, error) {
				return opAssign.Eval(r)
			},
		}
	}
	return opComma.Eval(r)
}

func evalEquation(s string) (int, error) {
	value, err := readEquation(strings.NewReader(replacer.Replace(s)))
	if err == io.EOF {
		err = nil
	}
	return value.Get(), err
}
