package prism

import (
	"fmt"

	"github.com/alecthomas/repr"
	"github.com/llir/llvm/ir/types"
)

func (a AtomicType) String() string {
	return a.Name.String()
}

func (v VectorType) String() string {
	return "[" + v.ElementType.String() + "]"
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

func (a AtomicType) Width() int64 {
	return int64(a.WidthInBytes)
}

func (v VectorType) Width() int64 {
	return v.ElementType.Width() + 16
}

func (s StructType) Width() (acc int64) {
	for _, v := range s.FieldTypes {
		acc += v.Width()
	}

	return acc
}

func (a AtomicType) Realise() types.Type {
	return a.Actual
}

func (v VectorType) Realise() types.Type {
	return types.NewStruct(
		types.I64, types.I64,
		types.NewPointer(v.ElementType.Realise()))
}

func (s StructType) Realise() types.Type {
	acc := []types.Type{}
	for _, v := range s.FieldTypes {
		acc = append(acc, v.Realise())
	}

	return types.NewStruct(acc...)
}

func (f DFunction) String() (s string) {
	s += "Δ " + f.AlphaType.String() + " " + f.Name.String() + " " +
		f.OmegaType.String() + " -> " + f.Returns.String() + "\n"

	if f.Body != nil {
		for _, p := range f.Body {
			s += " " + p.String() + "\n"
		}
	} else if f.PreBody != nil {
		repr.String(f.PreBody)
	}

	return s + "∇\n"
}

func (f MFunction) String() (s string) {
	s += "Δ " + f.Name.String() + " " +
		f.OmegaType.String() + " -> " + f.Returns.String() + "\n"

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
	return d.Operator.String() + " (" + d.Left.String() + ", " + d.Right.String() + ")"
}

func (m MApplication) String() string {
	return m.Operator.String() + " (" + m.Operand.String() + ")"
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
