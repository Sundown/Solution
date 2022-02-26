package prism

/* This is perfect in every single way */
func Delegate(mould, cast *Type) (determined Type, failure *string) {
	if mould == nil {
		return nil, Ref("mould is nil")
	} else if cast == nil {
		return nil, Ref("cast is nil")
	}

	if _, sd := (*cast).(GenericType); sd {
		return nil, Ref("Cast cannot be T: " + (*cast).String())
	}
	if _, sdg := (*cast).(SumType); sdg {
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

	if group, tgp := (*mould).(SumType); tgp {
		for _, elm := range group.Types {
			typ, fail := Delegate(&elm, cast)

			// Errors are expected, don't check for them
			if fail == nil {
				*mould = *cast
				return typ, nil
			}
		}

		return nil, Ref("Cast does not fit within algebraic group")
	}

	panic("unreachable")
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

// Derive concrete type based on likeness of generic/sum type
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
			if e.Equals(like) {
				return e
			}
		}
	case GenericType:
		return like
	}

	panic("Unreachable")
}
