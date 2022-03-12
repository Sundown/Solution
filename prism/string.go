package prism

import "github.com/alecthomas/repr"

var StringType = VectorType{Type: CharType}

// Interface Type comparison
func (t String) Equals(b Type) bool {
	return b.Kind() == TypeString
}

// Resolve composes Integrate with Derive,
// Should not be used on concrete types
func (i String) Resolve(t Type) Type {
	Panic("Unreachable")
	panic(nil)
}

type String struct {
	Value string
}

// String function for interface
func (s String) String() string {
	return "\"" + s.Value + "\""
}

// Type property for interface
func (s String) Type() Type {
	repr.Println(s)
	return StringType
}

// Interface Type width for LLVM codegen
func (t String) Width() int64 {
	return t.Type().(VectorType).Width()
}

// Interface Type algebraic predicate
func (t String) IsAlgebraic() bool {
	return false
}

func (t String) _atomicflag() { return }
