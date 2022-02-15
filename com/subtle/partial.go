package subtle

import (
	"sundown/solution/palisade"
	"sundown/solution/prism"
)

func (env Environment) AnalysePartial(d *palisade.Dyadic) prism.MonadicFunction {
	op := env.FetchVerb(d.Verb)
	if _, ok := op.(prism.DyadicFunction); !ok {
		panic("Verb is not a dyadic function")
	}
	var left prism.Expression
	if d.Monadic != nil {
		left = env.AnalyseMonadic(d.Monadic)
	} else if d.Morphemes != nil {
		left = env.AnalyseMorphemes(d.Morphemes)
	} else {
		panic("Dyadic expression has no left operand")
	}

	fn := op.(prism.DyadicFunction)

	pred := fn.AlphaType.IsAlgebraic()

	tmp := left.Type()
	resolved_left, err := prism.Delegate(&fn.AlphaType, &tmp)
	if err != nil {
		prism.Panic(*err)
	}

	if fn.Returns.IsAlgebraic() {
		fn.Returns = fn.Returns.Resolve(resolved_left)
	}

	var takes prism.Type
	if pred {
		// TODO sort this out
		takes = fn.OmegaType.Resolve(resolved_left)
		fn.OmegaType = fn.OmegaType.Resolve(resolved_left)
	}

	dapp := prism.DyadicApplication{
		Operator: fn,
		Left:     left,
		Right:    nil,
	}
	//------
	mon := prism.MonadicFunction{
		Special:   false,
		Name:      prism.Ident{Package: "_", Name: "inlined_partial"},
		Returns:   dapp.Operator.Returns,
		OmegaType: dapp.Operator.OmegaType,
	}

	mon.Body = []prism.Expression{
		prism.MonadicApplication{
			Operator: prism.MonadicFunction{
				Special: false,
				Name:    prism.Ident{Package: "_", Name: "Return"},
				Returns: dapp.Operator.Returns,
			},
			Operand: prism.DyadicApplication{
				Operator: dapp.Operator,
				Left:     dapp.Left,
				Right:    prism.Omega{TypeOf: takes},
			},
		},
	}

	env.MonadicFunctions[mon.Ident()] = &mon

	return mon
}
