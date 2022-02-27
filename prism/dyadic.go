package prism

import (
	"fmt"

	"github.com/sundown/solution/palisade"
)

type DyadicOperator struct {
	Operator int
	Left     Expression
	Right    Expression
	Returns  Type
}

type DyadicFunction struct {
	Special     bool
	SkipBuilder bool
	Inline      bool
	Name        Ident
	AlphaType   Type
	OmegaType   Type
	Returns     Type
	PreBody     *[]palisade.Expression
	Body        []Expression
}

type DyadicApplication struct {
	Operator DyadicFunction
	Left     Expression
	Right    Expression
}

// Type property for interface
func (f DyadicFunction) Type() Type {
	return f.Returns
}

// Type property for interface
func (d DyadicApplication) Type() Type {
	return d.Operator.Type()
}

// String function for interface
func (do DyadicOperator) String() string {
	return do.Left.String() + " " + fmt.Sprint(do.Operator) + " " + do.Right.String()
}

// String function for interface
func (f DyadicFunction) String() (s string) {
	s += "Δ " + f.AlphaType.String() + " " + f.Name.String() + " " +
		f.OmegaType.String() + " → " + f.Returns.String() + "\n"

	if f.Body != nil {
		for _, p := range f.Body {
			s += " " + p.String() + "\n"
		}
	} else if f.PreBody != nil {
		fmt.Println("PREBODY")
	}

	return s + "∇\n"
}

// String function for interface
func (d DyadicApplication) String() string {
	return d.Left.String() + " " + d.Operator.Name.String() + " " + d.Right.String()
}

// Type property for interface
//
// Operators each return a type dependant on a different input
func (do DyadicOperator) Type() Type {
	switch do.Operator {
	case KindMapOperator:
		return VectorType{do.Left.Type()}
	case KindReduceOperator:
		return do.Left.Type()
	}

	Panic("Need to impl type")
	panic(nil)
}

func (d DyadicFunction) IsSpecial() bool {
	return d.Special
}

func (d DyadicFunction) Ident() Ident {
	return d.Name
}

func (f DyadicFunction) LLVMise() string {
	return f.Name.Package + "::" + f.Name.Name + "_" + f.AlphaType.String() + "," + f.OmegaType.String() + "->" + f.Returns.String()
}

func (f DyadicFunction) ShouldInline() bool {
	return f.Inline
}
