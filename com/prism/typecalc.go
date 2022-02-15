package prism

func Integrate(this, from Type) Type {
	switch j := this.(type) {
	case VectorType:
		if j.IsAlgebraic() {
			return Integrate(j, from)
		} else {
			Panic("Vector is not algebraic")
		}
	case SumType:
		for _, e := range j.Types {
			if e.Equals(from) {
				return e
			}
		}
	case GenericType:
		return from
	}
	panic("Unreachable")

}

func Derive(this, like Type) Type {
	switch j := this.(type) {
	case VectorType:
		if !j.IsAlgebraic() {
			Panic("Vector is not algebraic")
		} else if like.IsAlgebraic() {
			Panic("Cannot derive algebraic type from algebraic type")
		}

		if v, ok := like.(VectorType); ok {
			return Derive(j.Type, v.Type)
		}

		return nil
	case SumType:
		for _, e := range j.Types {
			if e.Equals(like) { // TODO might be shortsighted
				return e
			}
		}
	case GenericType:
		return like
	}

	panic("Unreachable")
}
