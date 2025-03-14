package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/nyaosorg/go-readline-ny"

	"github.com/nyaosorg/nyagos/internal/completion"
	"github.com/nyaosorg/nyagos/internal/go-ignorecase-sorted"
	"github.com/nyaosorg/nyagos/internal/shell"
)

// ReadStdinAsFile is the flat to read commands from stdin as a file stream
var ReadStdinAsFile = false

type Bool interface {
	Usage() string
	NoUsage() string
	Set(value bool)
	Get() bool
}

type BoolPtr = ConfigPtr[bool]
type BoolFunc = ConfigFunc[bool]

var OptionPredictColor = true

// Bools are the all global option list.
var Bools = ignoreCaseSorted.MapToDictionary(map[string]Bool{
	"completion_hidden": &BoolPtr{
		ptr:     &completion.IncludeHidden,
		usage:   "Include hidden files on completion",
		noUsage: "Do not include hidden files on completion",
	},
	"completion_slash": &BoolPtr{
		ptr:     &completion.UseSlash,
		usage:   "use forward slash on completion",
		noUsage: "Do not use slash on completion",
	},
	"glob": &BoolPtr{
		ptr:     &shell.WildCardExpansionAlways,
		usage:   "Enable to expand wildcards",
		noUsage: "Disable to expand wildcards",
	},
	"glob_slash": &BoolPtr{
		ptr:     &shell.GlobUseSlash,
		usage:   "Use forward slash on wildcard expansion",
		noUsage: "Do not Use forward slash on wildcard expansion",
	},
	"noclobber": &BoolPtr{
		ptr:     &shell.NoClobber,
		usage:   "forbide to overwrite files on redirect",
		noUsage: "Do not forbide to overwrite files no redirect",
	},
	"usesource": &BoolPtr{
		ptr:     &shell.UseSourceRunBatch,
		usage:   "allow batchfile to change environment variables of nyagos",
		noUsage: "forbide batchfile to change environment variables of nyagos",
	},
	"tilde_expansion": &BoolPtr{
		ptr:     &shell.TildeExpansion,
		usage:   "Enable Tilde Expansion",
		noUsage: "Disable Tilde Expansion",
	},
	"read_stdin_as_file": &BoolPtr{
		ptr:     &ReadStdinAsFile,
		usage:   "Read commands from stdin as a file stream. Disable to edit line",
		noUsage: "Read commands from stdin as Windows Console(tty). Enable to edit line",
	},
	"output_surrogate_pair": &BoolFunc{
		Setter:  readline.EnableSurrogatePair,
		Getter:  readline.IsSurrogatePairEnabled,
		usage:   "Output surrogate pair characters as it is",
		noUsage: "Output surrogate pair characters like <NNNNN>",
	},
	"predict": &BoolPtr{
		ptr:     &OptionPredictColor,
		usage:   "Enable prediction on readline",
		noUsage: "Disable prediction on readline",
	},
})

func toLuaLiteral(s string) string {
	var buf strings.Builder
	for _, c := range s {
		if c < ' ' {
			fmt.Fprintf(&buf, "\\%03d", c)

		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

func Dump(w io.Writer) {
	for p := Bools.Front(); p != nil; p = p.Next() {
		v := p.Value.Get()
		if v {
			fmt.Fprintf(w, "-- %s\n", p.Value.Usage())
		} else {
			fmt.Fprintf(w, "-- %s\n", p.Value.NoUsage())
		}
		fmt.Fprintf(w, "nyagos.option.%s=%v\n", p.Key, v)
	}
	for p := Strings.Front(); p != nil; p = p.Next() {
		fmt.Fprintf(w, "-- %s\n", p.Value.Usage())
		fmt.Fprintf(w, "nyagos.option.%s='%s'\n", p.Key, toLuaLiteral(p.Value.Get()))
	}
}
