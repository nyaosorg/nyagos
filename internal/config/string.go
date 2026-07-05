package config

import (
	"github.com/nyaosorg/nyagos/internal/go-ignorecase-sorted"
)

type stringInterface interface {
	Usage() string
	NoUsage() string
	Set(value string)
	Get() string
}

var Strings = ignoreCaseSorted.MapToDictionary(map[string]stringInterface{})

func String(name, defaultv, usage, noUsage string) *string {
	value := &ConfigPtr[string]{
		ptr:     &defaultv,
		usage:   usage,
		noUsage: noUsage,
	}
	Strings.Set(name, value)
	return &defaultv
}
