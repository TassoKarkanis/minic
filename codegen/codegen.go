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
		val := NewLocalValue(typ, offset, true)
		c.bind(val, c.integerReg[integerParameterOrder[i]], true)
		values = append(values, val)

		fmt.Fprintf(c.out, "\t; param %s %s -> %s [rsp - %d]\n",
			typ.String(), paramNames[i], val.register.Name(4), val.offset) // TODO
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

func (c *Codegen) CreateLocalValue(name string, typ types.Type) *Value {
	// allocate a stack offset
	offset := c.stack.Alloc4() // TODO: size

	// allocate a register
	reg := c.allocateRegister()

	// make the value value
	val := &Value{
		typ:     typ,
		backing: LocalBacking,
		offset:  offset,
		lvalue:  true,
	}

	// bind the register
	c.bind(val, reg, false)

	fmt.Fprintf(c.out, "\t; var %s %s -> %s [rsp - %d]\n",
		typ.String(), name, val.register.Name(4), val.offset) // TODO

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

// TransferValue reregisters a value with a new key.
func (c *Codegen) TransferValue(destKey, srcKey GetText) {
	if c.values[destKey] != nil {
		c.fail("TransferValue(): destkey(%s) already exists", destKey.GetText())
	}

	if c.values[srcKey] == nil {
		c.fail("TransferValue(): srckey(%s) not found", destKey.GetText())
	}

	c.values[destKey] = c.values[srcKey]
	delete(c.values, srcKey)
}

// Move copies the source to the destination.
func (c *Codegen) Move(dest *Value, srcKey GetText) {
	src := c.values[srcKey]
	if src == nil {
		c.fail("Move(): source not found: %s", srcKey.GetText())
	}

	// dest should be an lvalue
	if !dest.lvalue {
		c.fail("Move(): moving into rvalue")
	}

	// unbind destination register, if any
	c.unbindValue(dest)

	// move src -> dest
	fmt.Fprintf(c.out, "\tmov %s, %s\n", dest.Source(), src.Source())
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
	val := c.allocateTransientValue(v1.typ, nil)

	// TODO: use proper register size
	fmt.Fprintf(c.out, "\tmov %s, %s\n", val.Source(), v1.Source())
	fmt.Fprintf(c.out, "\tadd %s, %s\n", val.Source(), v2.Source())

	c.setValue(key, val)
}

// Subtract generates code for the difference of two values and registers the result under the given key.
func (c *Codegen) Subtract(key GetText, v1, v2 *Value) {
	val := c.allocateTransientValue(v1.typ, nil)

	// TODO: use proper register size
	fmt.Fprintf(c.out, "\tmov %s, %s\n", val.Source(), v1.Source())
	fmt.Fprintf(c.out, "\tsub %s, %s\n", val.Source(), v2.Source())

	c.setValue(key, val)
}

// Multiply generates code for the product of two values and registers the result under the given key.
func (c *Codegen) Multiply(key GetText, v1, v2 *Value) {
	// Only RAX can do multiplication.  Also, the high bits end up in RDX.

	// unbind rdx
	rdx := c.integerReg[RDX]
	c.unbindRegister(rdx)

	// make rax contain the left side
	rax := c.integerReg[RAX]
	if v1.register == nil || v1.register != rax {
		c.unbindRegister(rax)
		fmt.Fprintf(c.out, "\tmov %s, %s\n", rax.dwordName, v1.Source())
	}

	// perform the multiplication
	fmt.Fprintf(c.out, "\tmul %s\n", v2.Source())

	// allocate and set the value
	val := c.allocateTransientValue(v1.typ, rax)
	c.setValue(key, val)
}

// Divide generates code for the dividend of two values and registers the result under the given key.
func (c *Codegen) Divide(key GetText, v1, v2 *Value) {
	c.fail("divide not yet implemented")
}

// Modulus generates code for the modulus of two values and registers the result under the given key.
func (c *Codegen) Modulus(key GetText, v1, v2 *Value) {
	c.fail("modulus not yet implemented")
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

func (c *Codegen) allocateTransientValue(typ types.Type, reg *Register) *Value {
	// allocate a register if necessary
	if reg == nil {
		reg = c.allocateRegister()
	}

	// allocate a stack value
	offset := c.stack.Alloc4() // TODO: proper size
	val := NewLocalValue(typ, offset, false)

	// bind the register to the value (we assume the register will be written)
	c.bind(val, reg, true)

	return val
}

func (c *Codegen) bind(val *Value, reg *Register, dirty bool) {
	val.register = reg
	val.dirty = dirty
	reg.binding = val
}

func (c *Codegen) unbindValue(val *Value) {
	if val.register == nil {
		return
	}

	c.unbindRegister(val.register)
}

func (c *Codegen) unbindRegister(reg *Register) {
	if reg.binding == nil {
		return // nothing to do
	}

	// unbind
	val := reg.binding
	reg.binding = nil
	val.register = nil

	// store the register value to its backing, if necessary
	if val.dirty {
		switch val.backing {
		case LocalBacking:
			fmt.Fprintf(c.out, "\tmov %s, %s ; flush to local\n", val.Source(), reg.dwordName)

		default:
			c.fail("unhandled case in unbindRegister()")
		}

		val.dirty = false
	}
}

func (c *Codegen) endFrameLabel() string {
	return fmt.Sprintf("%s.end", c.funcName)
}

func (c *Codegen) fail(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	panic(msg)
}
