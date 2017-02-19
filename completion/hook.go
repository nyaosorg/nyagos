package completion

import (
	"strings"

	"../interpreter"
	"../readline"
)

var command_listupper = []func() []string{
	func() []string { return listUpAllExecutableOnEnv("PATH") },
	func() []string { return listUpAllExecutableOnEnv("NYAGOSPATH") },
}

func AppendCommandLister(f func() []string) {
	command_listupper = append(command_listupper, f)
}

var HookToList = []func(*readline.Buffer, *CompletionList) (*CompletionList, error){
	luaHook,
}

var PercentFuncs = []func(string) []string{
	listUpOsEnv,
	listUpDynamicEnv,
}

func listUpDynamicEnv(name string) []string {
	matches := []string{}
	for envName, _ := range interpreter.PercentFunc {
		if strings.HasPrefix(envName, name) {
			matches = append(matches, "%"+envName+"%")
		}
	}
	return matches
}
