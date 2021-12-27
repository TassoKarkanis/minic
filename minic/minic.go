package main

import (
	"fmt"
	"io"
	"os"

	"github.com/TassoKarkanis/minic/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/spf13/cobra"
)

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

func runMainPass(inputFile string, output io.Writer) error {
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
	codegen := NewMainPass(output)
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

		err = runMainPass(module, output)
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
