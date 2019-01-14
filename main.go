package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mbotezatu/slr1-table/grammar"
	"github.com/mbotezatu/slr1-table/slr1"
	"github.com/olekukonko/tablewriter"
)

func main() {

	// var nonTerm = regexp.MustCompile(`[A-Z]{1}[a-zA-Z0-9]*`)
	// var term = regexp.MustCompile(`[a-z]{1}[a-zA-Z0-9]*`)
	var imp = regexp.MustCompile(`->`)
	var prod = regexp.MustCompile(`^((\s*?[A-Z]{1}[a-zA-Z0-9]*\s*?->(\s*?|[A-Z]{1}[a-zA-Z0-9]*|[a-z]{1}[a-zA-Z0-9]*)*?)|(\s*))$`)

	var file *os.File

	if len(os.Args) > 1 {
		var err error
		file, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}()
	} else {
		file = os.Stdin
	}

	s := bufio.NewScanner(file)

	cfg := grammar.NewCFG()
	stopAssign := false
	firstRule := true

	for s.Scan() {
		if prod.MatchString(s.Text()) {
			if stopAssign {
				continue
			}

			lrhs := imp.Split(s.Text(), -1)
			lrhs[0] = strings.TrimSpace(lrhs[0])
			if len(lrhs) < 2 {
				continue
			}

			rhsSymbols := strings.Fields(lrhs[1])
			if len(rhsSymbols) == 0 {
				rhsSymbols = append(rhsSymbols, "")
			}

			rule := grammar.NewProduction(lrhs[0], rhsSymbols)
			cfg.Rules[lrhs[0]] = append(cfg.Rules[lrhs[0]], rule)
			if firstRule {
				cfg.StartRule = rule
				firstRule = false
			}
			cfg.NonTerminals.Add(lrhs[0])
			for _, v := range rule.RHS() {
				if grammar.IsTerminal(v) {
					cfg.Terminals.Add(v)
				} else {
					cfg.NonTerminals.Add(v)
				}
			}

		} else {
			fmt.Fprintf(os.Stderr, "%s does not match the production regexp.\n", s.Text())
			stopAssign = true
		}
	}

	if stopAssign {
		os.Exit(1)
	}

	if cfg.StartRule == nil {
		fmt.Fprint(os.Stderr, "No grammar defined.\n")
		os.Exit(1)
	}

	firstSets, err := grammar.First(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("First Sets:")
	for k, v := range firstSets {
		if k == "" {
			k = "λ"
		}
		fmt.Printf("\t%s: { ", k)
		next := v.Iterator()
		for elem, ok := next(); ok; elem, ok = next() {
			if elem == "" {
				elem = "λ"
			}
			fmt.Printf("%s ", elem)
		}
		fmt.Println("}")
	}

	followSets := grammar.Follow(cfg, firstSets)
	fmt.Println("Follow Sets: ")
	for k, v := range followSets {
		fmt.Printf("\t%s: { ", k)
		next := v.Iterator()
		for elem, ok := next(); ok; elem, ok = next() {
			fmt.Printf("%s ", elem)
		}
		fmt.Println("}")
	}

	fmt.Println()

	actionTable, gotoTable, err := slr1.GenSLR1Table(cfg, followSets)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	actionTableWriter := tablewriter.NewWriter(os.Stdout)
	var actionTableHeader []string
	next := cfg.Terminals.Iterator()
	actionTableHeader = append(actionTableHeader, "State")
	for elem, ok := next(); ok; elem, ok = next() {
		if elem != "" {
			actionTableHeader = append(actionTableHeader, elem)
		}
	}
	actionTableHeader = append(actionTableHeader, "$")
	actionTableWriter.Append(actionTableHeader)
	actionTableWriter.SetRowLine(true)
	prtActionTable := make([][]string, len(actionTable))
	for i := 0; i < len(actionTable); i++ {
		prtActionTable[i] = make([]string, len(actionTableHeader))
		prtActionTable[i][0] = fmt.Sprintf("%d", i)
		for j := 1; j < len(actionTableHeader); j++ {
			prtActionTable[i][j] = actionTable[i][actionTableHeader[j]]
		}
	}
	actionTableWriter.AppendBulk(prtActionTable)
	fmt.Println("ACTION Table:")
	actionTableWriter.Render()
	fmt.Println()

	gotoTableWriter := tablewriter.NewWriter(os.Stdout)
	gotoTableHeader := make([]string, cfg.NonTerminals.Len()+1)
	next = cfg.NonTerminals.Iterator()
	gotoTableHeader[0] = "State"
	i := 1
	for elem, ok := next(); ok; elem, ok = next() {
		gotoTableHeader[i] = elem
		i++
	}
	gotoTableWriter.Append(gotoTableHeader)
	gotoTableWriter.SetRowLine(true)
	prtGotoTable := make([][]string, len(gotoTable))
	for i := 0; i < len(gotoTable); i++ {
		prtGotoTable[i] = make([]string, len(gotoTableHeader))
		prtGotoTable[i][0] = fmt.Sprintf("%d", i)
		for j := 1; j < len(gotoTableHeader); j++ {
			var state string
			if gotoTable[i][gotoTableHeader[j]] > 0 {
				state = fmt.Sprintf("%d", gotoTable[i][gotoTableHeader[j]])
			}
			prtGotoTable[i][j] = state
		}
	}
	gotoTableWriter.AppendBulk(prtGotoTable)
	fmt.Println("GOTO Table:")
	gotoTableWriter.Render()
}
