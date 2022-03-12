package prism

var (
	Numeric = SumType{[]Type{IntType, RealType, CharType, BoolType}}

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
		Special:            true,
		disallowAutoVector: true,
		Name:               Ident{Package: "_", Name: ","},
		AlphaType:          VectorType{GenericType{}},
		OmegaType:          VectorType{GenericType{}},
		Returns:            VectorType{GenericType{}},
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
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	}

	SubSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "-"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	}

	MulSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "*"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	}

	DivSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "÷"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            RealType,
	}

	MaxSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Max"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	}

	AndSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "&"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            BoolType,
	}

	OrSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "|"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            BoolType,
	}

	CeilSpecial = MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Max"},
		OmegaType:          Numeric,
		Returns:            Numeric,
	}

	MinSpecial = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Min"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	}

	FloorSpecial = MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Min"},
		OmegaType:          Numeric,
		Returns:            Numeric,
	}

	RightHook = DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "⊢"},
		AlphaType:          SumType{[]Type{RealType, IntType}},
		OmegaType:          SumType{[]Type{RealType, IntType}},
		Returns:            SumType{[]Type{RealType, IntType}},
	}
)
