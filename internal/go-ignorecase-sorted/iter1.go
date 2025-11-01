package ignoreCaseSorted

type Iterator[T any] struct {
	dic   *Dictionary[T]
	index int
	Key   string
	Value T
}

func (d *Dictionary[T]) Front() *Iterator[T] {
	if len(d.maps) <= 0 {
		return nil
	}
	p := d.maps[d.order[0]]
	return &Iterator[T]{
		dic:   d,
		index: 0,
		Key:   p.Key,
		Value: p.Value,
	}
}

func (iter *Iterator[T]) Next() *Iterator[T] {
	iter.index++
	dic := iter.dic
	if iter.index >= len(dic.order) {
		return nil
	}
	p := dic.maps[dic.order[iter.index]]
	iter.Key = p.Key
	iter.Value = p.Value
	return iter
}

func (d *Dictionary[T]) Back() *Iterator[T] {
	if d.maps == nil || len(d.maps) <= 0 {
		return nil
	}
	index := len(d.order) - 1
	p := d.maps[d.order[index]]
	return &Iterator[T]{
		dic:   d,
		index: index,
		Key:   p.Key,
		Value: p.Value,
	}
}

func (iter *Iterator[T]) Prev() *Iterator[T] {
	iter.index--
	if iter.index < 0 {
		return nil
	}
	dic := iter.dic
	p := dic.maps[dic.order[iter.index]]
	iter.Key = p.Key
	iter.Value = p.Value
	return iter
}
