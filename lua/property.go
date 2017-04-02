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

type VirtualTable struct {
	Name     string
	Index    func(Lua) int
	NewIndex func(Lua) int
	Call     func(Lua) int
	Len      func(Lua) int
}

func (this *VirtualTable) Push(L Lua) int {
	dict := map[string]Pushable{}
	if this.Index != nil {
		dict["__index"] = TGoFunction(this.Index)
	}
	if this.NewIndex != nil {
		dict["__newindex"] = TGoFunction(this.NewIndex)
	}
	if this.Call != nil {
		dict["__call"] = TGoFunction(this.Call)
	}
	if this.Len != nil {
		dict["__len"] = TGoFunction(this.Len)
	}
	return L.Push(&MetaOnlyTableT{
		Name: this.Name,
		Table: TTable{
			Dict:  dict,
			Array: map[int]Pushable{},
		},
	})
}
