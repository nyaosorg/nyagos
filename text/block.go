package text

import (
	"strings"
)

func ReadBlock(reader func() (string, error), unreader func(string)) []string {
	kekka_count := 0
	result := make([]string, 0, 10)
	for {
		line, err := reader()
		if err != nil {
			return result
		}
		args := SplitQ(line)
		for i, arg1 := range args {
			if arg1 == "(" {
				kekka_count++
			} else if arg1 == ")" {
				kekka_count--
				if kekka_count < 0 {
					if i > 0 {
						result = append(result, strings.Join(args[:i], " "))
					}
					if i+1 < len(args) {
						unreader(strings.Join(args[i+1:], " "))
					}
					return result
				}
			}
		}
		if len(args) > 0 {
			result = append(result, strings.Join(args, " "))
		}
	}
}
