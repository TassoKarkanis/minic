package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	OutputFile  string
	CompileOnly bool
	CodegenOnly bool
}

func mainE(options Options, modules []string) error {
	m := NewMiniC()

	for _, module := range modules {
		err := m.CompileFile(module, options.OutputFile)
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
