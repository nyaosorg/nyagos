package nodos

import (
	"fmt"
	"time"
)

var clockmark = []rune{'/', '-', '\u2216', '|'}

func Progress() func() {
	done := make(chan struct{})
	fmt.Print(" ")
	go func() {
		ticker := time.NewTicker(time.Second / 2)
		i := 0
		for {
			select {
			case <-done:
				ticker.Stop()
				close(done)
				fmt.Print(" \b\b")
				return
			case <-ticker.C:
				if i >= len(clockmark) {
					i = 0
				}
				fmt.Printf("%c\b", clockmark[i])
				i++
			}
		}
	}()

	return func() {
		done <- struct{}{}
	}
}
