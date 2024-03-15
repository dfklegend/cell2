package cmd

type FuncCmd struct {
	name string
	cb   func(ctx IContext, args []string, cbFin func(result string))
}

func NewFuncCmd(name string, cb func(ctx IContext, args []string, cbFin func(result string))) ICmd {
	return &FuncCmd{
		name: name,
		cb:   cb,
	}
}

func (f *FuncCmd) GetName() string {
	return f.name
}

func (f *FuncCmd) Do(ctx IContext, args []string, cbFin func(result string)) {
	f.cb(ctx, args, cbFin)
}
