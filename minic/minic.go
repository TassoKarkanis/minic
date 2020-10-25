package main

import (
	"fmt"
	"io"
	"os"

	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/symbols"
	"github.com/TassoKarkanis/minic/types"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/spf13/cobra"
)

type Codegen struct {
	*parser.BaseCListener
	Output                   io.Writer
	Symbols                  *symbols.Table
	LastDeclaratorIdentifier string
	LastType                 int // i.e. CParserVoid
	LastFunction             *types.FunctionType
}

func NewCodegen(output io.Writer) *Codegen {
	return &Codegen{
		Output:  output,
		Symbols: symbols.NewTable(),
	}
}

func (c *Codegen) EnterCompilationUnit(ctx *parser.CompilationUnitContext) {
	fmt.Printf("EnterCompilationUnit(): %s\n", ctx.GetText())
}

func (c *Codegen) ExitCompilationUnit(ctx *parser.CompilationUnitContext) {
	fmt.Printf("ExitCompilationUnit()\n")
}

func (c *Codegen) EnterExternalDeclaration(ctx *parser.ExternalDeclarationContext) {
	fmt.Printf("EnterExternalDeclaration(): %s\n", ctx.GetText())
}

func (c *Codegen) EnterFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	fmt.Printf("EnterFunctionDefinition(): %s\n", ctx.GetText())
	c.LastFunction = &types.FunctionType{}
	c.Symbols.PushScope()
}

func (c *Codegen) ExitFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	fmt.Printf("ExitFunctionDefinition(): %s\n", ctx.GetText())
	fmt.Fprintf(c.Output, "\tret\n\n")
	c.Symbols.PopScope()
	c.LastFunction = nil
}

func (c *Codegen) EnterDeclarator(ctx *parser.DeclaratorContext) {
	fmt.Printf("EnterDeclarator(): %s\n", ctx.GetText())
}

func (c *Codegen) EnterDirectDeclaratorFunction(ctx *parser.DirectDeclaratorFunctionContext) {
	fmt.Printf("EnterDirectDeclaratorFunction(): %s\n", ctx.GetText())
}

func (c *Codegen) ExitDirectDeclaratorFunction(ctx *parser.DirectDeclaratorFunctionContext) {
	fmt.Printf("ExitDirectDeclaratorFunction(): %s\n", ctx.GetText())
}

func (c *Codegen) EnterDirectDeclaratorIdentifier(ctx *parser.DirectDeclaratorIdentifierContext) {
	fmt.Printf("EnterDirectDeclaratorIdentifier(): %s\n", ctx.GetText())
	c.LastDeclaratorIdentifier = ctx.Identifier().GetText()
	fmt.Printf("LastDeclaratorIdentifier: %s\n", c.LastDeclaratorIdentifier)
}

func (c *Codegen) EnterTypeSpecifierSimple(ctx *parser.TypeSpecifierSimpleContext) {
	fmt.Printf("EnterTypeSpecifierSimple(): %s\n", ctx.GetText())
	fmt.Printf("  GetStart(): %s\n", ctx.GetStart().GetText())
	c.LastType = ctx.GetStart().GetTokenType()
	fmt.Printf("LastType: %v\n", c.LastType)
}

func (c *Codegen) EnterParameterTypeList(ctx *parser.ParameterTypeListContext) {
	fmt.Printf("EnterParameterTypeList(): %s\n", ctx.GetText())
	cp := ctx.GetParser().(*parser.CParser)
	name := c.LastDeclaratorIdentifier
	c.LastFunction.ReturnType = types.NewBasicType(c.LastType, cp)
	c.LastFunction.Name = name
	c.Symbols.AddSymbol(name, c.LastFunction)
	fmt.Fprintf(c.Output, "%s:\n", name)
}

func (c *Codegen) EnterCompoundStatement(ctx *parser.CompoundStatementContext) {
	fmt.Printf("EnterCompoundStatement(): %s\n", ctx.GetText())
}

type Options struct {
	OutputFile  string
	CompileOnly bool
	CodegenOnly bool
}

func runPrepass(inputFile string, output io.Writer) error {
	// setup the input
	is, err := antlr.NewFileStream(inputFile)
	if err != nil {
		return err
	}

	// Create the Lexer
	lexer := parser.NewCLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewCParser(stream)

	// Finally parse the expression
	prepass := NewPrepass(output)
	antlr.ParseTreeWalkerDefault.Walk(prepass, p.CompilationUnit())
	return nil
}

func runCodegen(inputFile string, output io.Writer) error {
	// setup the input
	is, err := antlr.NewFileStream(inputFile)
	if err != nil {
		return err
	}

	// Create the Lexer
	lexer := parser.NewCLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewCParser(stream)

	// Finally parse the expression
	codegen := NewCodegen(output)
	antlr.ParseTreeWalkerDefault.Walk(codegen, p.CompilationUnit())
	return nil
}

func mainE(options Options, modules []string) error {
	// open the output file
	output, err := os.Create(options.OutputFile)
	if err != nil {
		return err
	}

	for _, module := range modules {
		err = runPrepass(module, output)
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "\n")

		err = runCodegen(module, output)
		if err != nil {
			return err
		}
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

	flags.BoolVarP(
		&options.CompileOnly,
		"compile-only",
		"c",
		false,
		"compile only and don't run linker")
	flags.StringVarP(
		&options.OutputFile,
		"output",
		"o",
		"a.out",
		"output to this file")
	flags.BoolVarP(
		&options.CodegenOnly,
		"codegen-only",
		"S",
		false,
		"generate code only and don't assemble or run linker")

	err := cobraCmd.Execute()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
