package texts_test

import (
	"testing"

	"github.com/zetamatta/nyagos/texts"
)

func TestReplaceIgnoreCase(t *testing.T) {
	result := texts.ReplaceIgnoreCase("AHAHAahahaAhaha", "aHaha", "<>")
	if result != "<><><>" {
		t.Fatal("Error: " + result)
		return
	}
}
