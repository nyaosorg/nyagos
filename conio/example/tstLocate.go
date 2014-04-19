package main

import ".."
import "fmt"

func main() {
	x, y := conio.GetLocate()
	fmt.Printf("X=%d Y=%d\n", x, y)
	conio.Locate(0, 0)
}
