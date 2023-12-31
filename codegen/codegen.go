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
		for k, v := range c.values {
			c.fail("unreleased value: (%s) (%v)", k.GetText(), v)
			break
		}

	}
}

// CreateFunction returns a value that represents a function.
func (c *Codegen) CreateFunction(name string, typ types.Type) *Value {
	return newGlobalValue(c, name, typ)
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
		val := newLocalValue(c, typ, offset, true)
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
		owner:   c,
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
		owner:   c,
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
	val.Retain()
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

// ReleaseValue unregisters a value and unbinds the associated register if necessary.
func (c *Codegen) ReleaseValue(key GetText) {
	if c.values[key] == nil {
		c.fail("value to release not found: %s", key.GetText())
	}

	val := c.values[key]
	delete(c.values, key)
	val.Release()
}

func (c *Codegen) MoveValue(destKey, srcKey GetText) {
	dest := c.values[destKey]
	if dest == nil {
		c.fail("MoveValue(): destination not found")
	}

	c.Move(dest, srcKey)
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

	// ignore statements like "x = x;"
	if dest == src {
		return
	}

	// generate move instructions, but note that only one param can be an address
	destBound := dest.IsBound()
	srcBound := src.IsBound()
	switch {
	case destBound && srcBound:
		fallthrough
	case destBound && !srcBound:
		// move to the destination register
		fmt.Fprintf(c.out, "\tmov %s, %s ; update l-value register\n", dest.Source(), src.Source())

		// move to the destination address
		fmt.Fprintf(c.out, "\tmov %s, %s ; store l-value\n", dest.RawSource(), dest.Source())

	case !destBound && srcBound:
		fmt.Fprintf(c.out, "\tmov %s, %s ; store l-value\n", dest.Source(), src.Source())

	case !destBound && !srcBound:
		// allocate a register for the source
		c.allocateValueRegister(src)

		// load source value
		fmt.Fprintf(c.out, "\tmov %s, %s ; load mov source\n", src.Source(), src.RawSource())

		// move to the destination
		fmt.Fprintf(c.out, "\tmov %s, %s ; store l-value\n", dest.Source(), src.Source())
	}
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

// UnaryMinus generates code to negate a value and registers the result under the given key.
func (c *Codegen) UnaryMinus(key GetText, v *Value) {
	// allocate a new value
	val := c.allocateTransientValue(v.typ, nil)

	// load the source
	// TODO: use proper register size
	fmt.Fprintf(c.out, "\tmov %s, %s ; load for negation\n", val.Source(), v.Source())

	// negate the destination
	fmt.Fprintf(c.out, "\tneg %s\n", val.Source())

	c.setValue(key, val)
}

// Multiply generates code for the product of two values and registers the result under the given key.
func (c *Codegen) Multiply(key GetText, v1, v2 *Value) {
	// Only RAX can do multiplication.  Also, the high bits end up in RDX.

	// unbind rdx
	rdx := c.integerReg[RDX]
	c.unbindRegister(rdx, false)

	// make rax contain the left side
	rax := c.integerReg[RAX]
	if v1.register == nil || v1.register != rax {
		c.unbindRegister(rax, false)
		fmt.Fprintf(c.out, "\tmov %s, %s\n", rax.dwordName, v1.Source())
	}

	// perform the multiplication
	fmt.Fprintf(c.out, "\tmul %s\n", v2.Source())

	// allocate and set the value
	val := c.allocateTransientValue(v1.typ, rax)
	c.setValue(key, val)
}

// Divide generates code for the quotient of two values and registers the result under the given key.
func (c *Codegen) Divide(key GetText, v1, v2 *Value) {
	c.divide(key, v1, v2)

	// allocate and set the value
	rax := c.integerReg[RAX]
	val := c.allocateTransientValue(v1.typ, rax)
	c.setValue(key, val)
}

// Modulus generates code for the modulus of two values and registers the result under the given key.
func (c *Codegen) Modulus(key GetText, v1, v2 *Value) {
	c.divide(key, v1, v2)

	// allocate and set the value
	rdx := c.integerReg[RDX]
	val := c.allocateTransientValue(v1.typ, rdx)
	c.setValue(key, val)
}

func (c *Codegen) divide(key GetText, v1, v2 *Value) {
	// Only RAX can do division.  Also, the remainder end up in RDX.

	// unbind rdx and clear it
	rdx := c.integerReg[RDX]
	c.unbindRegister(rdx, false)
	fmt.Fprintf(c.out, "\tmov %s, 0\n", rdx.fullName)

	// make rax contain the left side
	rax := c.integerReg[RAX]
	if v1.register == nil || v1.register != rax {
		c.unbindRegister(rax, false)
		fmt.Fprintf(c.out, "\tmov %s, %s\n", rax.dwordName, v1.Source())
	}

	// perform the division
	fmt.Fprintf(c.out, "\tcdq\n")
	fmt.Fprintf(c.out, "\tidiv %s\n", v2.Source())
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
		// make rax available
		c.unbindRegister(rax, true)

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

func (c *Codegen) releaseValue(val *Value) {
	c.unbindValue(val, true)
	val.owner = nil
}

func (c *Codegen) setValue(key GetText, val *Value) {
	if c.values[key] != nil {
		c.fail("value already set for %s", key.GetText())
	}

	c.values[key] = val
	val.Retain()
}

func (c *Codegen) allocateValueRegister(val *Value) {
	if val.IsBound() {
		c.fail("allocateValueRegister() called for bound value")
	}

	reg := c.allocateRegister()
	c.bind(val, reg, false)
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
	val := newLocalValue(c, typ, offset, false)

	// bind the register to the value (we assume the register will be written)
	c.bind(val, reg, true)

	return val
}

func (c *Codegen) bind(val *Value, reg *Register, dirty bool) {
	val.register = reg
	val.dirty = dirty
	reg.binding = val
}

func (c *Codegen) unbindValue(val *Value, discard bool) {
	if val.register == nil {
		return
	}

	c.unbindRegister(val.register, discard)
}

func (c *Codegen) unbindRegister(reg *Register, discard bool) {
	if reg.binding == nil {
		return // nothing to do
	}

	// unbind
	val := reg.binding
	reg.binding = nil
	val.register = nil

	// store the register value to its backing, if necessary
	if val.dirty && !discard {
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
