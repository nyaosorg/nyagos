package nodos_test

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/nyaosorg/go-windows-mbcs"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

// On Windows Terminal once wsl.exe is executed, this test sometimes fails
// because CMD.EXE outputs with UTF-8 encoding even when code page is 932.

func TestTimeFormatOsLayout(t *testing.T) {
	expectBin, err := exec.Command("cmd.exe", "/c", "echo", "%DATE%").Output()
	if err != nil {
		t.Fatalf("exec.Command: %s", err.Error())
	}
	expect, err := mbcs.AnsiToUtf8(expectBin, mbcs.ConsoleCP())
	if err != nil {
		t.Fatalf("mbcs.AnsiToUTf8: %s", err.Error())
	}
	expect = strings.TrimSpace(expect)
	result, err := nodos.TimeFormatOsLayout(time.Now())
	if err != nil {
		t.Fatal(err.Error())
	}
	//println("expect:",expect)
	//println("result:",result)
	if expect != result {
		t.Fatalf("TestOsDateLayout: expect '%s', but '%s'", expect, result)
	}
}
