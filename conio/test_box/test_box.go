package main

import (
	"fmt"

	conio ".."
	"github.com/mattn/go-colorable"
)

var Console = colorable.NewColorableStderr()

func main() {
	result := conio.BoxChoice([]string{
		"ahaha1", "ihihi1", "ufufuf1", "ehehe1", "ohoho1", "uhaha1",
		"ahaha2", "ihihi2", "ufufuf2", "ehehe2", "ohoho2", "uhaha2",
		"ahaha3", "ihihi3", "ufufuf3", "ehehe3", "ohoho3", "uhaha3",
		"ahaha4", "ihihi4", "ufufuf4", "ehehe4", "ohoho4", "uhaha4",
		"ahaha5", "ihihi5", "ufufuf5", "ehehe5", "ohoho5", "uhaha5",
		"ahaha6", "ihihi6", "ufufuf6", "ehehe6", "ohoho6", "uhaha6",
		"ahaha7", "ihihi7", "ufufuf7", "ehehe7", "ohoho7", "uhaha7",
	}, Console)
	fmt.Printf("\n--> %s\n", result)
}
