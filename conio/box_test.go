package conio

import (
	"os"
	"testing"
)

func TestBoxPrint(t *testing.T) {
	BoxPrint([]string{
		"aaaa", "bbbb", "cccc", "fjdaksljflkdajfkljsalkfjdlkf",
		"jfkldsjflkjdsalkfjlkdsajflkajds",
		"fsdfsdf"}, os.Stdout)
}
