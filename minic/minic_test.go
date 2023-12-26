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

func TestReturnArg(t *testing.T) {
	runTest(t, "return-arg")
}

func TestTypes(t *testing.T) {
	runTest(t, "types")
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
	expected, err := os.ReadFile(r.getFile(FileExpectedAsm))
	r.NoError(err)

	// compare the generated output with the expected
	actual, err := os.ReadFile(r.getFile(FileOutputAsm))
	r.NoError(err)

	if string(actual) != string(expected) {
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
	cmd := exec.Command(r.getFile(FileOutputBin))
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	fmt.Print(stdout.String())
	fmt.Print(stderr.String())
	r.NoError(err)
}

func (r *testRunner) getFile(fileType FileType) string {
	dir := "testdata"
	return filepath.Join(dir, r.name+filePostFix[fileType])
}