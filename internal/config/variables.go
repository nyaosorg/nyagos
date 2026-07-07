package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/nyaosorg/go-readline-ny"

	"github.com/nyaosorg/nyagos/internal/completion"
	"github.com/nyaosorg/nyagos/internal/shell"
)

func init() {
	BoolVar(&completion.IncludeHidden,
		"completion_hidden",
		"Include hidden files on completion",
		"Do not include hidden files on completion")

	BoolVar(&completion.UseSlash,
		"completion_slash",
		"use forward slash on completion",
		"Do not use slash on completion")

	BoolVar(&shell.WildCardExpansionAlways,
		"glob",
		"Enable to expand wildcards",
		"Disable to expand wildcards")

	BoolVar(&shell.GlobUseSlash,
		"glob_slash",
		"Use forward slash on wildcard expansion",
		"Do not Use forward slash on wildcard expansion")

	BoolVar(&shell.NoClobber,
		"noclobber",
		"forbide to overwrite files on redirect",
		"Do not forbide to overwrite files no redirect")

	BoolVar(&shell.UseSourceRunBatch,
		"usesource",
		"allow batchfile to change environment variables of nyagos",
		"forbide batchfile to change environment variables of nyagos")

	BoolVar(&shell.TildeExpansion,
		"tilde_expansion",
		"Enable Tilde Expansion",
		"Disable Tilde Expansion")

	Bools.Set("output_surrogate_pair", &configFunc[bool]{
		Setter:  readline.EnableSurrogatePair,
		Getter:  readline.IsSurrogatePairEnabled,
		usage:   "Output surrogate pair characters as it is",
		noUsage: "Output surrogate pair characters like <NNNNN>",
	})
}

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
