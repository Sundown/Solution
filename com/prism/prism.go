package prism

import (
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
	Type() Type
	Ident() Ident
	String() string
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

type Morpheme interface {
	_atomicflag()
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

type Alpha struct{}
type Omega struct{}
type EOF struct{}
