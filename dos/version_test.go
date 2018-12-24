package dos_test

import (
	"github.com/zetamatta/nyagos/dos"
	"testing"
)

func TestVerify(t *testing.T) {
	if dos.Windows10.Verify() {
		println("Windows10")
	} else if dos.Windows81.Verify() {
		println("Windows8.1")
	} else if dos.Windows8.Verify() {
		println("Windows8")
	} else if dos.Windows7.Verify() {
		println("Windows7")
	} else if dos.WindowsVista.Verify() {
		println("WindowsVista")
	} else {
		println("Unknown Windows")
	}
}
