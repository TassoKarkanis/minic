package main

import (
	"fmt"
	"io"

	"github.com/TassoKarkanis/minic/types"
)

type GetText interface {
	GetText() string
}

type Codegen struct {
	out        io.Writer
	integerReg []*Register // rax is first
	values     map[GetText]*Value
}

func NewCodegen(out io.Writer) *Codegen {
	return &Codegen{
		out: out,
		integerReg: []*Register{
			&Register{
				integer:   true,
				fullName:  "rax",
				dwordName: "eax",
				wordName:  "ax",
				byteName:  "al",
			},
			&Register{
				integer:   true,
				fullName:  "rdi",
				dwordName: "edi",
				wordName:  "di",
				byteName:  "dl",
			},
			&Register{
				integer:   true,
				fullName:  "rsi",
				dwordName: "esx",
				wordName:  "si",
				byteName:  "sl",
			},
			&Register{
				integer:   true,
				fullName:  "rdx",
				dwordName: "edx",
				wordName:  "dx",
				byteName:  "dl",
			},
			&Register{
				integer:   true,
				fullName:  "rcx",
				dwordName: "ecx",
				wordName:  "cx",
				byteName:  "cl",
			},
			&Register{
				integer:   true,
				fullName:  "r8",
				dwordName: "r8d",
				wordName:  "r8w",
				byteName:  "r8b",
			},
			&Register{
				integer:   true,
				fullName:  "r9",
				dwordName: "r9d",
				wordName:  "r9w",
				byteName:  "r9b",
			},
		},
		values: make(map[GetText]*Value),
	}
}

func (c *Codegen) Close() {
	if len(c.values) > 0 {
		var unreleasedVal GetText
		for key, _ := range c.values {
			unreleasedVal = key
			break
		}
		msg := fmt.Sprintf("unreleased value: %s", unreleasedVal)
		panic(msg)
	}
}

func (c *Codegen) CreateIntValue(key GetText, typ types.Type, value string) *Value {
	val := &Value{
		typ:      typ,
		storage:  ConstantStorage,
		constant: value,
	}

	c.setValue(key, val)
	return val
}

func (c *Codegen) GetValue(key GetText) *Value {
	val := c.values[key]
	if val == nil {
		panic("value not found!")
	}
	return val
}

func (c *Codegen) MoveValue(destKey, srcKey GetText) {
	if c.values[destKey] != nil {
		c.fail("MoveValue(): destkey(%s) already exists", destKey.GetText())
	}

	if c.values[srcKey] == nil {
		c.fail("MoveValue(): srckey(%s) not found", destKey.GetText())
	}

	c.values[destKey] = c.values[srcKey]
	delete(c.values, srcKey)
}

func (c *Codegen) ReleaseValue(key GetText) {
	if c.values[key] == nil {
		msg := fmt.Sprintf("value to release not found: %s", key.GetText)
		panic(msg)
	}

	// clear the register binding, if any
	val := c.values[key]
	if val.register != nil {
		reg := val.register
		val.register = nil
		reg.binding = nil
	}

	delete(c.values, key)
}

func (c *Codegen) Add(key GetText, v1, v2 *Value) {
	reg := c.allocateIntRegister()
	val := &Value{
		typ:      v1.typ,
		storage:  RegisterStorage,
		register: reg,
	}
	reg.binding = val

	fmt.Fprintf(c.out, "\tmov %s, %s\n", val.Source(), v1.Source())
	fmt.Fprintf(c.out, "\tadd %s, %s\n", val.Source(), v2.Source())

	c.setValue(key, val)
}

func (c *Codegen) Subtract(key GetText, v1, v2 *Value) {
	reg := c.allocateIntRegister()
	val := &Value{
		typ:      v1.typ,
		storage:  RegisterStorage,
		register: reg,
	}
	reg.binding = val

	fmt.Fprintf(c.out, "\tmov %s, %s\n", val.Source(), v1.Source())
	fmt.Fprintf(c.out, "\tsub %s, %s\n", val.Source(), v2.Source())

	c.setValue(key, val)
}

// Return a value from a function
func (c *Codegen) ReturnValue(key GetText) {
	val := c.GetValue(key)
	fmt.Fprintf(c.out, "\tmov %s, %s\n", "eax", val.Source())
	fmt.Fprintf(c.out, "\tret\n")
}

// Return void from a function
func (c *Codegen) Return() {
	fmt.Fprintf(c.out, "\tret\n")
}

func (c *Codegen) setValue(key GetText, val *Value) {
	if c.values[key] != nil {
		msg := fmt.Sprintf("value already set for %s", key.GetText())
		panic(msg)
	}

	c.values[key] = val
}

func (c *Codegen) allocateIntRegister() *Register {
	for i := len(c.integerReg) - 1; i > 0; i-- {
		reg := c.integerReg[i]
		if reg.binding == nil {
			return reg
		}
	}

	c.fail("failed to allocate integer register!")
	return nil
}

func (c *Codegen) fail(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	panic(msg)
}
