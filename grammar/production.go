package grammar

import (
	"bytes"
	"strings"
)

// Production defines the structure of
// production rules for context free grammars.
type Production struct {
	lhs string
	rhs []string
}

// NewProduction creates and returns a Production pointer
func NewProduction(lhs string, rhs []string) (p *Production) {
	p = &Production{lhs, rhs}
	return
}

// LHS returns the left-hand side of a production rule
func (p *Production) LHS() string {
	return p.lhs
}

// RHS returns the right-hand side of a production rule
func (p *Production) RHS() []string {
	return p.rhs
}

func (p *Production) String() string {
	var b bytes.Buffer
	b.WriteString(p.lhs)
	b.WriteString(" -> ")
	rhs := strings.Join(p.rhs, " ")
	b.WriteString(rhs)

	return b.String()
}
