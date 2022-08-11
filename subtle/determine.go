package subtle

import (
	"github.com/sundown/solution/palisade"
)

func (env Environment) determineArity(f *palisade.Function) bool {
	if f.Tacit != nil {
		return env.determineExpression(f.Tacit)
	}

	for _, expr := range *f.Body {
		if env.determineExpression(&expr) {
			return true
		}
	}

	return false
}

func (env Environment) determineExpression(e *palisade.Expression) bool {
	if e.Monadic != nil {
		return env.determineMonadic(e.Monadic)
	} else if e.Dyadic != nil {
		return env.determineDyadic(e.Dyadic)
	} else if e.Morphemes != nil {
		return env.determineMorphemes(e.Morphemes)
	}

	panic("Unreachable")
}

func (env Environment) determineMorphemes(ms *palisade.Morpheme) bool {
	return ms.Alpha != nil
}

func (env Environment) determineDyadic(d *palisade.Dyadic) bool {
	var alpha bool
	if d.Monadic != nil {
		alpha = false
	} else if d.Morphemes != nil {
		alpha = env.determineMorphemes(d.Morphemes)
	}

	omega := env.determineExpression(d.Expression)

	return alpha || omega
}

func (env Environment) determineMonadic(d *palisade.Monadic) bool {
	return env.determineExpression(d.Expression)
}
