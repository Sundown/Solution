package prism

// I have no clue what this does
func Delegate(mould, cast *Type) (determined Type, failure *string) {
	if mould == nil {
		return nil, Ref("mould is nil")
	} else if cast == nil {
		return nil, Ref("cast is nil")
	}

	if _, sd := (*cast).(GenericType); sd {
		return nil, Ref("Cast cannot be T: " + (*cast).String())
	}

	if _, sdg := (*cast).(Group); sdg {
		return nil, Ref("Cast cannot be sum: " + (*cast).String())
	}

	// First
	switch (*mould).(type) {
	case AtomicType:
		if mt, ok := (*cast).(AtomicType); ok {
			if (*mould).(AtomicType).ID == mt.ID {
				temp := Type(mt)
				return temp, nil // Success; matched atomic types together
			} else {
				return nil, Ref("Atomic type mismatch")
			}
		} else {
			return nil, Ref("Type class mismatch, cast is not atomic")
		}
	case VectorType:
		if vt, vtp := (*cast).(VectorType); vtp {
			temp := (*mould).(VectorType).Type
			del, err := Delegate(&temp, &vt.Type)
			if err != nil {
				return nil, err
			}

			*mould = *cast
			return del, nil
		} else {
			return nil, Ref("Type class mismatch, cast is not a vector")
		}
	}

	// Second
	// T has been matched with a determined type directly
	if _, tp := (*mould).(GenericType); tp {
		*mould = *cast
		return *cast, nil
	}

	if group, tgp := (*mould).(Group); tgp {
		if group.Universal() {
			panic("TODO")
		}

		for _, elm := range group.(TypeGroup).Set {
			typ, fail := Delegate(&elm, cast)

			// Errors are expected, don't check for them
			if fail == nil {
				*mould = *cast
				return typ, nil
			}
		}

		return nil, Ref("Cast does not fit within algebraic group")
	}

	Panic("unreachable")
	panic("Unknown error")
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
