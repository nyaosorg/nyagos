package shell_test

import (
	"testing"

	"github.com/zetamatta/nyagos/shell"
)

func TestQuote(t *testing.T) {
	testCases := [][2]string{
		{`123`, `"123"`},
		{`123"456`, `"123\"456"`},
		{`foo\bar\`, `"foo\bar\\"`},
		{`foo\"bar`, `"foo\\\"bar"`},
	}

	for _, testCase1 := range testCases {
		source := testCase1[0]
		expect := testCase1[1]
		actual := shell.Quote(source)
		if actual != expect {
			t.Fatalf("Quote: expect `%s` for `%s`, but `%s`", expect, source, actual)
		}
	}
}
