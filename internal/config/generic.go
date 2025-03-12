package config

type ConfigPtr[T any] struct {
	ptr     *T
	usage   string
	noUsage string
}

func (C *ConfigPtr[T]) Usage() string   { return C.usage }
func (C *ConfigPtr[T]) NoUsage() string { return C.noUsage }
func (C *ConfigPtr[T]) Set(value T)     { *C.ptr = value }
func (C *ConfigPtr[T]) Get() T          { return *C.ptr }

type ConfigFunc[T any] struct {
	Setter  func(value T)
	Getter  func() T
	usage   string
	noUsage string
}

func (C *ConfigFunc[T]) Usage() string   { return C.usage }
func (C *ConfigFunc[T]) NoUsage() string { return C.noUsage }
func (C *ConfigFunc[T]) Set(value T)     { C.Setter(value) }
func (C *ConfigFunc[T]) Get() T          { return C.Getter() }
