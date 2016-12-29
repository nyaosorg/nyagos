package text

import (
	"io"
	"strings"
	"testing"
)

func TestReadBlock(t *testing.T) {
	sample := []string{
		"ahaha",
		"( ihihi )",
		"ufufu ehehe )",
		"ohoho",
	}

	results := ReadBlock(func() (result string, err error) {
		if len(sample) <= 0 {
			result = ""
			err = io.EOF
		} else {
			result = sample[0]
			sample = sample[1:]
			err = nil
		}
		return
	}, func(line string) {
		temp := make([]string, 0, len(sample)+1)
		temp = append(temp, line)
		sample = append(temp, sample...)
	})

	if len(results) != 3 ||
		results[0] != "ahaha" ||
		results[1] != "( ihihi )" ||
		results[2] != "ufufu ehehe" {

		t.Fail()
	}
	if sample[0] != "ohoho" {
		t.Fail()
	}

	sample2 := []string{
		"11111",
		"(",
		"22222",
		")",
		"33333",
		")",
		"44444",
	}
	result2 := ReadBlock(func() (result string, err error) {
		if len(sample2) <= 0 {
			result = ""
			err = io.EOF
		} else {
			result = sample2[0]
			sample2 = sample2[1:]
			err = nil
		}
		return
	}, func(line string) {
		temp := make([]string, 0, len(sample2)+1)
		temp = append(temp, line)
		sample2 = append(temp, sample2...)
	})

	if len(result2) != 5 ||
		result2[0] != "11111" ||
		result2[1] != "(" ||
		result2[2] != "22222" ||
		result2[3] != ")" ||
		result2[4] != "33333" ||
		len(sample2) != 1 ||
		sample2[0] != "44444" {

		t.Log("result2=" + strings.Join(result2, "\n"))
		t.Log("sample2=" + strings.Join(sample2, "\n"))

		t.Fail()
	}
}
