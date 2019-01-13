package grammar

import (
	"bytes"

	"github.com/mbotezatu/slr1-table/containers"
)

// CFG defines a context-free grammar
type CFG struct {
	Terminals    *containers.Set
	NonTerminals *containers.Set
	StartRule    *Production
	Rules        map[string][]*Production
}

// NewCFG creates and returns a Context-Free Grammar pointer
func NewCFG() (cfg *CFG) {
	cfg = &CFG{
		containers.NewSet(),
		containers.NewSet(),
		nil,
		make(map[string][]*Production)}
	return
}

func (cfg *CFG) String() string {
	var b bytes.Buffer

	b.WriteString("Terminals: ")
	b.WriteString(cfg.Terminals.String())
	b.WriteString("\n")
	b.WriteString("Non-Terminals: ")
	b.WriteString(cfg.NonTerminals.String())
	b.WriteString("\n\nProductions: \n\t")
	for _, p := range cfg.Rules {
		for _, v := range p {
			b.WriteString(v.String())
			b.WriteString("\n\t")
		}
	}

	return b.String()
}
