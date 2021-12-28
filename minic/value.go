package main

import (
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
	typ      types.Type
	storage  Storage
	constant string
	register *Register
}

func (v *Value) GetType() types.Type {
	return v.typ
}

func (v *Value) Source() string {
	if v.register != nil {
		return v.register.Name(4)
	} else {
		return v.constant
	}
}
