package prism

import "fmt"

type MonadicOperator struct {
	Operator int
	Fn       Function
	ExprType Type
	Returns  Type
}

type OperatorApplication struct {
	Op   MonadicOperator
	Expr Expression
}

// String function for interface
func (do MonadicOperator) String() string {
	return do.Fn.String() + " " + fmt.Sprint(do.Operator)
}

// Type property for interface
//
// Operators each return a type dependant on a different input
func (do MonadicOperator) Type() Type {
	switch do.Operator {
	case KindMapOperator:
		return VectorType{do.Fn.Type()}
	case KindReduceOperator:
		return do.Fn.Type()
	}

	Panic("Need to impl type")
	panic("Unknown error")
}

func (d OperatorApplication) String() string {
	return d.Op.String() + " " + d.Expr.String()
}
func (d OperatorApplication) Type() Type {
	return d.Op.Type()
}

func (do MonadicOperator) LLVMise() string {
	return do.Fn.LLVMise() + "_" + fmt.Sprint(do.Operator)
}

func (do MonadicOperator) IsSpecial() bool {
	return false
}
func (do MonadicOperator) ShouldInline() bool {
	return true
}

func (do MonadicOperator) Ident() Ident {
	return Ident{}
}
