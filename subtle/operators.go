package subtle

import (
	"fmt"

	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env *Environment) createMapOperator(function prism.MonadicFunction, rType prism.Type) prism.MonadicOperator {
	if !rType.Equals(function.OmegaType) {
		if !prism.QueryCast(rType, function.OmegaType) {
			tmp := rType
			_, err := prism.Delegate(&function.OmegaType, &tmp)
			if err != nil {
				panic(*err)
			}
		}
	}

	if function.Returns.IsAlgebraic() {
		function.Returns = function.Returns.Resolve(function.OmegaType)
	}

	return prism.MonadicOperator{
		Operator: prism.KindMapOperator,
		Fn:       function,
		ExprType: prism.VectorType{Type: rType},
		Returns:  function.Type(),
	}
}

func (env *Environment) createReduceOperator(function prism.DyadicFunction, rType prism.Type) prism.MonadicOperator {
	if !rType.Equals(function.OmegaType) {
		if !prism.QueryCast(rType, function.OmegaType) {
			tmp := rType
			_, err := prism.Delegate(&function.OmegaType, &tmp)
			if err != nil {
				prism.Panic(*err)
			}
		}
	}

	if !rType.Equals(function.AlphaType) {
		if !prism.QueryCast(rType, function.AlphaType) {
			tmp := rType
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
	return prism.MonadicOperator{
		Operator: prism.KindReduceOperator,
		Fn:       function,
		ExprType: prism.VectorType{Type: rType},
		Returns:  function.Type(),
	}
}

func (env *Environment) analyseMonadicOperator(app palisade.Applicable, rType prism.Type) prism.MonadicOperator {
	switch *app.Operator.Operator {
	case "¨":
		typ, ok := rType.(prism.VectorType)
		if !ok {
			prism.Panic("Right operand is not a vector")
		}

		function := env.analysePrimeApplicable(app, nil, typ)

		if _, ok := function.(prism.MonadicFunction); !ok {
			prism.Panic("Right operand is not a monadic function")
		}

		return env.createMapOperator(
			function.(prism.MonadicFunction),
			rType.(prism.VectorType).Type)
	case "/":
		typ, ok := rType.(prism.VectorType)
		if !ok {
			prism.Panic("Right operand is not a vector")
		}

		function := env.analysePrimeApplicable(app, typ.Type, typ.Type)

		if _, ok := function.(prism.DyadicFunction); !ok {
			prism.Panic("Left operand is not a function")
		}

		return env.createReduceOperator(
			function.(prism.DyadicFunction),
			rType.(prism.VectorType).Type)
	}

	panic("Unknown operator")
}

/* func (env Environment) analyseDyadicApplication(d *palisade.Operator) prism.OperatorApplication {
	rexpr := env.analyseExpression(d.Expression)
	return prism.OperatorApplication{
		Op:   env.analyseMonadicOperator(d, rexpr.Type()),
		Expr: rexpr,
	}
} */

func (env Environment) monadicOperatorToFunction(op prism.MonadicOperator) prism.MonadicFunction {
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
