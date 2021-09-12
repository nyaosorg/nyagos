package texts_test

import (
	"testing"

	"github.com/nyaosorg/nyagos/texts"
)

func TestReplaceIgnoreCase(t *testing.T) {
	result := texts.ReplaceIgnoreCase("AHAHAahahaAhaha", "aHaha", "<>")
	if result != "<><><>" {
		t.Fatal("Error: " + result)
		return
	}
}
