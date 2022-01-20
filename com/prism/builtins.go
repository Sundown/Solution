package prism

var (
	ReturnSpecial = MFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Return"},
		OmegaType: AtomicType{AnyType: true},
		Returns:   AtomicType{AnyType: true},
	}
)
