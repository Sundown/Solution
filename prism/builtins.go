package prism

var (
	Numeric   = SumType{[]Type{RealType, IntType, CharType, BoolType}}
	Countable = SumType{[]Type{IntType, CharType, BoolType}}
)

func (env Environment) InternBuiltins() {

	env.Intern(MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Return"},
		OmegaType: GenericType{},
		Returns:   GenericType{},
	})

	env.Intern(MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Println"},
		OmegaType: GenericType{},
		Returns:   VoidType,
	})

	env.Intern(MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "Print"},
		OmegaType: GenericType{},
		Returns:   VoidType,
	})

	env.Intern(MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "≢"},
		OmegaType: VectorType{GenericType{}},
		Returns:   IntType,
	})

	env.Intern(MonadicFunction{
		Special:   true,
		Name:      Ident{Package: "_", Name: "__Cap"},
		OmegaType: VectorType{GenericType{}},
		Returns:   IntType,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: true,
		Name:               Ident{Package: "_", Name: "⊃"},
		AlphaType:          IntType,
		OmegaType:          VectorType{GenericType{}},
		Returns:            GenericType{},
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: true,
		Name:               Ident{Package: "_", Name: ","},
		AlphaType:          VectorType{GenericType{}},
		OmegaType:          VectorType{GenericType{}},
		Returns:            VectorType{GenericType{}},
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "="},
		AlphaType:          GenericType{},
		OmegaType:          GenericType{},
		Returns:            BoolType,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "+"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "-"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	})

	env.Intern(MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "-"},
		OmegaType:          Numeric,
		Returns:            Numeric,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "×"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "÷"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            RealType,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "*"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            RealType,
	})

	env.Intern(MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "*"},
		OmegaType:          Numeric,
		Returns:            RealType,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Max"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "∧"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            BoolType,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "∨"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            BoolType,
	})

	env.Intern(MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Max"},
		OmegaType:          Numeric,
		Returns:            Numeric,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Min"},
		AlphaType:          Numeric,
		OmegaType:          Numeric,
		Returns:            Numeric,
	})

	env.Intern(MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "Min"},
		OmegaType:          Numeric,
		Returns:            Numeric,
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "⊢"},
		AlphaType:          GenericType{},
		OmegaType:          GenericType{},
		Returns:            GenericType{},
	})

	env.Intern(MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "⊢"},
		OmegaType:          GenericType{},
		Returns:            GenericType{},
	})

	env.Intern(DyadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "⊣"},
		AlphaType:          GenericType{},
		OmegaType:          GenericType{},
		Returns:            GenericType{},
	})

	env.Intern(MonadicFunction{
		Special:            true,
		disallowAutoVector: true,
		Name:               Ident{Package: "_", Name: "⊂"},
		OmegaType:          GenericType{},
		Returns:            VectorType{Type: GenericType{}},
	})

	env.Intern(MonadicFunction{
		Special:            true,
		disallowAutoVector: true,
		Name:               Ident{Package: "_", Name: "⍳"},
		OmegaType:          Countable,
		Returns:            VectorType{Type: Countable},
	})

	env.Intern(MonadicFunction{
		Special:            true,
		disallowAutoVector: false,
		Name:               Ident{Package: "_", Name: "~"},
		OmegaType:          Countable,
		Returns:            BoolType,
	})
}
