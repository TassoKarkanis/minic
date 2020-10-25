package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/symbols"
	"github.com/TassoKarkanis/minic/types"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/spf13/cobra"
)

type FunctionInfo struct {
}

type Listener struct {
	*parser.BaseCListener
	Symbols                  *symbols.Table
	LastDeclaratorIdentifier string
	LastType                 int // i.e. CParserVoid
	LastFunction             *types.FunctionType
}

func NewListener() *Listener {
	return &Listener{
		Symbols: symbols.NewTable(),
	}
}

func (l *Listener) EnterCompilationUnit(c *parser.CompilationUnitContext) {
	fmt.Printf("EnterCompilationUnit(): %s\n", c.GetText())
}

func (l *Listener) ExitCompilationUnit(c *parser.CompilationUnitContext) {
	fmt.Printf("ExitCompilationUnit()\n")
}

func (l *Listener) EnterExternalDeclaration(c *parser.ExternalDeclarationContext) {
	fmt.Printf("EnterExternalDeclaration(): %s\n", c.GetText())
}

func (l *Listener) EnterFunctionDefinition(c *parser.FunctionDefinitionContext) {
	fmt.Printf("EnterFunctionDefinition(): %s\n", c.GetText())
	l.LastFunction = &types.FunctionType{}
	l.Symbols.PushScope()
}

func (l *Listener) ExitFunctionDefinition(c *parser.FunctionDefinitionContext) {
	fmt.Printf("ExitFunctionDefinition(): %s\n", c.GetText())
	l.Symbols.PopScope()
	l.LastFunction = nil
}

func (l *Listener) EnterDeclarator(c *parser.DeclaratorContext) {
	fmt.Printf("EnterDeclarator(): %s\n", c.GetText())
}

func (l *Listener) EnterDirectDeclaratorFunction(c *parser.DirectDeclaratorFunctionContext) {
	fmt.Printf("EnterDirectDeclaratorFunction(): %s\n", c.GetText())
	l.Symbols.AddSymbol(l.LastFunction.Name, l.LastFunction)
}

func (l *Listener) ExitDirectDeclaratorFunction(c *parser.DirectDeclaratorFunctionContext) {
	fmt.Printf("ExitDirectDeclaratorFunction(): %s\n", c.GetText())
}

func (l *Listener) EnterDirectDeclaratorIdentifier(c *parser.DirectDeclaratorIdentifierContext) {
	fmt.Printf("EnterDirectDeclaratorIdentifier(): %s\n", c.GetText())
	l.LastDeclaratorIdentifier = c.Identifier().GetText()
	fmt.Printf("LastDeclaratorIdentifier: %s\n", l.LastDeclaratorIdentifier)
}

func (l *Listener) EnterTypeSpecifierSimple(c *parser.TypeSpecifierSimpleContext) {
	fmt.Printf("EnterTypeSpecifierSimple(): %s\n", c.GetText())
	fmt.Printf("  GetStart(): %s\n", c.GetStart().GetText())
	l.LastType = c.GetStart().GetTokenType()
	fmt.Printf("LastType: %v\n", l.LastType)
}

func (l *Listener) EnterParameterTypeList(c *parser.ParameterTypeListContext) {
	fmt.Printf("EnterParameterTypeList(): %s\n", c.GetText())
	cp := c.GetParser().(*parser.CParser)
	l.LastFunction.ReturnType = types.NewBasicType(l.LastType, cp)
	l.LastFunction.Name = l.LastDeclaratorIdentifier
}

func (l *Listener) EnterCompoundStatement(c *parser.CompoundStatementContext) {
	fmt.Printf("EnterCompoundStatement(): %s\n", c.GetText())
}

type Options struct {
	ProgramFile string
}

func mainE(options Options, modules []string) error {
	for _, module := range modules {
		// read the module
		data, err := ioutil.ReadFile(module)
		if err != nil {
			return err
		}

		// setup the input
		is := antlr.NewInputStream(string(data))

		// Create the Lexer
		lexer := parser.NewCLexer(is)
		stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

		// Create the Parser
		p := parser.NewCParser(stream)

		// Finally parse the expression
		listener := NewListener()
		antlr.ParseTreeWalkerDefault.Walk(listener, p.CompilationUnit())
	}

	return nil
}

func main() {
	options := Options{}

	cobraCmd := &cobra.Command{
		Use:  "minic module.c [module1.c ...] -o program",
		Args: cobra.MinimumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return mainE(options, args)
		},
	}

	flags := cobraCmd.Flags()
	flags.StringVarP(&options.ProgramFile, "program", "o", "a.out", "output file name")

	err := cobraCmd.Execute()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
