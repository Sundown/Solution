package reform

import "fmt"

func (b Application) String() string {
	return b.Operator + " (" + b.Operand.String() + ")"
}

func (u Dangle) String() string {
	return "(" + u.Outer.String() + " " + u.Inner.String() + ")"
}

func (i Int) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i Ident) String() string {
	return i.Value
}

func (e EOF) String() string {
	return "EOF"
}

func (s Subexpression) String() string {
	return "(" + s.Expression.String() + ")"
}

func (r Real) String() string {
	return fmt.Sprintf("%f", r.Value)
}

func (s String) String() string {
	return "\"" + s.Value + "\""
}

func (c Char) String() string {
	return "'" + string(c.Value) + "'"
}

func (a Alpha) String() string {
	return "α"
}

func (o Omega) String() string {
	return "ω"
}
