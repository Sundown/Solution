package subtle

import (
	"fmt"

	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env *Environment) createMapOperator(function prism.MonadicFunction, rType prism.Type) prism.MonadicOperator {

	if !rType.Equals(function.OmegaType) {
		if !prism.QueryCast(rType, function.OmegaType) {
			_, err := prism.Delegate(function.OmegaType, rType)
			if err != nil {
				panic(*err)
			}
		}
	}

	if function.Returns.IsAlgebraic() {
		function.Returns = function.Returns.Resolve(function.OmegaType)
	}

	retType := prism.Type(prism.VectorType{Type: function.Type()})
	if retType.(prism.VectorType).SubIsVoid() {
		retType = prism.VoidType
	}

	return prism.MonadicOperator{
		Operator: prism.KindMapOperator,
		Fn:       function,
		ExprType: prism.VectorType{Type: rType},
		Returns:  retType,
	}
}

func (env *Environment) createReduceOperator(function prism.DyadicFunction, rType prism.Type) prism.MonadicOperator {
	if !rType.Equals(function.OmegaType) {
		if !prism.QueryCast(rType, function.OmegaType) {

			j, err := prism.Delegate(function.OmegaType, rType)
			function.OmegaType = j
			if err != nil {
				prism.Panic(*err)
			}
		}
	}

	if !rType.Equals(function.AlphaType) {
		if !prism.QueryCast(rType, function.AlphaType) {
			j, err := prism.Delegate(function.AlphaType, rType)
			function.AlphaType = j
			if err != nil {
				prism.Panic(*err)
			}
		}
	}

	if j, err := prism.Delegate(function.AlphaType, function.OmegaType); err == nil {
		function.AlphaType = j
	} else {
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

		function := env.analysePrimeApplicable(app, nil, typ.Type)

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

func (env *Environment) monadicOperatorToFunction(op prism.MonadicOperator) prism.MonadicFunction {
	fn := prism.MonadicFunction{
		Attribute: prism.Attribute{
			SkipBuilder: true,
			ForceInline: true,
		},
		Name:      prism.Ident{Package: "_", Name: "m_op_" + fmt.Sprint(env.Iterate())},
		OmegaType: op.ExprType,
		Returns:   op.Returns,
		PreBody:   nil,
		Body: []prism.Expression{
			prism.MonadicApplication{
				Operator: prism.MonadicFunction{
					Name:    prism.Ident{Package: "_", Name: "←"},
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
