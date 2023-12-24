package main

import (
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestEmptyFunction(t *testing.T) {
	runTest(t, "empty-function")
}

func TestReturnOne(t *testing.T) {
	runTest(t, "return-one")
}

func TestAdd(t *testing.T) {
	runTest(t, "add")
}

func TestReturnArg(t *testing.T) {
	runTest(t, "return-arg")
}

func TestTypes(t *testing.T) {
	runTest(t, "types")
}

func runTest(t *testing.T, testName string) {
	r := require.New(t)
	log.SetLevel(log.DebugLevel)

	// make the output filename
	dir := "testdata"
	inputFile := filepath.Join(dir, testName+".c")
	outputFile := filepath.Join(dir, testName+".out.asm")

	// call the compiler
	opts := Options{
		OutputFile:  outputFile,
		CompileOnly: true,
		CodegenOnly: true,
	}
	err := mainE(opts, []string{inputFile})
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

}
