package dos

import "testing"

func TestWNetEnum(t *testing.T) {
	err := WNetEnum(func(local string, remote string) {
		println(local + ":" + remote)
	})
	if err != nil {
		println(err.Error())
	}
}
