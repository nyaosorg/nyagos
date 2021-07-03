package commands_test

import (
	"os"
	"strings"
	"testing"

	"github.com/zetamatta/nyagos/commands"
)

func testFixPathCase(t *testing.T, path string) string {
	newpath, err := commands.CorrectCase(path)
	if err != nil {
		t.Fatalf("CorrectCase: %v", err.Error())
	}
	return newpath
}

func TestFixPathCase(t *testing.T) {
	orgPath, err := os.Getwd()
	if err != nil {
		t.Errorf("os.Getwd(): %v", err)
	}
	chgPath := strings.ToUpper(orgPath)
	actPath := testFixPathCase(t, chgPath)
	if actPath != orgPath {
		t.Fatalf("CorrectCase('%s') == %s", chgPath, actPath)
	}
	actPath = testFixPathCase(t, "c:\\")
	if actPath != `C:\` {
		t.Fatalf("CorrectCase('c:\\') == '%s'", actPath)
	}
}
