package codegen

import (
	"fmt"

	"github.com/TassoKarkanis/minic/types"
)

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
