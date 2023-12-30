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
	owner    *Codegen   // owning code generator
	refCount int        // reference count, see Retain() and Release()
	typ      types.Type // type of value
	backing  Backing    // how the value is stored
	value    string     // constant or global address
	offset   int        // relative to BP, stored positively, set for LocalBacking
	register *Register  // set if register is "bound" to this value
	dirty    bool       // value in register is more recent than at storage location
	lvalue   bool       // true if the Value is an l-value
}

func newGlobalValue(owner *Codegen, name string, typ types.Type) *Value {
	return &Value{
		owner:   owner,
		typ:     typ,
		backing: GlobalBacking,
		value:   name,
		lvalue:  true,
	}
}

func newLocalValue(owner *Codegen, typ types.Type, offset int, lvalue bool) *Value {
	return &Value{
		owner:   owner,
		typ:     typ,
		backing: LocalBacking,
		offset:  offset,
		lvalue:  lvalue,
	}
}

func (v *Value) Retain() {
	if v.owner == nil {
		panic("Retain() on dead value")
	}

	v.refCount++
}

func (v *Value) Release() {
	if v.owner == nil {
		panic("Release() on dead value")
	}

	v.refCount--
	if v.refCount == 0 {
		v.owner.releaseValue(v)
	}
}

func (v *Value) GetType() types.Type {
	return v.typ
}

// Source returns the string representation of the bound register,
// constant, or address of the value.
func (v *Value) Source() string {
	if v.register != nil {
		return v.register.Name(4) // TODO: size
	} else {
		return v.RawSource()
	}
}

// RawSource returns the string representation of the value ignoring
// the register assignment.
func (v *Value) RawSource() string {
	switch {
	case v.backing == ConstantBacking:
		return v.value

	case v.backing == LocalBacking:
		return fmt.Sprintf("dword [rsp - %d]", v.offset) // TODO: size

	default:
		panic("Value.Source() undefined case!")
	}
}

// IsBound returns true if the value currently has a register assigment.
func (v *Value) IsBound() bool {
	return v.register != nil
}
