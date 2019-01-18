package slr1

import (
	"fmt"

	"github.com/mbotezatu/slr1-table/containers"
	"github.com/mbotezatu/slr1-table/grammar"
)

// Closure returns the closure set
func Closure(i []*LR0item, cfg *grammar.CFG) (j []*LR0item) {
	j = append(j, i...)
	added := true
	for added {
		added = false
		for k := 0; k < len(j); k++ {
			if j[k].AtPPos() != "" {
				for _, v := range cfg.Rules[j[k].AtPPos()] {
					auxItem := NewLR0item(0, v)
					if !containsItem(j, auxItem) {
						j = append(j, auxItem)
						added = true
					}
				}
			}
		}
	}
	return
}

// Goto returns the closure set of
// all the items of i that accept x
func Goto(i []*LR0item, x string, cfg *grammar.CFG) (j []*LR0item) {
	var aux []*LR0item
	for _, v := range i {
		if v.AtPPos() == x {
			aux = append(aux, NewLR0item(v.pPos+1, v.rule))
		}
	}

	j = Closure(aux, cfg)
	return
}

// CanonicalCollection returns the canonical colection of LR0 items
func CanonicalCollection(cfg *grammar.CFG) (cc [][]*LR0item, gotos map[int]map[string]int) {
	cc = append(cc, Closure([]*LR0item{(NewLR0item(0, cfg.StartRule))}, cfg))
	gotos = make(map[int]map[string]int)
	for i := 0; i < len(cc); i++ {
		gotos[i] = make(map[string]int)

		next := cfg.NonTerminals.Iterator()
		for j, ok := next(); ok; j, ok = next() {
			if gt := Goto(cc[i], j, cfg); len(gt) > 0 {
				index, yes := containsItemSet(cc, gt)
				if !yes {
					cc = append(cc, gt)
					gotos[i][j] = len(cc) - 1
				} else {
					gotos[i][j] = index
				}
			}
		}

		next = cfg.Terminals.Iterator()
		for j, ok := next(); ok; j, ok = next() {
			if j == "" {
				continue
			}

			if gt := Goto(cc[i], j, cfg); len(gt) > 0 {
				index, yes := containsItemSet(cc, gt)
				if !yes {
					cc = append(cc, gt)
					gotos[i][j] = len(cc) - 1
				} else {
					gotos[i][j] = index
				}
			}
		}
	}

	/* 	for i := 0; i < len(cc); i++ {
		var b strings.Builder
		b.WriteString("State: ")
		b.WriteString(fmt.Sprintf("%d\n", i))
		for _, v := range cc[i] {
			b.WriteString(v.rule.LHS())
			b.WriteString(" -> ")
			pointPlaced := false
			for i, k := range v.rule.RHS() {
				if i == v.pPos {
					b.WriteString(".")
					pointPlaced = true
				}
				b.WriteString(k)
				b.WriteString(" ")
			}

			if !pointPlaced {
				b.WriteString(".")
			}
			b.WriteString("\n")
		}

		fmt.Println(b.String())
	} */

	return
}

// GenSLR1Table generates the tables
func GenSLR1Table(cfg *grammar.CFG, followSets map[string]*containers.Set) (actionTable map[int]map[string]string, gotoTable map[int]map[string]int, err error) {
	cc, gotos := CanonicalCollection(cfg)
	actionTable = make(map[int]map[string]string)
	for i := 0; i < len(cc); i++ {
		actionTable[i] = make(map[string]string)
		for _, v := range cc[i] {
			if v.Equals(NewLR0item(len(cfg.StartRule.RHS()), cfg.StartRule)) {
				if actionTable[i]["$"] != "" {
					err = fmt.Errorf("The grammar is not SLR(1)")
					fmt.Println(actionTable[i]["$"])
					return
				}
				actionTable[i]["$"] = "ACCEPT"
			} else if v.pPos == len(v.rule.RHS()) {
				next := followSets[v.rule.LHS()].Iterator()
				for elem, ok := next(); ok; elem, ok = next() {
					if actionTable[i][elem] != "" {
						err = fmt.Errorf("The grammar is not SLR(1)")
						fmt.Println(actionTable[i][elem])
						return
					}
					actionTable[i][elem] = "REDUCE " + v.rule.String()
				}
			} else if cfg.Terminals.Contains(v.rule.RHS()[v.pPos]) /* && v.rule.RHS()[v.pPos] != "" */ {
				aux := fmt.Sprintf("SHIFT %d", gotos[i][v.rule.RHS()[v.pPos]])
				if actionTable[i][v.rule.RHS()[v.pPos]] != "" {
					if actionTable[i][v.rule.RHS()[v.pPos]] != aux {
						err = fmt.Errorf("The grammar is not SLR(1)")
						// fmt.Printf("State: %d, Terminal: %s, There Is: %s", i, v.rule.RHS()[v.pPos], actionTable[i][v.rule.RHS()[v.pPos]])
						// fmt.Println(actionTable[i][v.rule.RHS()[v.pPos]])
						return
					}
				} else {
					actionTable[i][v.rule.RHS()[v.pPos]] = aux
				}

			}
		}
	}

	gotoTable = make(map[int]map[string]int)
	for i := 0; i < len(cc); i++ {
		gotoTable[i] = make(map[string]int)
		next := cfg.NonTerminals.Iterator()
		for elem, ok := next(); ok; elem, ok = next() {
			if elem != cfg.StartRule.LHS() {
				gotoTable[i][elem] = gotos[i][elem]
			}
		}
	}

	return
}

func containsItem(i []*LR0item, item *LR0item) (yes bool) {
	for _, v := range i {
		if item.Equals(v) {
			yes = true
			return
		}
	}
	return
}

// brute-force (for now)
func containsItemSet(cc [][]*LR0item, i []*LR0item) (index int, yes bool) {
	yes = false
	for ind, v := range cc {
		for k := 0; k < len(i) && containsItem(v, i[k]); k++ {
			if k == len(i)-1 {
				index = ind
				yes = true
				return
			}
		}
	}
	return
}
