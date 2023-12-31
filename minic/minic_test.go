package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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

func TestAddMany(t *testing.T) {
	runTest(t, "add-many")
}

func TestReturnArg(t *testing.T) {
	runTest(t, "return-arg")
}

func TestSub(t *testing.T) {
	runTest(t, "sub")
}

func TestMult(t *testing.T) {
	runTest(t, "mult")
}

func TestDiv(t *testing.T) {
	runTest(t, "div")
}

func TestVar(t *testing.T) {
	runTest(t, "var")
}

func TestVarScope(t *testing.T) {
	runTest(t, "var-scope")
}

func TestUnaryMinus(t *testing.T) {
	runTest(t, "unary-minus")
}

func runTest(t *testing.T, name string) {
	r := testRunner{
		name:       name,
		miniC:      NewMiniC(),
		Assertions: require.New(t),
	}

	r.run()
}

type FileType int

const (
	FileDriverC      FileType = iota // (stored) C file of test driver
	FileDriverObject                 // (generated) C file of test driver
	FileExpectedAsm                  // (stored) expected assembly output
	FileInputC                       // (stored) test input C file
	FileOutputAsm                    // (generated) actual assembly output
	FileOutputObject                 // (generated) assembled output file
	FileOutputBin                    // (generated) test binary
)

var filePostFix = map[FileType]string{
	FileDriverC:      ".driver.c",
	FileDriverObject: ".driver.o",
	FileExpectedAsm:  ".asm",
	FileInputC:       ".c",
	FileOutputAsm:    ".out.asm",
	FileOutputObject: ".o",
	FileOutputBin:    ".out",
}

type testRunner struct {
	name  string
	miniC *MiniC
	*require.Assertions
}

func (r *testRunner) run() {
	log.SetLevel(log.DebugLevel)

	r.compile()
	r.compareToExpected()
	r.assemble()
	r.buildBin()
	r.runTest()
}

func (r *testRunner) compile() {
	// compile to assembly
	err := r.miniC.CompileFile(r.getFile(FileInputC), r.getFile(FileOutputAsm))
	r.NoError(err)
}

func (r *testRunner) compareToExpected() {
	// read the expected file
	expectedFile := r.getFile(FileExpectedAsm)
	expected, err := os.ReadFile(expectedFile)
	r.NoError(err)

	// compare the generated output with the expected
	actualFile := r.getFile(FileOutputAsm)
	actual, err := os.ReadFile(actualFile)
	r.NoError(err)

	if string(actual) != string(expected) {
		// diff the two files
		cmdLine := []string{"diff", "-c", expectedFile, actualFile}
		output, _ := r.runCommand(cmdLine)
		fmt.Println(output)
		r.Fail("result not as expected")
	}
}

func (r *testRunner) assemble() {
	err := r.miniC.AssembleFile(r.getFile(FileOutputAsm), r.getFile(FileOutputObject))
	r.NoError(err)
}

func (r *testRunner) buildBin() {
	// compile the driver with gcc
	gccCmd := []string{
		"gcc",
		"-c",
		r.getFile(FileDriverC),
		"-o",
		r.getFile(FileDriverObject),
	}
	cmd := exec.Command(gccCmd[0], gccCmd[1:]...)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		fmt.Print(stdout.String())
		fmt.Print(stderr.String())
	}
	r.NoError(err)

	// link with minic
	modules := []string{
		r.getFile(FileDriverObject),
		r.getFile(FileOutputObject),
	}
	err = r.miniC.Link(modules, r.getFile(FileOutputBin))
	r.NoError(err)
}

func (r *testRunner) runTest() {
	cmdLine := []string{r.getFile(FileOutputBin)}
	output, err := r.runCommand(cmdLine)
	fmt.Print(output)
	r.NoError(err)
}

func (r *testRunner) runCommand(cmdLine []string) (string, error) {
	cmd := exec.Command(cmdLine[0], cmdLine[1:]...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	output := stdout.String() + stderr.String()
	return output, err
}

func (r *testRunner) getFile(fileType FileType) string {
	dir := "testdata"
	return filepath.Join(dir, r.name+filePostFix[fileType])
}
