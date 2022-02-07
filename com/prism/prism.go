package prism

import (
	"fmt"
	"sundown/solution/palisade"

	"github.com/llir/llvm/ir"
)

type Environment struct {
	LexResult *palisade.PalisadeResult
	//
	MonadicFunctions map[Ident]*MonadicFunction
	DyadicFunctions  map[Ident]*DyadicFunction
	Types            map[Ident]Type
	//
	EmitFormat   string
	Output       string
	Verbose      *bool
	Optimisation *int64
	File         string
	//
	EntryFunction      MonadicFunction
	Module             *ir.Module
	Block              *ir.Block
	LLDyadicFunctions  map[string]*ir.Func
	LLMonadicFunctions map[string]*ir.Func
	Specials           map[string]*ir.Func
	CurrentFunction    *ir.Func
	CurrentFunctionIR  Expression
	PanicStrings       map[string]*ir.Global
}

type Ident struct {
	Package string
	Name    string
}

type Function interface {
	LLVMise() string
	IsSpecial() bool
	Type() Type
	Ident() Ident
	String() string
}

func (d DyadicFunction) IsSpecial() bool {
	return d.Special
}

func (d MonadicFunction) IsSpecial() bool {
	return d.Special
}

func (d DyadicFunction) Ident() Ident {
	return d.Name
}

func (m MonadicFunction) Ident() Ident {
	return m.Name
}

func (f DyadicFunction) LLVMise() string {
	return f.Name.Package + "::" + f.Name.Name + "_" + f.AlphaType.String() + "," + f.OmegaType.String() + "->" + f.Returns.String()
}

func (f MonadicFunction) LLVMise() string {
	return f.Name.Package + "::" + f.Name.Name + "_" + f.OmegaType.String() + "->" + f.Returns.String()
}

type Vector struct {
	ElementType VectorType
	Body        *[]Expression
}

type Expression interface {
	Type() Type
	String() string
}

type Void struct{}

func (Void) Type() Type {
	return VoidType
}

func (Void) String() string {
	return "Void"
}

type DyadicOperator struct {
	Operator int
	Left     Expression
	Right    Expression
	Returns  Type
}

func (do DyadicOperator) Type() Type {
	switch do.Operator {
	case KindMapOperator:
		return VectorType{do.Left.Type()}
	case KindFoldlOperator:
		return do.Left.Type()
	}

	panic("Need to impl type")
}

func (do DyadicOperator) String() string {
	return do.Left.String() + " " + fmt.Sprint(do.Operator) + " " + do.Right.String()
}

type Morpheme interface {
	_atomicflag()
}

type Cast struct {
	Value  Expression
	ToType Type
}

func (c Cast) Type() Type {
	return c.ToType
}

func (c Cast) String() string {
	return "<" + c.ToType.String() + ">" + c.Value.String()
}

type DyadicFunction struct {
	Special   bool
	Name      Ident
	AlphaType Type
	OmegaType Type
	Returns   Type
	PreBody   *[]palisade.Expression
	Body      []Expression
}

type MonadicFunction struct {
	Special   bool
	Name      Ident
	OmegaType Type

	Returns Type
	PreBody *[]palisade.Expression
	Body    []Expression
}

type MApplication struct {
	Operator MonadicFunction
	Operand  Expression
	// Algebraic T refers to which type in given application
	T_Refers Type
}

type DApplication struct {
	Operator DyadicFunction
	Left     Expression
	Right    Expression
	// Algebraic T refers to which type in given application
	// Possibility of requiring A and B algebraic types once tuples are present
	T_Refers Type
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

type Alpha struct {
	TypeOf Type
}
type Omega struct {
	TypeOf Type
}

func (a Alpha) Type() Type {
	return a.TypeOf
}

func (a Omega) Type() Type {
	return a.TypeOf
}

type EOF struct{}
