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

// Store is same as Set that is a compatible method with "sync".Map
func (d *Dictionary[T]) Store(key string, val T) {
	d.Set(key, val)
}

func (d *Dictionary[T]) Set(key string, val T) {
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

// Load is same as Get that is a compatible method with "sync".Map
func (d *Dictionary[T]) Load(key string) (val T, ok bool) {
	return d.Get(key)
}

func (d *Dictionary[T]) Get(key string) (val T, ok bool) {
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

func MapToDictionary[T any](source map[string]T) *Dictionary[T] {
	var d Dictionary[T]
	for key, val := range source {
		d.Set(key, val)
	}
	return &d
}

type Ascending[T any] struct {
	Key   string
	Value T
	maps  map[string]KeyValue[T]
	order []string
}

func (d *Dictionary[T]) Ascend() *Ascending[T] {
	d.makeOrder()
	if d.maps == nil || len(d.maps) <= 0 {
		return nil
	}
	p := d.maps[d.order[0]]
	return &Ascending[T]{
		maps:  d.maps,
		order: d.order[1:],
		Key:   p.Key,
		Value: p.Value,
	}
}

func (a *Ascending[T]) Next() *Ascending[T] {
	if a == nil || a.order == nil || len(a.order) < 1 {
		return nil
	}
	p := a.maps[a.order[0]]
	a.Key = p.Key
	a.Value = p.Value
	a.order = a.order[1:]
	return a
}
