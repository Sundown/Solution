package prism

// I have no clue what this does
func Delegate(algebraic, concrete Type) (determined Type, failure *string) {
	if algebraic == nil {
		return nil, Ref("mould is nil")
	} else if concrete == nil {
		return nil, Ref("cast is nil")
	}

	if _, sd := (concrete).(GenericType); sd {
		return nil, Ref("Cast cannot be T: " + (concrete).String())
	}

	if _, sdg := (concrete).(Group); sdg {
		return nil, Ref("Cast cannot be sum: " + (concrete).String())
	}

	// First
	switch m := algebraic.(type) {
	case AtomicType:
		if _, ok := concrete.(AtomicType); !ok {
			return nil, Ref("Type class mismatch, cast is not atomic")
		}

		if !algebraic.Equals(concrete) {
			return nil, Ref("Atomic type mismatch")
		}

		return concrete, nil
	case VectorType:
		if _, ok := concrete.(VectorType); !ok {
			return nil, Ref("Type class mismatch, cast is not a vector")
		}

		// TODO maybe this should be `return Vector(Delegate(...))?`
		return Delegate(concrete.(VectorType), m.Type)
	case GenericType:
		return concrete, nil
	case Group:
		if m.Universal() {
			return concrete, nil
		}

		if !m.Has(concrete) {
			return nil, Ref("Cast does not fit within algebraic group")
		}

		return concrete, nil

	}

	panic("Unreachable")
}

// Integrate a concrete type into a sum or generic type
func Integrate(this, from Type) Type {
	switch j := this.(type) {
	case VectorType:
		if j.IsAlgebraic() {
			return Integrate(j, from)
		} else {
			Panic("Vector is not algebraic")
		}
	case Group:
		if j.Universal() {
			return from
		}

		if j.Has(from) {
			return from
		}
	case GenericType:
		return from
	}

	return nil
}

// Derive concrete type based on likeness of generic/sum type
func Derive(this, like Type) Type {
	switch j := this.(type) {
	case VectorType:
		if !j.IsAlgebraic() {
			Panic("Vector is not algebraic")
		} else if like.IsAlgebraic() {
			panic("Cannot derive algebraic type from algebraic type")
		}

		if v, ok := like.(VectorType); ok {
			return Derive(j.Type, v.Type)
		} else {
			return Derive(j.Type, like)
		}
	case Group:
		if j.Universal() {
			return like
		}

		if like.IsAlgebraic() {
			panic("Cannot derive algebraic type from algebraic type")
		}

		if j.Has(like) {
			return like
		}
	case GenericType:
		return like
	}

	return nil
}
