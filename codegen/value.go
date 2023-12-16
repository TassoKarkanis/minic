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
	offset   int        // typically negative, set local (stack) storage
	register *Register  // set if register storage
}

func NewGlobalValue(name string, typ types.Type) *Value {
	return &Value{
		typ:     typ,
		storage: GlobalStorage,
		value:   name,
	}
}

func (v *Value) GetType() types.Type {
	return v.typ
}

func (v *Value) Source() string {
	switch v.storage {
	case ConstantStorage:
		return v.value

	case RegisterStorage:
		return v.register.Name(4)

	case LocalStorage:
		// TODO: use type
		return fmt.Sprintf("dword [rsp - %d]", -v.offset)

	case GlobalStorage:
		return v.value

	default:
		panic("Value.Source() undefined case!")
	}
}
