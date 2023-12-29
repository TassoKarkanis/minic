package codegen

import (
	"fmt"

	"github.com/TassoKarkanis/minic/types"
)

// Backup storage for a Value when not in a register.
type Backing int

const (
	ConstantBacking Backing = iota // value is a constant
	GlobalBacking                  // value refers to a global variable
	LocalBacking                   // value refers to a parameter or automatic variable
)

type Value struct {
	typ      types.Type // type of value
	backing  Backing    // how the value is stored
	value    string     // constant or global address
	offset   int        // relative to BP, stored positively, set for LocalBacking
	register *Register  // set if register is "bound" to this value
	dirty    bool       // value in register is more recent than at storage location
	lvalue   bool       // true if the Value is an l-value
}

func NewGlobalValue(name string, typ types.Type) *Value {
	return &Value{
		typ:     typ,
		backing: GlobalBacking,
		value:   name,
		lvalue:  true,
	}
}

func NewLocalValue(typ types.Type, offset int, lvalue bool) *Value {
	return &Value{
		typ:     typ,
		backing: LocalBacking,
		offset:  offset,
		lvalue:  lvalue,
	}
}

func (v *Value) GetType() types.Type {
	return v.typ
}

func (v *Value) Source() string {
	switch {
	case v.register != nil:
		return v.register.Name(4) // TODO: size

	case v.backing == ConstantBacking:
		return v.value

	case v.backing == LocalBacking:
		// TODO: use type/size
		return fmt.Sprintf("dword [rsp - %d]", v.offset)

	default:
		panic("Value.Source() undefined case!")
	}
}
