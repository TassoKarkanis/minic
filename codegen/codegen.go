package codegen

import (
	"fmt"
	"io"

	"github.com/TassoKarkanis/minic/types"
)

type GetText interface {
	GetText() string
}

// Used to generate code for each C function.
type Codegen struct {
	out        io.Writer
	integerReg map[string]*Register
	values     map[GetText]*Value
	funcName   string // name of current function
	stack      *Stack // tracks allocate stack space
}

// NewCodegen returns a new generate with a writer for the assembly code.
func NewCodegen(out io.Writer) *Codegen {

	c := &Codegen{
		out:        out,
		integerReg: AllocateIntegerRegisters(),
		values:     make(map[GetText]*Value),
		stack:      NewStack(),
	}

	return c
}

// Close checks that all values have been released.
func (c *Codegen) Close() {
	if len(c.values) > 0 {
		var unreleasedVal GetText
		for key := range c.values {
			unreleasedVal = key
			break
		}
		c.fail("unreleased value: %s", unreleasedVal)
	}
}

// StartStackFrame writes the function program and creates values for the parameters.
func (c *Codegen) StartStackFrame(name string, paramNames []string, paramTypes []types.Type) []*Value {
	// we can only handle parameters in registers for now
	if len(paramNames) > 6 {
		c.fail("cannot handle more than 6 parameters!")
	}

	// store the name
	c.funcName = name
	fmt.Fprintf(c.out, "%s:\n", c.funcName)

	// allocate values for the parameters
	var values []*Value
	for i, typ := range paramTypes {
		offset := c.stack.Alloc4() // TODO: based on type

		// create the value
		val := NewLocalValue(typ, offset)
		c.bind(val, c.integerReg[integerParameterOrder[i]], true)
		values = append(values, val)

		fmt.Fprintf(c.out, "\t; param %s %s -> %s\n", typ.String(), paramNames[i], val.register.Name(4)) // TODO
	}

	// store the callee-save register
	for _, name := range integerRegisterAllocationOrder {
		reg := c.integerReg[name]
		if !reg.callerSave {
			fmt.Fprintf(c.out, "\tpush %s\n", reg.fullName)
		}
	}

	return values
}

// EndStackFrame writes the function epilogue.
func (c *Codegen) EndStackFrame() {
	// write epilogue
	fmt.Fprintf(c.out, "%s:\n", c.endFrameLabel())

	// restore the caller-save registers
	for i := len(integerRegisterAllocationOrder) - 1; i >= 0; i-- {
		name := integerRegisterAllocationOrder[i]
		reg := c.integerReg[name]
		if !reg.callerSave {
			fmt.Fprintf(c.out, "\tpop %s\n", reg.fullName)
		}
	}

	fmt.Fprintf(c.out, "\tret\n\n")
}

// CreateIntLiteralValue creates a value from an integer literal, registered
// to the given key.
func (c *Codegen) CreateIntLiteralValue(key GetText, typ types.Type, value string) *Value {
	val := &Value{
		typ:     typ,
		backing: ConstantBacking,
		value:   value,
	}

	c.setValue(key, val)
	return val
}

// CreateValue registers a value against a key.
func (c *Codegen) CreateValue(key GetText, val *Value) {
	if c.values[key] != nil {
		c.fail("CreateValue(): key(%s) already exists", key.GetText())
	}

	c.values[key] = val
}

// GetValue returns the value associated with the key.
func (c *Codegen) GetValue(key GetText) *Value {
	val := c.values[key]
	if val == nil {
		c.fail("value not found!")
	}
	return val
}

// MoveValue reregisters a value with a new key.
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

// ReleaseValue unregisters a value and unbinds the associated register if necessary.
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

// Add generates code for the sum of two values and registers the result under the given key.
func (c *Codegen) Add(key GetText, v1, v2 *Value) {
	val := c.allocateTransientValue(v1.typ)

	// TODO: use proper register size
	fmt.Fprintf(c.out, "\tmov %s, %s\n", val.Source(), v1.Source())
	fmt.Fprintf(c.out, "\tadd %s, %s\n", val.Source(), v2.Source())

	c.setValue(key, val)
}

// Subtract geneates code for the difference of two values and registers the result under the given key.
func (c *Codegen) Subtract(key GetText, v1, v2 *Value) {
	val := c.allocateTransientValue(v1.typ)

	// TODO: use proper register size
	fmt.Fprintf(c.out, "\tmov %s, %s\n", val.Source(), v1.Source())
	fmt.Fprintf(c.out, "\tsub %s, %s\n", val.Source(), v2.Source())

	c.setValue(key, val)
}

// ReturnValue generates code to return a value from the function.
func (c *Codegen) ReturnValue(key GetText) {
	// look up the value
	val := c.GetValue(key)

	// grab the rax register
	rax := c.integerReg[RAX]

	if val.register == rax {
		// the value is already in RAX, nothing to do
	} else {
		if rax.binding != nil {
			c.fail("rax bound during return")
		}

		// move the value to rax
		fmt.Fprintf(c.out, "\tmov %s, %s\n", rax.dwordName, val.Source())
	}

	// jmp to the epilogue
	c.Return()
}

// Return generates code to return nothing from the function.
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
	for _, regName := range integerRegisterAllocationOrder {
		reg := c.integerReg[regName]
		if reg.binding == nil {
			return reg
		}
	}

	c.fail("no integer register found")

	return nil
}

func (c *Codegen) allocateTransientValue(typ types.Type) *Value {
	// allocate a register
	reg := c.allocateRegister()

	// allocate a stack value
	offset := c.stack.Alloc4() // TODO: proper size
	val := NewLocalValue(typ, offset)

	// bind the register to the value (we assume the register will be written)
	c.bind(val, reg, true)

	return val
}

func (c *Codegen) bind(val *Value, reg *Register, dirty bool) {
	val.register = reg
	val.dirty = dirty
	reg.binding = val
}

func (c *Codegen) endFrameLabel() string {
	return fmt.Sprintf("%s.end", c.funcName)
}

func (c *Codegen) fail(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	panic(msg)
}
