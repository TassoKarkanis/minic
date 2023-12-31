package codegen

import (
	"fmt"
)

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
