package prism

import "fmt"

type DyadicOperator struct {
	Operator int
	Fn       Function
	ExprType Type
	Returns  Type
}

type OperatorApplication struct {
	Op   DyadicOperator
	Expr Expression
}

// String function for interface
func (do DyadicOperator) String() string {
	return do.Fn.String() + " " + fmt.Sprint(do.Operator)
}

// Type property for interface
//
// Operators each return a type dependant on a different input
func (do DyadicOperator) Type() Type {
	switch do.Operator {
	case KindMapOperator:
		return VectorType{do.Fn.Type()}
	case KindReduceOperator:
		return do.Fn.Type()
	}

	Panic("Need to impl type")
	panic(nil)
}

func (d OperatorApplication) String() string {
	return d.Op.String() + " " + d.Expr.String()
}
func (d OperatorApplication) Type() Type {
	return d.Op.Type()
}

func (do DyadicOperator) LLVMise() string {
	return do.Fn.LLVMise() + "_" + fmt.Sprint(do.Operator)
}

func (do DyadicOperator) IsSpecial() bool {
	return false
}
func (do DyadicOperator) ShouldInline() bool {
	return true
}

func (do DyadicOperator) Ident() Ident {
	return Ident{}
}
