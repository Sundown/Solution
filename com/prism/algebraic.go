package prism

func PredicateGenericType(t Type) bool {
	if t.Kind() == TypeKindSemiDetermined || t.Kind() == TypeKindSemiDeterminedGroup {
		return true
	}

	if t.Kind() == TypeKindVector {
		return PredicateGenericType(t.(VectorType).Type)
	}

	// TODO probably more cases
	return false
}

func DeriveGenericType(a, b Type) Type {
	if a.Kind() == TypeKindSemiDetermined && b.Kind() == TypeKindSemiDetermined {
		panic("Impossible to derive T from T and T, must provide fully substantiated type")
	}

	if a.Kind() == TypeKindSemiDetermined {
		return b
	}

	if b.Kind() == TypeKindSemiDetermined {
		return a
	}

	if a.Kind() == TypeKindSemiDeterminedGroup {
		for _, t := range a.(SumType).Types {
			if c := DeriveGenericType(t, b); c != nil {
				// c is t and b, guess this is tidy or something
				return c
			}
		}
	}

	if b.Kind() == TypeKindSemiDeterminedGroup {
		for _, t := range b.(SumType).Types {
			if c := DeriveGenericType(t, b); c != nil {
				// c is t and a, guess this is tidy or something
				return c
			}
		}
	}

	switch a.(type) {
	case AtomicType:
		return a
	case VectorType:
		return DeriveGenericType(a.(VectorType).Type, b.(VectorType).Type)

	}

	return nil
}

func IntegrateGenericType(derived Type, semidet Type) Type {
	if semidet.Kind() == TypeKindSemiDetermined {
		return derived
	}

	if semidet.Kind() == TypeKindSemiDeterminedGroup {
		for _, t := range semidet.(SumType).Types {
			if c := DeriveGenericType(derived, t); c != nil {
				// c is t and b, guess this is tidy or something
				return c
			}
		}
	}

	switch semidet.Kind() {
	case TypeKindAtomic:
		panic("????? what why")
	case TypeKindVector:
		return VectorType{IntegrateGenericType(derived, semidet.(VectorType).Type)}
	case TypeKindStruct:
		// TODO
	}

	panic("Failed to integrate")
}
