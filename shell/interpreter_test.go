package shell_test

import (
	"context"
	"os"
	"testing"
	"path/filepath"

	"github.com/zetamatta/nyagos/shell"
)

func TestInterpret(t *testing.T) {
	ctx := context.Background()
	tempFilePath := filepath.Join(os.TempDir(),"hogehoge")
	_, err := shell.New().Interpret(ctx, `cmd /c "echo 12345" > "%TEMP%\hogehoge"`)
	if err != nil {
		t.Fatalf("Fail: %s",err.Error())
	}

	tempFileData,err := os.ReadFile(tempFilePath)
	if err != nil {
		t.Fatalf("Fail: `%s` not found",tempFilePath)
	}
	defer os.Remove(tempFilePath)
	if data := string(tempFileData) ; data != "12345\r\n" {
		t.Fatalf("Fail: %s's contents is expected as \"12345\\r\\n\", but %v",
				tempFilePath, string(tempFileData))
	}
}
