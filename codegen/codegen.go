package codegen

import (
	"fmt"
	"io"

	"github.com/TassoKarkanis/minic/types"
)

type GetText interface {
	GetText() string
}

type Scope struct {
	offset int // negative, offset of current end of scope
}

const RAX = "rax"
const RDI = "rdi"
const RSI = "rsi"
const RDX = "rdx"
const RCX = "rcx"
const R8 = "r8"
const R9 = "r9"

var registerAllocationOrder = [...]string{
	RDI, RSI, RDX, RCX, R8, R9,
}

type Codegen struct {
	out        io.Writer
	integerReg map[string]*Register
	values     map[GetText]*Value
	funcName   string // name of current function
	scopes     []*Scope
}

func NewCodegen(out io.Writer) *Codegen {
	// make the map of registers
	registers := map[string]*Register{
		RAX: {
			integer:   true,
			fullName:  "rax",
			dwordName: "eax",
			wordName:  "ax",
			byteName:  "al",
		},
		RDI: {
			integer:   true,
			fullName:  "rdi",
			dwordName: "edi",
			wordName:  "di",
			byteName:  "dl",
		},
		RSI: {
			integer:   true,
			fullName:  "rsi",
			dwordName: "esx",
			wordName:  "si",
			byteName:  "sl",
		},
		RDX: {
			integer:   true,
			fullName:  "rdx",
			dwordName: "edx",
			wordName:  "dx",
			byteName:  "dl",
		},
		RCX: {
			integer:   true,
			fullName:  "rcx",
			dwordName: "ecx",
			wordName:  "cx",
			byteName:  "cl",
		},
		R8: {
			integer:   true,
			fullName:  "r8",
			dwordName: "r8d",
			wordName:  "r8w",
			byteName:  "r8b",
		},
		R9: {
			integer:   true,
			fullName:  "r9",
			dwordName: "r9d",
			wordName:  "r9w",
			byteName:  "r9b",
		},
	}

	c := &Codegen{
		out:        out,
		values:     make(map[GetText]*Value),
		integerReg: registers,
	}

	return c
}

func (c *Codegen) Close() {
	if len(c.values) > 0 {
		var unreleasedVal GetText
		for key := range c.values {
			unreleasedVal = key
			break
		}
		c.fail("unreleased value: %s", unreleasedVal)
	}

	if len(c.scopes) != 0 {
		c.fail("must close scopes and stack frame")
	}
}

func (c *Codegen) StartStackFrame(name string, params []types.Type) []*Value {
	// we don't handle nested functions
	if len(c.scopes) > 0 {
		c.fail("cannot handle nested scopes!")
	}

	// we can only handle parameters in registers for now
	if len(params) > 6 {
		c.fail("cannot handle more than 6 parameters!")
	}

	// store the name
	c.funcName = name
	fmt.Fprintf(c.out, "%s:\n", c.funcName)

	// create the curent scope
	scope := &Scope{}
	c.scopes = append(c.scopes, scope)

	// allocate values for the parameters
	var values []*Value
	for i, typ := range params {
		offset := scope.offset - 4 // TODO: based on type
		scope.offset = offset

		// create the value
		val := &Value{
			typ:      typ,
			storage:  LocalStorage,
			offset:   offset,
			register: c.integerReg[registerAllocationOrder[i]],
			dirty:    true,
		}
		values = append(values, val)
	}

	return values
}

func (c *Codegen) EndStackFrame() {
	// scopes should be terminated
	if len(c.scopes) != 1 {
		c.fail("ending stack frame with open scope")
	}

	// write epilogue
	fmt.Fprintf(c.out, "%s:\n", c.endFrameLabel())
	fmt.Fprintf(c.out, "\tret\n\n")

	c.scopes = nil
}

func (c *Codegen) StartScope() {
	if len(c.scopes) == 0 {
		c.fail("no stack frame started")
	}

	// push a scope
	scope := &Scope{
		offset: c.scopes[len(c.scopes)-1].offset,
	}
	c.scopes = append(c.scopes, scope)
}

func (c *Codegen) EndScope() {
	if len(c.scopes) == 1 {
		c.fail("cannot end last scope")
	}

	c.scopes = c.scopes[0 : len(c.scopes)-1]
}

func (c *Codegen) CreateIntLiteralValue(key GetText, typ types.Type, value string) *Value {
	val := &Value{
		typ:     typ,
		storage: ConstantStorage,
		value:   value,
	}

	c.setValue(key, val)
	return val
}

func (c *Codegen) CreateValue(key GetText, val *Value) {
	if c.values[key] != nil {
		c.fail("CreateValue(): key(%s) already exists", key.GetText())
	}

	c.values[key] = val
}

func (c *Codegen) GetValue(key GetText) *Value {
	val := c.values[key]
	if val == nil {
		c.fail("value not found!")
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
		c.fail("value to release not found: %s", key.GetText())
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
	val := c.allocateLocalStorage(v1.typ)

	// TODO: use proper register size
	fmt.Fprintf(c.out, "\tmov %s, %s\n", val.Source(), v1.Source())
	fmt.Fprintf(c.out, "\tadd %s, %s\n", val.Source(), v2.Source())

	c.setValue(key, val)
}

func (c *Codegen) Subtract(key GetText, v1, v2 *Value) {
	val := c.allocateLocalStorage(v1.typ)

	// TODO: use proper register size
	fmt.Fprintf(c.out, "\tmov %s, %s\n", val.Source(), v1.Source())
	fmt.Fprintf(c.out, "\tsub %s, %s\n", val.Source(), v2.Source())

	c.setValue(key, val)
}

// Return a value from a function
func (c *Codegen) ReturnValue(key GetText) {
	val := c.GetValue(key)
	fmt.Fprintf(c.out, "\tmov eax, %s\n", val.Source())
	c.Return()
}

// Return void from a function
func (c *Codegen) Return() {
	fmt.Fprintf(c.out, "\tjmp %s\n", c.endFrameLabel())
}

func (c *Codegen) setValue(key GetText, val *Value) {
	if c.values[key] != nil {
		c.fail("value already set for %s", key.GetText())
	}

	c.values[key] = val
}

func (c *Codegen) allocateRegister() *Register {
	for _, regName := range registerAllocationOrder {
		reg := c.integerReg[regName]
		if reg.binding == nil {
			return reg
		}
	}

	return nil
}

func (c *Codegen) allocateLocalStorage(typ types.Type) *Value {
	scope := c.scope()
	scope.offset -= 4 // TODO: proper type

	// allocate a register
	reg := c.allocateRegister()

	val := &Value{
		typ:      typ,
		storage:  LocalStorage,
		offset:   scope.offset,
		register: c.allocateRegister(),
	}

	// bind the register to the value
	val.register = reg
	reg.binding = val

	return val
}

func (c *Codegen) scope() *Scope {
	if len(c.scopes) == 0 {
		c.fail("no scopes!")
	}

	return c.scopes[len(c.scopes)-1]
}

func (c *Codegen) endFrameLabel() string {
	return fmt.Sprintf("%s.end", c.funcName)
}

func (c *Codegen) fail(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	panic(msg)
}
