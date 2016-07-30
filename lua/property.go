package lua

type Property struct {
	Pointer *Pushable
}

func (this Property) Push(L Lua) int {
	return (*this.Pointer).Push(L)
}

func (this Property) Set(L Lua, index int) error {
	value, err := L.ToPushable(index)
	if err == nil {
		*this.Pointer = value
	}
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

type BoolProperty struct {
	Pointer *bool
}

func (this BoolProperty) Push(L Lua) int {
	L.PushBool(*this.Pointer)
	return 1
}

func (this BoolProperty) Set(L Lua, index int) error {
	*this.Pointer = L.ToBool(index)
	return nil
}

type MetaOnlyTableT struct {
	Name  string
	Table TTable
}

func (this MetaOnlyTableT) Push(L Lua) int {
	L.NewUserData(0)
	if this.Name == "" {
		L.NewTable()
	} else {
		L.NewMetaTable(this.Name)
	}
	for key, val := range this.Table.Dict {
		L.Push(val)
		L.SetField(-2, key)
	}
	L.SetMetaTable(-2)
	return 1
}

func NewVirtualTable(name string, getter func(Lua) int, setter func(Lua) int) Pushable {
	return &MetaOnlyTableT{
		Name: name,
		Table: TTable{
			Dict: map[string]Pushable{
				"__index":    &TGoFunction{getter},
				"__newindex": &TGoFunction{setter},
				"__call":     &TGoFunction{setter},
			},
			Array: map[int]Pushable{},
		},
	}
}
