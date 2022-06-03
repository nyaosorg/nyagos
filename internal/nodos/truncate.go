package nodos

import (
	"io"
)

func Truncate(folder string, whenError func(string, error) bool, out io.Writer) error {
	return truncate(folder, whenError, out)
}
