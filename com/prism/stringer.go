package prism

import (
	"fmt"

	"github.com/alecthomas/repr"
)

func (a AtomicType) String() string {
	return a.Name.String()
}

func (v VectorType) String() string {
	return "[" + v.Type.String() + "]"
}

func (s StructType) String() (acc string) {
	acc = "("
	for i, v := range s.FieldTypes {
		if i > 0 {
			acc += ", "
		}

		acc += v.String()
	}

	return acc + ")"
}

func (i Ident) String() string {
	if i.Package == "_" {
		return i.Name
	}

	return i.Package + "::" + i.Name
}

func (Void) String() string {
	return "Void"
}

func (c Cast) String() string {
	return "<" + c.ToType.String() + ">" + c.Value.String()
}

func (do DyadicOperator) String() string {
	return do.Left.String() + " " + fmt.Sprint(do.Operator) + " " + do.Right.String()
}

func (f DyadicFunction) String() (s string) {
	s += "Δ " + f.AlphaType.String() + " " + f.Name.String() + " " +
		f.OmegaType.String() + " → " + f.Returns.String() + "\n"

	if f.Body != nil {
		for _, p := range f.Body {
			s += " " + p.String() + "\n"
		}
	} else if f.PreBody != nil {
		repr.String(f.PreBody)
	}

	return s + "∇\n"
}

func (f MonadicFunction) String() (s string) {
	s += "Δ " + f.Name.String() + " " +
		f.OmegaType.String() + " → " + f.Returns.String() + "\n"

	if f.Body != nil {
		for _, p := range f.Body {
			s += " " + p.String() + "\n"
		}
	} else if f.PreBody != nil {
		repr.String(f.PreBody)
	}

	return s + "∇\n"
}

func (d DApplication) String() string {
	return d.Left.String() + " " + d.Operator.Name.String() + " " + d.Right.String()
}

func (m MApplication) String() string {
	return m.Operator.Name.String() + " " + m.Operand.String()
}

func (i Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (e EOF) String() string {
	return "EOF"
}

func (r Real) String() string {
	return fmt.Sprintf("%f", r.Value)
}

func (s String) String() string {
	return "\"" + s.Value + "\""
}

func (c Char) String() string {
	return "'" + string(c.Value) + "'"
}

func (b Bool) String() string {
	if b.Value {
		return "True"
	}

	return "False"
}

func (v Vector) String() string {
	var s string
	for _, v := range *v.Body {
		s += v.String() + " "
	}

	return s
}

func (a Alpha) String() string {
	return "α"
}

func (o Omega) String() string {
	return "ω"
}

func (e Environment) String() (s string) {
	for _, f := range e.DyadicFunctions {
		s += f.String()
	}
	for _, f := range e.MonadicFunctions {
		s += f.String()
	}

	return
}

func (s GenericType) String() string {
	return "T"
}

func (s SumType) String() (res string) {
	for i, t := range s.Types {
		if i > 0 {
			res += " | "
		}
		res += t.String()
	}

	return
}
