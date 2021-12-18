package prism

import (
	"fmt"

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

func (f Function) String() (s string) {
	s += "Δ " + f.AlphaType.String() + " " + f.Name.String() + " " +
		f.OmegaType.String() + " -> " + f.Returns.String() + "\n"

	if f.Body != nil {
		for _, p := range *f.Body {
			s += " " + p.String() + "\n"
		}
	} else if f.PreBody != nil {
		for _, p := range *f.PreBody {
			s += " " + p.String() + "\n\t"
		}
	}

	return s + "∇\n"
}

func (b Application) String() string {
	return b.Operator.Name.String() + " (" + b.Operand.String() + ")"
}

func (u Dangle) String() string {
	return "(" + u.Outer.String() + " " + u.Inner.String() + ")"
}

func (i Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (e EOF) String() string {
	return "EOF"
}

func (s Subexpression) String() string {
	return "(" + s.Expression.String() + ")"
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

func (a Alpha) String() string {
	return "α"
}

func (o Omega) String() string {
	return "ω"
}
