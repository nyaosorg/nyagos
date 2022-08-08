package ignoreCaseSorted

type AscIter[T any] struct {
	dic   *Dictionary[T]
	index int
	Key   string
	Value T
}

func (d *Dictionary[T]) Ascend() *AscIter[T] {
	return &AscIter[T]{
		dic:   d,
		index: -1,
	}
}

func (iter *AscIter[T]) Range() bool {
	iter.index++
	if iter.index >= len(iter.dic.order) {
		return false
	}
	pair := iter.dic.maps[iter.dic.order[iter.index]]
	iter.Key = pair.Key
	iter.Value = pair.Value
	return true
}

type DescIter[T any] struct {
	dic   *Dictionary[T]
	index int
	Key   string
	Value T
}

func (d *Dictionary[T]) Descend() *DescIter[T] {
	return &DescIter[T]{
		dic:   d,
		index: len(d.order),
	}
}

func (iter *DescIter[T]) Range() bool {
	iter.index--
	if iter.index < 0 {
		return false
	}
	pair := iter.dic.maps[iter.dic.order[iter.index]]
	iter.Key = pair.Key
	iter.Value = pair.Value
	return true
}
