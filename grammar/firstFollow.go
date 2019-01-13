package grammar

import (
	"fmt"
	"strings"

	"github.com/mbotezatu/slr1-table/containers"
)

// First returns the set of all terminal symbols that can start sentences
// derived from the right-hand side of a production rule.
func First(cfg *CFG) (f map[string]*containers.Set, err error) {
	f = make(map[string]*containers.Set)

	next := cfg.Terminals.Iterator()
	for elem, ok := next(); ok; elem, ok = next() {
		f[elem] = containers.NewSet()
		f[elem].Add(elem)
	}

	next = cfg.NonTerminals.Iterator()
	for elem, ok := next(); ok; elem, ok = next() {
		f[elem] = containers.NewSet()
	}

	changed := true
	for changed {
		changed = false
		for _, v := range cfg.Rules {
			for _, p := range v {
				rhs := containers.NewSet()
				rhs.AddAll(f[p.rhs[0]])
				rhs.Delete("")
				i := 0
				k := len(p.rhs) - 1
				for i <= k-1 && f[p.rhs[i]].Contains("") {
					aux := containers.NewSet()
					aux.AddAll(f[p.rhs[i+1]])
					aux.Delete("")
					rhs.AddAll(aux)
					i++
				}

				if i == k && f[p.rhs[k]].Contains("") {
					rhs.Add("")
				}
				prevLen := f[p.lhs].Len()
				f[p.lhs].AddAll(rhs)
				if prevLen != f[p.lhs].Len() {
					changed = true
				}
			}
		}
	}

	var b strings.Builder
	for k, v := range f {
		if v.Len() == 0 {
			b.WriteString(k)
			b.WriteString(" ")
		}
	}

	if b.Len() > 0 {
		err = fmt.Errorf("There are unrealizable nonterminals in the grammar: %s", b.String())
	}

	return
}

// Follow :
func Follow(cfg *CFG, first map[string]*containers.Set) (follow map[string]*containers.Set) {
	follow = make(map[string]*containers.Set)

	next := cfg.NonTerminals.Iterator()
	for elem, ok := next(); ok; elem, ok = next() {
		follow[elem] = containers.NewSet()
	}
	follow[cfg.StartRule.lhs].Add("$")

	changed := true
	for changed {
		changed = false
		for _, v := range cfg.Rules {
			for _, r := range v {
				trailer := containers.NewSet()
				trailer.AddAll(follow[r.lhs])
				for i := len(r.rhs) - 1; i >= 0; i-- {
					if cfg.NonTerminals.Contains(r.rhs[i]) {
						prevLen := follow[r.rhs[i]].Len()
						follow[r.rhs[i]].AddAll(trailer)
						if prevLen != follow[r.rhs[i]].Len() {
							changed = true
						}
						if first[r.rhs[i]].Contains("") {
							aux := containers.NewSet()
							aux.AddAll(first[r.rhs[i]])
							aux.Delete("")
							trailer.AddAll(aux)
						} else {
							trailer = containers.NewSet()
							trailer.AddAll(first[r.rhs[i]])
						}
					} else {
						trailer = containers.NewSet()
						trailer.AddAll(first[r.rhs[i]])
					}
				}
			}
		}
	}

	return
}
