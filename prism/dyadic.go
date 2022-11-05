package prism

import (
	"fmt"

	"github.com/sundown/solution/palisade"
)

type DyadicFunction struct {
	Name      Ident
	AlphaType Type
	OmegaType Type
	Returns   Type

	PreBody *[]palisade.Expression
	Body    []Expression

	Attribute Attribute
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
func (f DyadicFunction) String() (s string) {
	s += f.AlphaType.String() + " " + f.Name.String() + " " +
		f.OmegaType.String() + " -> " + f.Returns.String() + " {\n"

	if f.Body != nil {
		for _, p := range f.Body {
			s += " " + p.String() + "\n"
		}
	} else if f.PreBody != nil {
		fmt.Println("PREBODY")
	}

	return s + "}\n"
}

// String function for interface
func (d DyadicApplication) String() string {
	return d.Left.String() + " " + d.Operator.Name.String() + " " + d.Right.String()
}

func (d DyadicFunction) Attrs() Attribute {
	return d.Attribute
}

func (d DyadicFunction) Ident() Ident {
	return d.Name
}

func (f DyadicFunction) LLVMise() string {
	return f.Name.Package + "." + f.Name.Name + "_" + f.AlphaType.String() + "." + f.OmegaType.String() + "." + f.Returns.String()
}
