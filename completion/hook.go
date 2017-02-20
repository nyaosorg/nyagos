package completion

import (
	"strings"

	"../interpreter"
	"../readline"
)

var command_listupper = []func() []Element{
	func() []Element { return listUpAllExecutableOnEnv("PATH") },
	func() []Element { return listUpAllExecutableOnEnv("NYAGOSPATH") },
}

func AppendCommandLister(f func() []Element) {
	command_listupper = append(command_listupper, f)
}

var HookToList = []func(*readline.Buffer, *CompletionList) (*CompletionList, error){
	luaHook,
}

var PercentFuncs = []func(string) []Element{
	listUpOsEnv,
	listUpDynamicEnv,
}

func listUpDynamicEnv(name string) []Element {
	matches := []Element{}
	for envName, _ := range interpreter.PercentFunc {
		if strings.HasPrefix(envName, name) {
			value := "%" + envName + "%"
			matches = append(matches, Element{value, value})
		}
	}
	return matches
}
