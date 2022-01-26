package prism

var (
	ReturnSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Return"},
		OmegaType: SemiDeterminedType{},
		Returns:   SemiDeterminedType{},
	}
	PrintlnSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Println"},
		OmegaType: SemiDeterminedType{},
		Returns:   VoidType,
	}
	PrintSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Print"},
		OmegaType: SemiDeterminedType{},
		Returns:   VoidType,
	}

	LenSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Len"},
		OmegaType: VectorType{SemiDeterminedType{}},
		Returns:   IntType,
	}
	CapSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Cap"},
		OmegaType: VectorType{SemiDeterminedType{}},
		Returns:   IntType,
	}

	GEPSpecial = DyadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "GEP"},
		AlphaType: IntType,
		OmegaType: VectorType{SemiDeterminedType{}},
		Returns:   SemiDeterminedType{},
	}

	AppendSpecial = DyadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Append"},
		AlphaType: VectorType{SemiDeterminedType{}},
		OmegaType: VectorType{SemiDeterminedType{}},
		Returns:   VectorType{SemiDeterminedType{}},
	}

	SumSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Sum"},
		OmegaType: VectorType{SemiDeterminedTypeGroup{[]Type{IntType, RealType}}},
		Returns:   SemiDeterminedTypeGroup{[]Type{IntType, RealType}},
	}

	AddSpecial = DyadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Add"},
		AlphaType: SemiDeterminedTypeGroup{[]Type{IntType, RealType}},
		OmegaType: SemiDeterminedTypeGroup{[]Type{IntType, RealType}},
		Returns:   SemiDeterminedTypeGroup{[]Type{IntType, RealType}},
	}

	ProductSpecial = MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Product"},
		OmegaType: VectorType{SemiDeterminedTypeGroup{[]Type{IntType, RealType}}},
		Returns:   SemiDeterminedTypeGroup{[]Type{IntType, RealType}},
	}
)
