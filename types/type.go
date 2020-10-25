package types

import (
	"strings"

	"github.com/TassoKarkanis/minic/parser"
)

// Type category
const (
	Basic    = iota
	Pointer  = iota
	Struct   = iota
	Function = iota
)

type Type interface {
	GetCategory() int
	String() string
}

//
// BasicType
//

type BasicType struct {
	Type      int // i.e. CParserVoid, CParserInt
	stringRep string
}

func NewBasicType(type_ int, cp *parser.CParser) *BasicType {
	return &BasicType{
		Type:      type_,
		stringRep: cp.LiteralNames[type_],
	}
}

func (t *BasicType) GetCategory() int {
	return Basic
}

func (t *BasicType) String() string {
	return t.stringRep
}

//
// PointerType
//

type PointerType struct {
	DerefType  Type
	DerefCount int
}

func (t *PointerType) GetCategory() int {
	return Pointer
}

func (t *PointerType) String() string {
	return t.DerefType.String() + strings.Repeat("*", t.DerefCount)
}

//
// StructType
//

type Field struct {
	Name string
	Type Type
}

type StructType struct {
	Name   string
	Fields []Field
}

func (t *StructType) GetCategory() int {
	return Struct
}

func (t *StructType) String() string {
	return "struct " + t.Name
}

//
// FunctionType
//

type Param struct {
	Name string
	Type Type
}

type FunctionType struct {
	Name       string
	ReturnType Type
	Params     []Param
}

func (t *FunctionType) GetCategory() int {
	return Function
}

func (t *FunctionType) String() string {
	rep := ""
	for i, param := range t.Params {
		if i > 0 {
			rep += ", "
		}
		rep += param.Type.String()
	}
	rep = t.ReturnType.String() + t.Name + "(" + rep + ")"
	return rep
}
