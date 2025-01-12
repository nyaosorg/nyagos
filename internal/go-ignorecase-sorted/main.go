package ignoreCaseSorted

import (
	"sort"
	"strings"
)

type _Pair[T any] struct {
	Key   string
	Value T
}

type Dictionary[T any] struct {
	maps  map[string]_Pair[T]
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
		d.maps = make(map[string]_Pair[T])
	}
	if _, ok := d.maps[lowerKey]; !ok {
		at := sort.Search(len(d.order), func(i int) bool { return d.order[i] >= lowerKey })
		d.order = append(d.order, "")
		copy(d.order[at+1:], d.order[at:])
		d.order[at] = lowerKey
	}
	d.maps[lowerKey] = _Pair[T]{Key: key, Value: val}
}

func (d *Dictionary[T]) Delete(key string) {
	if d.maps == nil {
		return
	}
	lowerKey := strings.ToLower(key)
	delete(d.maps, lowerKey)
	if len(d.order) > 0 {
		at := sort.Search(len(d.order), func(i int) bool { return d.order[i] >= lowerKey })
		if at < len(d.order) && d.order[at] == lowerKey {
			copy(d.order[at:], d.order[at+1:])
			d.order = d.order[:len(d.order)-1]
		}
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
	var v _Pair[T]
	v, ok = d.maps[lowerKey]
	if ok {
		val = v.Value
	}
	return
}

func (d *Dictionary[T]) Range(f func(string, T) bool) {
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
