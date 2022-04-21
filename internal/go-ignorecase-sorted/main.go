package ignoreCaseSorted

import (
	"sort"
	"strings"
)

type KeyValue[T any] struct {
	Key   string
	Value T
}

type Dictionary[T any] struct {
	maps  map[string]KeyValue[T]
	order []string
}

func (d *Dictionary[T]) Len() int {
	return len(d.maps)
}

func (d *Dictionary[T]) Store(key string, val T) {
	lowerKey := strings.ToLower(key)
	if d.maps == nil {
		d.maps = make(map[string]KeyValue[T])
	}
	d.maps[lowerKey] = KeyValue[T]{Key: key, Value: val}
	if d.order != nil && len(d.order) > 0 {
		d.order = d.order[:0]
	}
}

func (d *Dictionary[T]) Delete(key string) {
	if d.maps == nil {
		return
	}
	lowerKey := strings.ToLower(key)
	delete(d.maps, lowerKey)
	if d.order != nil && len(d.order) > 0 {
		d.order = d.order[:0]
	}
}

func (d *Dictionary[T]) Load(key string) (val T, ok bool) {
	if d.maps == nil {
		return
	}
	lowerKey := strings.ToLower(key)
	var v KeyValue[T]
	v, ok = d.maps[lowerKey]
	if ok {
		val = v.Value
	}
	return
}

func (d *Dictionary[t]) makeOrder() {
	if d.order == nil {
		d.order = make([]string, 0, len(d.maps))
	} else if len(d.order) > 0 {
		return
	}
	for key := range d.maps {
		d.order = append(d.order, key)
	}
	sort.Strings(d.order)
}

func (d *Dictionary[T]) Keys() []string {
	d.makeOrder()
	return d.order
}

func (d *Dictionary[T]) Range(f func(string, T) bool) {
	d.makeOrder()
	for _, lowerKey := range d.order {
		p := d.maps[lowerKey]
		if !f(p.Key, p.Value) {
			break
		}
	}
}

type Enumerator[T any] struct {
	maps  map[string]KeyValue[T]
	order []string
	Key   string
	Value T
}

func (d *Dictionary[T]) Each() *Enumerator[T] {
	d.makeOrder()
	return &Enumerator[T]{
		maps:  d.maps,
		order: d.order,
	}
}

func (e *Enumerator[T]) Range() bool {
	if len(e.order) <= 0 {
		return false
	}
	p := e.maps[e.order[0]]
	e.order = e.order[1:]
	e.Key = p.Key
	e.Value = p.Value
	return true
}

func New[T any](source map[string]T) *Dictionary[T] {
	var d Dictionary[T]
	for key, val := range source {
		d.Store(key, val)
	}
	return &d
}
