package types

func (f *Flags) Set(new Flags) {
	*f |= new
}

func (f *Flags) Has(test Flags) bool {
	return *f&test != 0
}

func (f *Flags) Clear(clear Flags) {
	*f &= ^clear
}
