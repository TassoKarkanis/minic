package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMinic(t *testing.T) {
	r := require.New(t)
	log.SetLevel(log.DebugLevel)

	// find the paths of all input C files
	paths, err := filepath.Glob(filepath.Join("testdata", "*.c"))
	r.NoError(err)

	// generate code for each file
	for _, path := range paths {
		dir, inputFilename := filepath.Split(path)
		testName := inputFilename[:len(inputFilename)-len(filepath.Ext(path))]

		fmt.Printf("testName: %s\n", testName)
		if testName != "return-arg" {
			continue
		}

		t.Run(testName, func(t *testing.T) {
			// make the output filename
			outputFile := filepath.Join(dir, testName+".out.asm")

			// call the compiler
			opts := Options{
				OutputFile:  outputFile,
				CompileOnly: true,
				CodegenOnly: true,
			}
			err := mainE(opts, []string{path})
			r.NoError(err)

			// read the expected file
			expectedFilename := filepath.Join(dir, testName+".asm")
			expected, err := os.ReadFile(expectedFilename)
			r.NoError(err)

			// compare the generated output with the expected
			actual, err := os.ReadFile(outputFile)
			r.NoError(err)

			if string(actual) != string(expected) {
				r.Fail("result not as expected")
			}
		})
	}
}
