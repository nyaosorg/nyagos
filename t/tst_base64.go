package main

import (
	"encoding/base64"
	"os"
	"os/exec"
)

func main() {
	bstr := base64.StdEncoding.EncodeToString([]byte(`ls "c:\Program Files"`))
	// println(bstr)
	cmd := exec.Command(`..\nyagos.exe`, "-b", bstr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		println(err.Error())
	}
}
