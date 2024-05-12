package prism

import "github.com/llir/llvm/ir/types"

type Group interface {
	String() string
	Universal() bool
	Union(t Group) Group
	Intersection(t Group) Group
	Has(typ Type) bool
}

type Universal struct{}

type TypeGroup struct {
	Set []Type
}

func (g Universal) Kind() int {
	return TypeKindGroup
}
func (g TypeGroup) Kind() int {
	return TypeKindGroup
}

func (g Universal) IsAlgebraic() bool {
	return true
}
func (g TypeGroup) IsAlgebraic() bool {
	return true
}

func (g Universal) Width() int64 {
	panic("Impossible")
}
func (g TypeGroup) Width() int64 {
	panic("Impossible")
}

func (g Universal) Realise() types.Type {
	panic("Impossible")
}
func (g TypeGroup) Realise() types.Type {
	panic("Impossible")
}

// Resolve composes Integrate with Derive,
// Fills in sum/generic type based on a concrete type
func (s Universal) Resolve(t Type) Type {
	return t
}

// Resolve composes Integrate with Derive,
// Fills in sum/generic type based on a concrete type
func (s TypeGroup) Resolve(t Type) Type {
	return Integrate(s, Derive(s, t))
}

func (g Universal) Equals(t Type) bool {
	return t.Kind() == TypeKindGroup && t.(Group).Universal()
}
func (g TypeGroup) Equals(t Type) bool {
	if t.Kind() != TypeKindGroup || t.(Group).Universal() || len(t.(Group).String()) == len(g.String()) {
		return false
	}

	for _, tt := range g.Set {
		if !t.(Group).Has(tt) {
			return false
		}
	}

	return true
}

func (g TypeGroup) String() (s string) {
	for i, t := range g.Set {
		if i > 0 {
			s += " "
		}
		s += t.String()
	}

	return ("{" + s + "}")
}

func (Universal) String() string {
	return "T"
}

func (TypeGroup) Universal() bool {
	return false
}

func (Universal) Universal() bool {
	return true
}

func (g TypeGroup) Union(t Group) Group {
	if t.Universal() {
		return t
	}

	return TypeGroup{append(g.Set, t.(TypeGroup).Set...)}
}

func (g TypeGroup) Intersection(t Group) Group {
	if t.Universal() {
		return g
	}

	var new TypeGroup
	for _, tt := range g.Set {
		if t.Has(tt) {
			new.Set = append(new.Set, tt)
		}
	}

	return new
}

func (group TypeGroup) Has(typ Type) bool {
	for _, iter := range group.Set {
		if iter.Equals(typ) {
			return true
		}
	}

	return false
}

func (g Universal) Union(t Group) Group {
	return g
}

func (g Universal) Intersection(t Group) Group {
	return t
}

func (group Universal) Has(typ Type) bool {
	return true
}
