package prism

var (
	ReturnSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Return"},
		OmegaType: AtomicType{AnyType: true},
		Returns:   AtomicType{AnyType: true},
	}
	PrintlnSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Println"},
		OmegaType: AtomicType{AnyType: true},
		Returns:   VoidType,
	}
	PrintSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Print"},
		OmegaType: AtomicType{AnyType: true},
		Returns:   VoidType,
	}

	LenSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Len"},
		OmegaType: VectorType{AnyType: true},
		Returns:   IntType,
	}
	CapSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Cap"},
		OmegaType: VectorType{AnyType: true},
		Returns:   IntType,
	}

	GEPSpecial = DyadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "GEP"},
		AlphaType: IntType,
		OmegaType: VectorType{AnyType: true},
		Returns:   IntType,
	}
	SumSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Sum"},
		OmegaType: VectorType{Type: TypeGroup{Types: []Type{IntType, RealType}}},
		Returns:   TypeGroup{Types: []Type{IntType, RealType}},
	}

	ProductSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Product"},
		OmegaType: VectorType{Type: TypeGroup{Types: []Type{IntType, RealType}}},
		Returns:   TypeGroup{Types: []Type{IntType, RealType}},
	}
)
