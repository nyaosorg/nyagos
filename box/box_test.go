package box

import "os"
import "testing"

func TestPrint(t *testing.T) {
	Print([]string{
		"aaaa", "bbbb", "cccc", "fjdaksljflkdajfkljsalkfjdlkf",
		"jfkldsjflkjdsalkfjlkdsajflkajds",
		"fsdfsdf"}, 80, os.Stdout)
}
