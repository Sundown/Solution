package prism

import (
	"sundown/solution/palisade"

	"github.com/llir/llvm/ir/types"
)

type Environment struct {
	MFunctions map[Ident]*MFunction
	DFunctions map[Ident]*DFunction
	Types      map[Ident]Type
}

type Ident struct {
	Package string
	Name    string
}

func (f DFunction) LLVMise() string {
	return f.Name.Package + "::" + f.Name.Name + " " + f.AlphaType.String() + ", " + f.OmegaType.String() + "->" + f.Returns.String()
}

func (f MFunction) LLVMise() string {
	return f.Name.Package + "::" + f.Name.Name + " " + f.OmegaType.String() + "->" + f.Returns.String()
}

const (
	TypeKindAtomic = iota
	TypeKindVector
	TypeKindStruct
	KindFunction
	TypeInt
	TypeReal
	TypeChar
	TypeBool
	TypeVoid
	TypeString
)

type Type interface {
	Kind() int
	Width() int64
	String() string
	Realise() types.Type
}

type AtomicType struct {
	ID           int
	WidthInBytes int
	Name         Ident
	Actual       types.Type
}

type Vector struct {
	ElementType VectorType
	Body        *[]Expression
}

type VectorType struct {
	Type
}

type StructType struct {
	FieldTypes []Type
}

type Expression interface {
	Type() Type
	String() string
}

type Morpheme interface {
	_atomicflag()
}

type DFunction struct {
	Name      Ident
	AlphaType Type
	OmegaType Type
	Returns   Type
	PreBody   *[]palisade.Expression
	Body      []Expression
}

type MFunction struct {
	Name      Ident
	OmegaType Type
	Returns   Type
	PreBody   *[]palisade.Expression
	Body      []Expression
}

type MApplication struct {
	Operator MFunction
	Operand  Expression
}

type DApplication struct {
	Operator DFunction
	Left     Expression
	Right    Expression
}

type Int struct {
	Value int64
}

type Real struct {
	Value float64
}

type String struct {
	Value string
}

type Char struct {
	Value string
}

type Bool struct {
	Value bool
}

type Alpha struct{}
type Omega struct{}
type EOF struct{}
