package subtle

import (
	"fmt"

	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env *Environment) createMapOperator(function prism.MonadicFunction, exprType prism.Type) prism.DyadicOperator {
	if !exprType.Equals(function.OmegaType) {
		if !prism.QueryCast(exprType, function.OmegaType) {
			tmp := exprType
			_, err := prism.Delegate(&function.OmegaType, &tmp)
			if err != nil {
				prism.Panic(*err)
			}
		}
	}

	if function.Returns.IsAlgebraic() {
		function.Returns = function.Returns.Resolve(function.OmegaType)
	}

	return prism.DyadicOperator{
		Operator: prism.KindMapOperator,
		Fn:       function,
		ExprType: prism.VectorType{Type: exprType},
		Returns:  function.Type(),
	}
}

func (env *Environment) createReduceOperator(function prism.DyadicFunction, exprType prism.Type) prism.DyadicOperator {
	if !exprType.Equals(function.OmegaType) {
		if !prism.QueryCast(exprType, function.OmegaType) {
			tmp := exprType
			_, err := prism.Delegate(&function.OmegaType, &tmp)
			if err != nil {
				prism.Panic(*err)
			}
		}
	}

	if !exprType.Equals(function.AlphaType) {
		if !prism.QueryCast(exprType, function.AlphaType) {
			tmp := exprType
			_, err := prism.Delegate(&function.AlphaType, &tmp)
			if err != nil {
				prism.Panic(*err)
			}
		}
	}

	if _, err := prism.Delegate(&function.AlphaType, &function.OmegaType); err != nil {
		prism.Panic(*err)
	}

	if function.Returns.IsAlgebraic() {
		function.Returns = function.Returns.Resolve(function.AlphaType)
	}
	return prism.DyadicOperator{
		Operator: prism.KindReduceOperator,
		Fn:       function,
		ExprType: prism.VectorType{Type: exprType},
		Returns:  function.Type(),
	}
}

func (env *Environment) analyseDyadicOperator(d *palisade.Operator, exprType prism.Type) prism.DyadicOperator {
	switch *d.Operator {
	case "¨":
		if d.Subexpr != nil {
			prism.Panic("Not implemented")
		}

		if !prism.IsVector(exprType) {
			prism.Panic("Right operand is not a vector")
		}

		return env.createMapOperator(
			env.FetchMVerb(d.Verb),
			exprType.(prism.VectorType).Type)
	case "/":
		var lexpr prism.Expression
		if d.Verb != nil {
			lexpr = env.FetchDVerb(d.Verb)
		}
		if _, ok := exprType.(prism.VectorType); !ok {
			prism.Panic("Right operand is not a vector")
		}

		if d.Subexpr != nil {
			t := exprType.(prism.VectorType).Type
			lexpr = env.analyseDyadicPartial(d.Subexpr, t, t)
		}

		if _, ok := lexpr.(prism.DyadicFunction); !ok {
			prism.Panic("Left operand is not a function")
		}

		fn := lexpr.(prism.DyadicFunction)
		return env.createReduceOperator(
			fn,
			exprType.(prism.VectorType).Type)
	}

	panic("Unknown operator")
}

func (env Environment) analyseDyadicApplication(d *palisade.Operator) prism.OperatorApplication {
	rexpr := env.analyseExpression(d.Expression)
	return prism.OperatorApplication{
		Op:   env.analyseDyadicOperator(d, rexpr.Type()),
		Expr: rexpr,
	}
}

func (env Environment) operatorToFunction(op prism.DyadicOperator) prism.MonadicFunction {
	fn := prism.MonadicFunction{
		Special:     false,
		SkipBuilder: true,
		Inline:      true,
		Name:        prism.Ident{Package: "_", Name: "m_op_" + fmt.Sprint(env.Iterate())},
		OmegaType:   op.ExprType,
		Returns:     op.Returns,
		PreBody:     nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Special: false,
					Name:    prism.Ident{Package: "_", Name: "Return"},
					Returns: op.Returns,
				},
				Operand: prism.OperatorApplication{
					Op:   op,
					Expr: prism.Omega{TypeOf: op.ExprType},
				},
			},
		},
	}

	env.MonadicFunctions[fn.Name] = &fn
	return fn
}
