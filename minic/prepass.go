package main

import (
	"fmt"
	"io"

	"github.com/TassoKarkanis/minic/parser"
)

type Prepass struct {
	*parser.BaseCListener
	Output                   io.Writer
	LastDeclaratorIdentifier string
	GlobalSymbols            []string
}

func NewPrepass(output io.Writer) *Prepass {
	p := &Prepass{
		Output: output,
	}

	// add the data section preamble
	fmt.Fprintf(p.Output, "section .data\n")
	return p
}

func (p *Prepass) ExitCompilationUnit(c *parser.CompilationUnitContext) {
	// close the data section
	fmt.Fprintf(p.Output, "\n")

	// output the text section preamble
	fmt.Fprintf(p.Output, "section .text\n")
	for _, name := range p.GlobalSymbols {
		fmt.Fprintf(p.Output, "\tglobal %s\n", name)
	}
}

func (p *Prepass) EnterDirectDeclaratorIdentifier(c *parser.DirectDeclaratorIdentifierContext) {
	p.LastDeclaratorIdentifier = c.Identifier().GetText()
}

func (p *Prepass) EnterParameterTypeList(c *parser.ParameterTypeListContext) {
	functionName := p.LastDeclaratorIdentifier
	p.GlobalSymbols = append(p.GlobalSymbols, functionName)
}
