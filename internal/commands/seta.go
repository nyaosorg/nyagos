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
	assignAdd
	assignSub
	assignMul
	assignDiv
	assignMod
	assignAnd
	assignOr
	assignXor
	assignLshift
	assignRshift
)

var replacer = strings.NewReplacer(
	"<<=", string(assignLshift),
	">>=", string(assignRshift),
	"<<", string(lshift),
	">>", string(rshift),
	"+=", string(assignAdd),
	"-=", string(assignSub),
	"*=", string(assignMul),
	"/=", string(assignDiv),
	"%=", string(assignMod),
	"&=", string(assignAnd),
	"|=", string(assignOr),
	"^=", string(assignXor))

func skipSpace(r io.RuneScanner) error {
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		if !strings.ContainsRune(" \t\v\r\n", ch) {
			r.UnreadRune()
			return nil
		}
	}
}

type valueT interface {
	Get() int
}

type rValueT int

func (n rValueT) Get() int {
	return int(n)
}

type lValueT string

func (name lValueT) Get() int {
	n, err := strconv.Atoi(os.Getenv(string(name)))
	if err != nil {
		return 0
	}
	return n
}

func (name lValueT) Set(value int) {
	os.Setenv(string(name), fmt.Sprintf("%d", value))
}

func readValue(r io.RuneScanner) (valueT, error) {
	if err := skipSpace(r); err != nil {
		return rValueT(1), err
	}
	ch, _, err := r.ReadRune()
	if err != nil {
		return rValueT(1), err
	}
	if ch == '-' {
		value, err := readValue(r)
		return rValueT(-value.Get()), err
	}
	if ch == '+' {
		return readValue(r)
	}
	if ch == '~' {
		value, err := readValue(r)
		return rValueT(^value.Get()), err
	}
	if ch == '!' {
		value, err := readValue(r)
		if value.Get() != 0 {
			return rValueT(0), err
		}
		return rValueT(1), err
	}
	if ch == '0' {
		value := 0
		ch, _, err := r.ReadRune()
		if err != nil {
			return rValueT(0), err
		}
		if ch == 'x' || ch == 'X' {
			for {
				ch, _, err := r.ReadRune()
				if err != nil {
					return rValueT(value), err
				}
				m := strings.IndexRune("0123456789ABCDEFF", unicode.ToUpper(ch))
				if m < 0 {
					r.UnreadRune()
					return rValueT(value), nil
				}
				value = value*16 + m
			}
		} else {
			for {
				m := strings.IndexRune("01234567", ch)
				if m < 0 {
					r.UnreadRune()
					return rValueT(value), nil
				}
				value = value*8 + m
				ch, _, err = r.ReadRune()
				if err != nil {
					return rValueT(value), err
				}
			}
		}
	}
	if n := strings.IndexRune("0123456789", ch); n >= 0 {
		value := n
		for {
			ch, _, err := r.ReadRune()
			if err != nil {
				return rValueT(value), err
			}
			m := strings.IndexRune("0123456789", ch)
			if m < 0 {
				r.UnreadRune()
				return rValueT(value), nil
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
		return lValueT(name.String()), nil
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
	return rValueT(1), errors.New("syntax error")
}

type operation struct {
	Operator map[rune]func(v1, v2 valueT) valueT
	Sub      func(io.RuneScanner) (valueT, error)
}

func (op *operation) Eval(r io.RuneScanner) (valueT, error) {
	value, err := op.Sub(r)
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
		if f := op.Operator[ch]; f != nil {
			value2, err := op.Sub(r)
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

func readEquation(r io.RuneScanner) (valueT, error) {
	if opComma == nil {
		opMulDiv := &operation{
			Operator: map[rune]func(v1, v2 valueT) valueT{
				'*': func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() * v2.Get())
				},
				'/': func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() / v2.Get())
				},
				'%': func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() % v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (valueT, error) {
				return readValue(r)
			},
		}
		opAddSub := &operation{
			Operator: map[rune]func(v1, v2 valueT) valueT{
				'+': func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() + v2.Get())
				},
				'-': func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() - v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (valueT, error) {
				return opMulDiv.Eval(r)
			},
		}
		opShift := &operation{
			Operator: map[rune]func(v1, v2 valueT) valueT{
				lshift: func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() << uint(v2.Get()))
				},
				rshift: func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() >> uint(v2.Get()))
				},
			},
			Sub: func(r io.RuneScanner) (valueT, error) {
				return opAddSub.Eval(r)
			},
		}
		opBitAnd := &operation{
			Operator: map[rune]func(v1, v2 valueT) valueT{
				'&': func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() & v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (valueT, error) {
				return opShift.Eval(r)
			},
		}
		opBitXor := &operation{
			Operator: map[rune]func(v1, v2 valueT) valueT{
				'^': func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() ^ v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (valueT, error) {
				return opBitAnd.Eval(r)
			},
		}
		opBitOr := &operation{
			Operator: map[rune]func(v1, v2 valueT) valueT{
				'|': func(v1, v2 valueT) valueT {
					return rValueT(v1.Get() | v2.Get())
				},
			},
			Sub: func(r io.RuneScanner) (valueT, error) {
				return opBitXor.Eval(r)
			},
		}
		opAssign := &operation{
			Operator: map[rune]func(v1, v2 valueT) valueT{
				'=': func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v2.Get())
						return V1
					}
					return v2
				},
				assignAdd: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() + v2.Get())
						return V1
					}
					return v2
				},
				assignSub: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() - v2.Get())
						return V1
					}
					return v2
				},
				assignMul: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() * v2.Get())
						return V1
					}
					return v2
				},
				assignDiv: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() / v2.Get())
						return V1
					}
					return v2
				},
				assignMod: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() % v2.Get())
						return V1
					}
					return v2
				},
				assignAnd: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() & v2.Get())
						return V1
					}
					return v2
				},
				assignOr: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() | v2.Get())
						return V1
					}
					return v2
				},
				assignXor: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() ^ v2.Get())
						return V1
					}
					return v2
				},
				assignLshift: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() << uint(v2.Get()))
						return V1
					}
					return v2
				},
				assignRshift: func(v1, v2 valueT) valueT {
					if V1, ok := v1.(lValueT); ok {
						V1.Set(v1.Get() >> uint(v2.Get()))
						return V1
					}
					return v2
				},
			},
			Sub: func(r io.RuneScanner) (valueT, error) {
				return opBitOr.Eval(r)
			},
		}
		opComma = &operation{
			Operator: map[rune]func(v1, v2 valueT) valueT{
				',': func(v1, v2 valueT) valueT {
					return v2
				},
			},
			Sub: func(r io.RuneScanner) (valueT, error) {
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
