package prism

var (
	ReturnSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Return"},
		OmegaType: GenericType{},
		Returns:   GenericType{},
	}

	PrintlnSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Println"},
		OmegaType: GenericType{},
		Returns:   VoidType,
	}

	PrintSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Print"},
		OmegaType: GenericType{},
		Returns:   VoidType,
	}

	LenSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Len"},
		OmegaType: VectorType{GenericType{}},
		Returns:   IntType,
	}

	CapSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Cap"},
		OmegaType: VectorType{GenericType{}},
		Returns:   IntType,
	}

	GEPSpecial = DyadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "GEP"},
		AlphaType: IntType,
		OmegaType: VectorType{GenericType{}},
		Returns:   GenericType{},
	}

	AppendSpecial = DyadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Append"},
		AlphaType: VectorType{GenericType{}},
		OmegaType: VectorType{GenericType{}},
		Returns:   VectorType{GenericType{}},
	}

	EqSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "="},
		AlphaType:          GenericType{},
		OmegaType:          GenericType{},
		Returns:            BoolType,
	}

	AddSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "+"},
		AlphaType:          SumType{[]Type{IntType, RealType}},
		OmegaType:          SumType{[]Type{IntType, RealType}},
		Returns:            SumType{[]Type{IntType, RealType}},
	}

	SubSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "-"},
		AlphaType:          SumType{[]Type{IntType, RealType}},
		OmegaType:          SumType{[]Type{IntType, RealType}},
		Returns:            SumType{[]Type{IntType, RealType}},
	}

	MulSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "*"},
		AlphaType:          SumType{[]Type{IntType, RealType}},
		OmegaType:          SumType{[]Type{IntType, RealType}},
		Returns:            SumType{[]Type{IntType, RealType}},
	}

	DivSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "÷"},
		AlphaType:          SumType{[]Type{IntType, RealType}},
		OmegaType:          SumType{[]Type{IntType, RealType}},
		Returns:            RealType,
	}

	MaxSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Max"},
		AlphaType:          SumType{[]Type{IntType, RealType}},
		OmegaType:          SumType{[]Type{IntType, RealType}},
		Returns:            SumType{[]Type{IntType, RealType}},
	}

	AndSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "&"},
		AlphaType:          SumType{[]Type{IntType, RealType, BoolType}},
		OmegaType:          SumType{[]Type{IntType, RealType, BoolType}},
		Returns:            SumType{[]Type{IntType, RealType, BoolType}},
	}

	OrSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "|"},
		AlphaType:          SumType{[]Type{IntType, RealType, BoolType}},
		OmegaType:          SumType{[]Type{IntType, RealType, BoolType}},
		Returns:            SumType{[]Type{IntType, RealType, BoolType}},
	}

	CeilSpecial = MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Max"},
		OmegaType:          SumType{[]Type{IntType, RealType}},
		Returns:            SumType{[]Type{IntType, RealType}},
	}

	MinSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Min"},
		AlphaType:          SumType{[]Type{IntType, RealType}},
		OmegaType:          SumType{[]Type{IntType, RealType}},
		Returns:            SumType{[]Type{IntType, RealType}},
	}

	FloorSpecial = MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Min"},
		OmegaType:          SumType{[]Type{IntType, RealType}},
		Returns:            SumType{[]Type{IntType, RealType}},
	}
)
