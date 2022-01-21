package prism

var (
	ReturnSpecial = MFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Return"},
		OmegaType: AtomicType{AnyType: true},
		Returns:   AtomicType{AnyType: true},
	}
	PrintlnSpecial = MFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Println"},
		OmegaType: AtomicType{AnyType: true},
		Returns:   VoidType,
	}
	PrintSpecial = MFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Print"},
		OmegaType: AtomicType{AnyType: true},
		Returns:   VoidType,
	}

	LenSpecial = MFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Len"},
		OmegaType: VectorType{AnyType: true},
		Returns:   IntType,
	}
	CapSpecial = MFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Cap"},
		OmegaType: VectorType{AnyType: true},
		Returns:   IntType,
	}
	SumSpecial = MFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Sum"},
		OmegaType: VectorType{Type: SomeType{Types: []Type{IntType, RealType}}},
		Returns:   SomeType{Types: []Type{IntType, RealType}},
	}

	ProductSpecial = MFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Product"},
		OmegaType: VectorType{Type: SomeType{Types: []Type{IntType, RealType}}},
		Returns:   SomeType{Types: []Type{IntType, RealType}},
	}
)
