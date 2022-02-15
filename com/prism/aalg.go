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
