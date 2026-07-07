package config

import (
	"github.com/nyaosorg/nyagos/internal/go-ignorecase-sorted"
)

type configInterface[T any] interface {
	Usage() string
	NoUsage() string
	Set(value T)
	Get() T
}

type configVar[T any] struct {
	ptr     *T
	usage   string
	noUsage string
}

func (C *configVar[T]) Usage() string   { return C.usage }
func (C *configVar[T]) NoUsage() string { return C.noUsage }
func (C *configVar[T]) Set(value T)     { *C.ptr = value }
func (C *configVar[T]) Get() T          { return *C.ptr }

type configFunc[T any] struct {
	Setter  func(value T)
	Getter  func() T
	usage   string
	noUsage string
}

func (C *configFunc[T]) Usage() string   { return C.usage }
func (C *configFunc[T]) NoUsage() string { return C.noUsage }
func (C *configFunc[T]) Set(value T)     { C.Setter(value) }
func (C *configFunc[T]) Get() T          { return C.Getter() }

var Bools ignoreCaseSorted.Dictionary[configInterface[bool]]

func BoolVar(p *bool, name, usage, noUsage string) {
	value := &configVar[bool]{
		ptr:     p,
		usage:   usage,
		noUsage: noUsage,
	}
	Bools.Set(name, value)
}

func Bool(name string, defaultv bool, usage, noUsage string) *bool {
	BoolVar(&defaultv, name, usage, noUsage)
	return &defaultv
}

var Strings ignoreCaseSorted.Dictionary[configInterface[string]]

func String(name, defaultv, usage, noUsage string) *string {
	value := &configVar[string]{
		ptr:     &defaultv,
		usage:   usage,
		noUsage: noUsage,
	}
	Strings.Set(name, value)
	return &defaultv
}
