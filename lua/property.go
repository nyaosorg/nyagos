package lua

type Property struct {
	Pointer *Pushable
}

func (this Property) Push(L Lua) int {
	return (*this.Pointer).Push(L)
}

func (this Property) Set(L Lua, index int) error {
	var err error
	*this.Pointer, err = L.ToPushable(index)
	return err
}

type StringProperty struct {
	Pointer *string
}

func (this StringProperty) Push(L Lua) int {
	L.PushString(*this.Pointer)
	return 1
}

func (this StringProperty) Set(L Lua, index int) error {
	s, err := L.ToString(index)
	if err == nil {
		*this.Pointer = s
	}
	return err
}

type MetaOnlyTableT struct {
	Table TTable
}

func (this MetaOnlyTableT) Push(L Lua) int {
	L.NewTable()
	L.NewTable()
	for key, val := range this.Table.Dict {
		L.Push(val)
		L.SetField(-2, key)
	}
	L.SetMetaTable(-2)
	return 1
}

func NewVirtualTable(getter func(Lua) int, setter func(Lua) int) Pushable {
	return &MetaOnlyTableT{
		TTable{
			Dict: map[string]Pushable{
				"__index":    &TGoFunction{getter},
				"__newindex": &TGoFunction{setter},
				"__call":     &TGoFunction{setter},
			},
			Array: map[int]Pushable{},
		},
	}
}
