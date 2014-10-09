package conio

import "os"
import "os/signal"

var CtrlC = false

func DisableCtrlC() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {
			CtrlC = true
		}
	}()
}
