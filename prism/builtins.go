package prism

var (
	Numeric   = TypeGroup{[]Type{RealType, IntType, CharType, BoolType}}
	Countable = TypeGroup{[]Type{IntType, CharType, BoolType}}
)

func (env Environment) InternBuiltins() {
	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: true},
		Name:      Ident{Package: "_", Name: "←"},
		OmegaType: GenericType{},
		Returns:   GenericType{},
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "Println"},
		OmegaType: GenericType{},
		Returns:   VoidType,
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "Print"},
		OmegaType: GenericType{},
		Returns:   VoidType,
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: true},
		Name:      Ident{Package: "_", Name: "≢"},
		OmegaType: VectorType{GenericType{}},
		Returns:   IntType,
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "__Cap"},
		OmegaType: VectorType{GenericType{}},
		Returns:   IntType,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: true},
		Name:      Ident{Package: "_", Name: "⊃"},
		AlphaType: IntType,
		OmegaType: VectorType{GenericType{}},
		Returns:   GenericType{},
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: true},
		Name:      Ident{Package: "_", Name: ","},
		AlphaType: VectorType{GenericType{}},
		OmegaType: VectorType{GenericType{}},
		Returns:   VectorType{GenericType{}},
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "="},
		AlphaType: GenericType{},
		OmegaType: GenericType{},
		Returns:   BoolType,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "+"},
		AlphaType: Numeric,
		OmegaType: Numeric,
		Returns:   Numeric,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "-"},
		AlphaType: Numeric,
		OmegaType: Numeric,
		Returns:   Numeric,
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "-"},
		OmegaType: Numeric,
		Returns:   Numeric,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "×"},
		AlphaType: Numeric,
		OmegaType: Numeric,
		Returns:   Numeric,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "÷"},
		AlphaType: Numeric,
		OmegaType: Numeric,
		Returns:   RealType,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "*"},
		AlphaType: Numeric,
		OmegaType: Numeric,
		Returns:   RealType,
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "*"},
		OmegaType: Numeric,
		Returns:   RealType,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "⌈"},
		AlphaType: Numeric,
		OmegaType: Numeric,
		Returns:   Numeric,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "∧"},
		AlphaType: Numeric,
		OmegaType: Numeric,
		Returns:   BoolType,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "∨"},
		AlphaType: Numeric,
		OmegaType: Numeric,
		Returns:   BoolType,
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "⌈"},
		OmegaType: Numeric,
		Returns:   Numeric,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "⌊"},
		AlphaType: Numeric,
		OmegaType: Numeric,
		Returns:   Numeric,
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "⌊"},
		OmegaType: Numeric,
		Returns:   Numeric,
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "⊢"},
		AlphaType: GenericType{},
		OmegaType: GenericType{},
		Returns:   GenericType{},
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "⊢"},
		OmegaType: GenericType{},
		Returns:   GenericType{},
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "⊣"},
		AlphaType: GenericType{},
		OmegaType: GenericType{},
		Returns:   GenericType{},
	})

	env.Intern(DyadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "⍴"},
		AlphaType: GenericType{},
		OmegaType: GenericType{},
		Returns:   GenericType{},
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: true},
		Name:      Ident{Package: "_", Name: "⊂"},
		OmegaType: GenericType{},
		Returns:   VectorType{Type: GenericType{}},
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: true},
		Name:      Ident{Package: "_", Name: "⍳"},
		OmegaType: Countable,
		Returns:   VectorType{Type: Countable},
	})

	env.Intern(MonadicFunction{
		Attribute: Attribute{Special: true, DisallowAutoVector: false},
		Name:      Ident{Package: "_", Name: "~"},
		OmegaType: Countable,
		Returns:   BoolType,
	})
}
