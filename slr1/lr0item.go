package slr1

import (
	"fmt"

	"github.com/mbotezatu/slr1-table/grammar"
)

// LR0item defines the structure of a LR(0) item
type LR0item struct {
	pPos int
	rule *grammar.Production
}

// NewLR0item creates and returns a LR0item
func NewLR0item(pPos int, rule *grammar.Production) (l *LR0item) {
	if pPos > len(rule.RHS()) {
		pPos = len(rule.RHS())
	}
	if rule.RHS()[0] == "" {
		pPos = 1
	}
	l = &LR0item{pPos, rule}
	return
}

// AtPPos returns the terminal or non-terminal at that position,
// or empty string if there isn't any
func (l *LR0item) AtPPos() (s string) {
	if l.pPos < len(l.rule.RHS()) {
		s = l.rule.RHS()[l.pPos]
	}
	return
}

// Equals returns true if the two items are equal,
// false otherwise
func (l *LR0item) Equals(i *LR0item) (yes bool) {
	if l.pPos == i.pPos && l.rule == i.rule {
		yes = true
	}
	return
}

func (l *LR0item) String() string {
	return fmt.Sprintf("Position: %d\nRule: %s", l.pPos, l.rule.String())
}
