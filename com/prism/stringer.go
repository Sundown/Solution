package prism

import "sundown/solution/oversight"

func (a *AtomicType) String() string {
	switch a.ID {
	case TypeInt:
		return "Int"
	case TypeReal:
		return "Real"
	case TypeChar:
		return "Char"
	}

	oversight.Panic("Type name not implemented")
	return ""
}

func (v *VectorType) String() string {
	return "[" + v.ElementType.String() + "]"
}

func (s *StructType) String() (acc string) {
	acc = "("
	for i, v := range s.FieldTypes {
		if i > 0 {
			acc += ", "
		}

		acc += v.String()
	}

	return acc + ")"
}
