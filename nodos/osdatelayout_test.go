package nodos_test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/nyaosorg/nyagos/nodos"
)

func testOsDateLayout() error {
	d, err := nodos.OsDateLayout()
	if err != nil {
		return fmt.Errorf("OsDatelayout: %s", err)
	}
	expectBin, err := exec.Command("cmd.exe", "/c", "echo", "%DATE%").Output()
	if err != nil {
		return fmt.Errorf("exec.Command: %w", err)
	}
	expect := strings.TrimSpace(string(expectBin))
	result := time.Now().Format(d)
	//println("expect:",expect)
	//println("result:",result)
	if expect != result {
		return fmt.Errorf("OsDatelayout: differs: '%s' != '%s' (OsDatelayout=%s)", expect, result, d)
	}
	return err
}

func TestOsDatelayout(t *testing.T) {
	if err := testOsDateLayout(); err != nil {
		t.Fatal(err.Error())
	}
}
