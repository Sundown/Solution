package prism

import (
	"fmt"

	"github.com/sundown/solution/palisade"
)

type MonadicFunction struct {
	Name      Ident
	OmegaType Type

	Returns Type
	PreBody *[]palisade.Expression
	Body    []Expression

	Attribute Attribute
}

func (MonadicFunction) Arity() int {
	return 1
}

type MonadicApplication struct {
	Operator MonadicFunction
	Operand  Expression
}

// Type property for interface
func (f MonadicFunction) Type() Type {
	return f.Returns
}

// Type property for interface
func (m MonadicApplication) Type() Type {
	return m.Operator.Type()
}

// String function for interface
func (f MonadicFunction) String() (s string) {
	s += f.Name.String() + " " + f.OmegaType.String() + " → " + f.Returns.String() + " {\n"

	if f.Body != nil {
		for _, p := range f.Body {
			s += "  " + p.String() + "\n"
		}
	} else if f.PreBody != nil {
		fmt.Println("PREBODY")
	}

	return s + "}\n"
}

// String function for interface
func (m MonadicApplication) String() string {
	return m.Operator.Name.String() + " " + m.Operand.String()
}

func (m MonadicFunction) Ident() Ident {
	return m.Name
}

func (d MonadicFunction) Attrs() Attribute {
	return d.Attribute
}

func (f MonadicFunction) LLVMise() string {
	return f.Name.Package + "." + f.Name.Name + "_" + f.OmegaType.String() + "." + f.Returns.String()
}
