package subtle

import (
	"github.com/sundown/solution/palisade"
	"github.com/sundown/solution/prism"
)

func (env Environment) internFunction(f palisade.Function) {
	var alpha, omega, sigma prism.Type
	var dyadic = false
	var ident prism.Ident
	var body []palisade.Expression

	if p := f.TypedFunction; p != nil {
		if q := p.Dyadic; q != nil {
			dyadic = true
			ident = env.AwareIntern(*q.Ident)
			alpha = env.SubstantiateType(*q.Alpha)
			omega = env.SubstantiateType(*q.Omega)
		} else if q := p.Monadic; q != nil {
			ident = env.AwareIntern(*q.Ident)
			omega = env.SubstantiateType(*q.Omega)
		}

		sigma = env.SubstantiateType(*p.Returns)
	} else if p := f.AmbiguousFunction; p != nil {
		dyadic = env.determineArity(&f)
		alpha = prism.Universal{}
		omega = prism.Universal{}
		if p.Returns != nil {
			sigma = env.SubstantiateType(*p.Returns)
		} else {
			sigma = prism.Universal{}
		}
		ident = env.AwareIntern(*p.Ident)
	} else {
		panic("unreachable")
	}

	if _, ok := env.DyadicFunctions[ident]; ok {
		prism.Panic("Dyadic function " + ident.String() + " is already defined")
	}

	if _, ok := env.MonadicFunctions[ident]; ok {
		prism.Panic("Monadic function " +
			ident.String() + " is already defined")
	}

	if f.Tacit != nil {
		body = []palisade.Expression{*f.Tacit}
	} else if f.Body != nil {
		body = *f.Body
	} else {
		panic("unreachable")
	}

	if dyadic {
		fn := prism.DyadicFunction{
			Name:      ident,
			AlphaType: alpha,
			OmegaType: omega,
			Returns:   sigma,
			PreBody:   &body,
		}

		env.DyadicFunctions[fn.Name] = &fn
	} else {
		fn := prism.MonadicFunction{
			Name:      ident,
			OmegaType: omega,
			Returns:   sigma,
			PreBody:   &body,
		}

		env.MonadicFunctions[fn.Name] = &fn
	}
}
