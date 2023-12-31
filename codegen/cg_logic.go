package codegen

import (
	"fmt"
)

// Equal generates code to evaluate equality and registers the result under the given key.
func (c *Codegen) Equal(key GetText, v1, v2 *Value) {
	c.logicalCompare(key, "equal", "sete", v1, v2)
}

// NotEqual generates code to evaluate inequality and registers the result under the given key.
func (c *Codegen) NotEqual(key GetText, v1, v2 *Value) {
	c.logicalCompare(key, "not-equal", "setne", v1, v2)
}

func (c *Codegen) Less(key GetText, v1, v2 *Value) {
	c.logicalCompare(key, "less", "setl", v1, v2)
}

func (c *Codegen) LessEqual(key GetText, v1, v2 *Value) {
	c.logicalCompare(key, "less-equal", "setle", v1, v2)
}

func (c *Codegen) Greater(key GetText, v1, v2 *Value) {
	c.logicalCompare(key, "greater", "setg", v1, v2)
}

func (c *Codegen) GreaterEqual(key GetText, v1, v2 *Value) {
	c.logicalCompare(key, "greater-equal", "setge", v1, v2)
}

func (c *Codegen) logicalCompare(key GetText, cmpType, setInstr string, v1, v2 *Value) {
	// allocate a new value for the result
	val := c.allocateTransientValue(v1.typ, nil)

	// move the left operand into the value
	fmt.Fprintf(c.out, "\tmov %s, %s ; %s: load LHS\n", val.Source(), v1.Source(), cmpType)

	// compare with right operand
	fmt.Fprintf(c.out, "\tcmp %s, %s ; %s: compare\n", val.Source(), v2.Source(), cmpType)

	// put the result in the low byte
	fmt.Fprintf(c.out, "\t%s %s ; %s: set byte in result\n", setInstr, val.register.byteName, cmpType)

	// zero extend to rest of value
	fmt.Fprintf(c.out, "\tmovzx %s, %s ; %s: zero-extend\n", val.Source(), val.register.byteName, cmpType)

	c.setValue(key, val)
}
