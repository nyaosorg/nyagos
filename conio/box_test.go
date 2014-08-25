package conio

import "os"
import "testing"

func TestBoxPrint(t *testing.T) {
	BoxPrint([]string{
		"aaaa", "bbbb", "cccc", "fjdaksljflkdajfkljsalkfjdlkf",
		"jfkldsjflkjdsalkfjlkdsajflkajds",
		"fsdfsdf"}, os.Stdout)
}
