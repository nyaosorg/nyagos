package functions

import (
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-skk"

	"github.com/nyaosorg/nyagos/internal/onexit"
)

func CmdSkk(args []anyT) []anyT {
	cfg := &skk.Config{
		BindTo: readline.GlobalKeyMap,
	}
	for len(args) > 0 {
		if table, ok := args[0].(map[anyT]anyT); ok {
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
					return []anyT{nil, "key name not found"}
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
		return []any{nil, err.Error()}
	}
	onexit.Register(func() { skkMode.SaveUserJisyo() })
	return []anyT{true}
}
