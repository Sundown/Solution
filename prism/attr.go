package prism

type Attribute struct {
	// Function attributes
	ForceInline        bool
	Special            bool
	SkipBuilder        bool
	DisallowAutoVector bool

	// Expression attributes
	IsConstant bool
}
