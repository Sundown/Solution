package prism

func PredicateSemiDeterminedType(t Type) bool {
	if t.Kind() == TypeKindSemiDetermined || t.Kind() == TypeKindSemiDeterminedGroup {
		return true
	}

	if t.Kind() == TypeKindVector {
		return PredicateSemiDeterminedType(t.(VectorType).Type)
	}

	// TODO probably more cases
	return false
}

func DeriveSemiDeterminedType(a, b Type) Type {
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
		for _, t := range a.(SemiDeterminedTypeGroup).Types {
			if c := DeriveSemiDeterminedType(t, b); c != nil {
				// c is t and b, guess this is tidy or something
				return c
			}
		}
	}

	if b.Kind() == TypeKindSemiDeterminedGroup {
		for _, t := range b.(SemiDeterminedTypeGroup).Types {
			if c := DeriveSemiDeterminedType(t, b); c != nil {
				// c is t and a, guess this is tidy or something
				return c
			}
		}
	}

	switch a.(type) {
	case AtomicType:
		return a
	case VectorType:
		return DeriveSemiDeterminedType(a.(VectorType).Type, b.(VectorType).Type)

	}

	return nil
}

func IntegrateSemiDeterminedType(derived Type, semidet Type) Type {
	if semidet.Kind() == TypeKindSemiDetermined {
		return derived
	}

	if semidet.Kind() == TypeKindSemiDeterminedGroup {
		for _, t := range semidet.(SemiDeterminedTypeGroup).Types {
			if c := DeriveSemiDeterminedType(derived, t); c != nil {
				// c is t and b, guess this is tidy or something
				return c
			}
		}
	}

	switch semidet.Kind() {
	case TypeKindAtomic:
		panic("????? what why")
	case TypeKindVector:
		return VectorType{IntegrateSemiDeterminedType(derived, semidet.(VectorType).Type)}
	case TypeKindStruct:
		// TODO
	}

	panic("Failed to integrate")
}
