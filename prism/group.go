package prism

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
