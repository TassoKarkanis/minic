package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/TassoKarkanis/minic/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type MiniC struct {
}

func NewMiniC() *MiniC {
	return &MiniC{}
}

// Compile a .c file to .asm
func (m *MiniC) CompileFile(inputFile string, outputFile string) (err error) {
	// open the output file
	output, err := os.Create(outputFile)
	if err != nil {
		return
	}

	// recover panics and return the error
	defer func() {
		obj := recover()
		if obj != nil {
			err = fmt.Errorf("%s", obj.(string))
			output.Close()
		}
	}()

	// run the prepass
	err = m.runPrepass(inputFile, output)
	if err != nil {
		output.Close()
		return
	}

	fmt.Fprintf(output, "\n")

	// run the main pass
	err = m.runMainPass(inputFile, output)
	if err != nil {
		output.Close()
		return
	}

	// close the file
	err = output.Close()
	if err != nil {
		return
	}

	return
}

func (m *MiniC) AssembleFile(inputFile string, outputFile string) error {
	cmd := exec.Command("nasm", "-f", "elf64", "-o", outputFile, inputFile)
	err := cmd.Run()
	return err
}

func (m *MiniC) Link(inputFiles []string, output string) error {
	// link with gcc
	args := []string{"-o", output}
	args = append(args, inputFiles...)
	cmd := exec.Command("gcc", args...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()

	// output stdout and stderr
	fmt.Print(stdout.String())
	fmt.Print(stderr.String())

	return err
}

func (m *MiniC) runPrepass(inputFile string, output io.Writer) error {
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

func (m *MiniC) runMainPass(inputFile string, output io.Writer) error {
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
