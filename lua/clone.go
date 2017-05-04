package lua

func (L Lua) CloneTo(newL Lua) error {
	L.GetGlobal("_G")
	defer L.Pop(1)
	err := L.ForInDo(-1, func(src Lua) error {
		if !src.IsString(-2) {
			return nil
		}
		key, err := src.ToString(-2)
		if err != nil {
			return err
		}
		//println("KEY=", key)

		// If new instance has already the member, pass it.
		newL.GetGlobal(key)
		defer newL.Pop(1)
		if !newL.IsNil(-1) {
			return nil
		}
		//println("not found and copy")
		val, err := src.ToPushable(-1)
		if err != nil {
			return err
		}
		//println("push to new instance")
		val.Push(newL)
		newL.SetGlobal(key)
		return nil
	})
	return err
}

func (L Lua) Clone() (Lua, error) {
	newL, err := New()
	if err != nil {
		return Lua(0), err
	}
	newL.OpenLibs()
	err = L.CloneTo(newL)
	return newL, err
}
