package dos

import (
	"fmt"
	"testing"
)

func TestYenYen2Yen(t *testing.T){
	org := "C:\\WINDOWS\\System32\\WindowsPowerShell\\v1.0\\\\powershell.exe"
	neo := "C:\\WINDOWS\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"

	out := YenYen2Yen(org)
	if out != neo {
		t.Fail()
	}
	fmt.Printf("[%s] -> [%s]\n",org,out)
}
