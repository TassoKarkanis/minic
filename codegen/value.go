package codegen

import (
	"fmt"

	"github.com/TassoKarkanis/minic/types"
)

type Storage int

const (
	ConstantStorage = iota
	RegisterStorage
	LocalStorage
	GlobalStorage
)

type Value struct {
	typ      types.Type // type of value
	storage  Storage    // how the value is stored
	value    string     // constant, global address
	offset   int        // relative to BP, typically negative
	register *Register  // set if register storage
	dirty    bool       // value in register is more recent than at storage location
}

func NewGlobalValue(name string, typ types.Type) *Value {
	return &Value{
		typ:     typ,
		storage: GlobalStorage,
		value:   name,
	}
}

func NewIdentifier(name string, typ types.Type, offset int) *Value {
	return &Value{
		typ:     typ,
		storage: LocalStorage,
		offset:  offset,
	}
}

func (v *Value) GetType() types.Type {
	return v.typ
}

func (v *Value) Source() string {
	useRegister := v.storage == RegisterStorage || v.register != nil

	switch {
	case v.storage == ConstantStorage:
		return v.value

	case useRegister:
		return v.register.Name(4)

	case v.storage == LocalStorage:
		// TODO: use type
		return fmt.Sprintf("dword [rsp - %d]", -v.offset)

	case v.storage == GlobalStorage:
		return v.value

	default:
		panic("Value.Source() undefined case!")
	}
}
