package prism

func DeferDyadicApplicationTypes(function *DyadicFunction, x, y *Expression) {
	// Enclose all function-side types in a vector if operands are vectors and the function is not a vector-function
	// TOOD allow X f Y : A f [B] -> A f B.0, A f B.1, ... ,A f B.n
	if !function.NoAutoVector() &&
		QueryAutoVector(function.OmegaType, (*y).Type()) &&
		QueryAutoVector(function.AlphaType, (*x).Type()) {
		function.AlphaType = VectorType{Type: function.AlphaType}
		function.OmegaType = VectorType{Type: function.OmegaType}
		function.Returns = VectorType{Type: function.Returns}
	}

	// There is simple equality between function- and operand-side types (regardless of autovectorisation)
	// ... do nothing and return
	if (*x).Type().Equals(function.AlphaType) && (*y).Type().Equals(function.OmegaType) && !function.Returns.IsAlgebraic() {
		return
	}

	// ... otherwise, cast the operands to the function-side types
	// Method 1:
	//		If function-side type is sum type then substitute if valid, otherwise set if generic.
	// Method 2:
	// 		The function-side type is not algebraic (i.e. sum or generic), however, there is a mapping
	// 		from the operand-side type to the function-side type.
	// Method 3:
	// 		There is no mapping from the operand-side type to the function-side type, however, the
	// 		function-side type is a sum type and it is possible to map the operand-side type to one of
	// 		the types within the sum.
	if newX := function.AlphaType.Resolve((*x).Type()); newX != nil { // 1
		function.AlphaType = newX

	} else if QueryCast((*x).Type(), function.AlphaType) { // 2
		*x = DelegateCast(*x, function.AlphaType)
	} else if cast := RoundhouseCast(*x, (*y).Type(), function.AlphaType); cast != nil { // 3
		*x = cast
	} else {
		Panic("Cannot find mapping between ", (*x).Type(), " and ", function.AlphaType)
	}

	// Same as above
	if newY := function.OmegaType.Resolve((*y).Type()); newY != nil { // 1
		function.OmegaType = newY
	} else if QueryCast((*y).Type(), function.OmegaType) { // 2
		*y = DelegateCast(*y, function.OmegaType)
	} else if cast := RoundhouseCast(*y, (*x).Type(), function.OmegaType); cast != nil { // 3
		*y = *cast
	} else {
		Panic("Cannot find mapping between ", (*y).Type(), " and ", function.OmegaType)
	}

	// Function return type may be reliant on the input type, substitute.
	// TODO more vigorous substitution, giving consideration to LHS type.
	if function.Returns.IsAlgebraic() {
		function.Returns = function.Returns.Resolve((*y).Type())
	}
}

// RoundhouseCast attempts to find a mapping from a type to a type within a sum
// from is the Expression to have it's type changed (note: this will change to a type if casting is re-implemented this way)
// otherside is the opposite-hand type (for DyadicFunctions only), used to give hint
// to is the function-side type which contains the sum type.
func RoundhouseCast(from Expression, otherside Type, to Type) (res *Cast) {
	if !to.IsAlgebraic() {
		panic("RoundhouseCast: to is not algebraic")
	}

	if sum, ok := to.(SumType); ok {
		// Iterate through types within the sum, if one of them is equal to the
		// opposite-hand type, use that, otherwise the first one is used.
		// TODO this is broken currently and is non-deterministic depending on
		// the order of the types in the sum.
		// Probably need type-sorting algorithm
		for _, t := range sum.Types {
			if otherside != nil && otherside.Equals(t) && QueryCast(from.Type(), t) {
				res = &Cast{Value: from, ToType: t}
			} else if QueryCast(from.Type(), t) {
				res = &Cast{Value: from, ToType: t}
			}
		}
	}

	return
}

// May be used for type sorting in future
func IncrementType(t Type) Type {
	switch t.Kind() {
	case BoolType.ID:
		return CharType
	case CharType.ID:
		return IntType
	case IntType.ID:
		return RealType
	case TypeKindVector:
		return VectorType{IncrementType(t.(VectorType).Type)}
	default:
		return nil
	}
}
