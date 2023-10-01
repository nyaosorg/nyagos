package functions

import (
	"fmt"
	"os"
	"strings"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-skk"

	"github.com/nyaosorg/nyagos/internal/onexit"
)

func CmdSkk(args []any) []any {
	cfg := &skk.Config{
		BindTo: readline.GlobalKeyMap,
	}
	var exportVar string
	for len(args) > 0 {
		if table, ok := args[0].(map[any]any); ok {
			if value, okok := table["export"].(string); okok {
				exportVar = value
			}
			if value, okok := table["user"].(string); okok {
				cfg.UserJisyoPath = value
				// println("user:", value)
			}
			if value, okok := table["ctrlj"].(string); okok {
				value := keys.NormalizeName(value)
				if code, ok := keys.NameToCode[value]; ok {
					// println("ctrlj:", value)
					cfg.CtrlJ = code
				} else {
					return []any{nil, "key name not found"}
				}
			}
			i := 0
			for {
				i++
				value, okok := table[i].(string)
				if !okok {
					break
				}
				// println("jisyo:", value)
				cfg.SystemJisyoPaths = append(cfg.SystemJisyoPaths, value)
			}
		} else if value, ok := args[0].(string); ok {
			// println("jisyo:", value)
			cfg.SystemJisyoPaths = append(cfg.SystemJisyoPaths, value)
		}
		args = args[1:]
	}
	skkMode, err := cfg.Setup()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return []any{nil, err.Error()}
	}
	onexit.Register(func() { skkMode.SaveUserJisyo() })
	if exportVar != "" {
		var buffer strings.Builder
		for _, value := range cfg.SystemJisyoPaths {
			if buffer.Len() > 0 {
				buffer.WriteByte(os.PathListSeparator)
			}
			buffer.WriteString(value)
		}
		if cfg.UserJisyoPath != "" {
			if buffer.Len() > 0 {
				buffer.WriteByte(os.PathListSeparator)
			}
			buffer.WriteString("user=")
			buffer.WriteString(cfg.UserJisyoPath)
		}
		os.Setenv(exportVar, buffer.String())
	}
	return []any{true}
}
