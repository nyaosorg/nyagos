package texts

import (
	"testing"
)

func TestReplaceIgnoreCase(t *testing.T) {
	result := ReplaceIgnoreCase("AHAHAahahaAhaha", "aHaha", "<>")
	if result != "<><><>" {
		t.Fatal("Error: " + result)
		return
	}
}
