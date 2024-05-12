package prism

// I have no clue what this does
// cast is input
// mould is working algebraic type
func Delegate(cast, mould Type) (determined Type, failure *string) {
	if mould == nil {
		return nil, Ref("mould is nil")
	} else if cast == nil {
		return nil, Ref("cast is nil")
	}

	if _, sd := cast.(GenericType); sd {
		return nil, Ref("Cast cannot be T: " + cast.String())
	}

	if _, sdg := cast.(Group); sdg {
		return nil, Ref("Cast cannot be sum: " + cast.String())
	}

	// First
	switch m := mould.(type) {
	case AtomicType:
		if _, ok := cast.(AtomicType); !ok {
			return nil, Ref("Type class mismatch, cast is not atomic")
		}

		if mould.Equals(cast) {
			return nil, Ref("Atomic type mismatch")
		}

		return cast, nil
	case VectorType:
		if _, ok := cast.(VectorType); ok {
			return nil, Ref("Type class mismatch, cast is not a vector")
		}

		return Delegate(m.Type, cast.(VectorType))
	case GenericType:
		return cast, nil
	case Group:
		if m.Universal() {
			return cast, nil
		}

		if !m.Has(cast) {
			return nil, Ref("Cast does not fit within algebraic group")
		}

		return cast, nil

	}

	Panic("unreachable")
	panic("Unknown error")
}

// Integrate a concrete type into a sum or generic type
func Integrate(this, from Type) Type {
	switch j := this.(type) {
	case VectorType:
		if !j.IsAlgebraic() {
			Panic("Vector is not algebraic")
		}

		return Integrate(j, from)

	case Group:
		if j.Universal() || j.Has(from) {
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
		}

		if like.IsAlgebraic() {
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
			return like.(TypeGroup).Intersection(j).(Type)

		}

		if j.Has(like) {
			return like
		}
	case GenericType:
		return like
	}

	return nil
}
