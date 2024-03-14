package cmd

type FuncCmd struct {
	name string
	cb   func(args []string)
}

func NewFuncCmd(name string, cb func(args []string)) ICmd {
	return &FuncCmd{
		name: name,
		cb:   cb,
	}
}

func (f *FuncCmd) GetName() string {
	return f.name
}

func (f *FuncCmd) Do(args []string) {
	f.cb(args)
}
