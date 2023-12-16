package types

import (
	"strings"

	"github.com/TassoKarkanis/minic/parser"
)

// void, int, float
// int*
// int**
// int[3]
// void f(int)

// Type category
const (
	Basic = iota
	Pointer
	Struct
	Function
)

type Type interface {
	Category() int
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

func (t *BasicType) Category() int {
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

func (t *PointerType) Category() int {
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

func (t *StructType) Category() int {
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

func (t *FunctionType) Category() int {
	return Function
}

func (t *FunctionType) String() string {
	rep := ""
	for i, param := range t.Params {
		if i > 0 {
			rep += ", "
		}

		paramStr := "<none>"
		if param.Type != nil {
			paramStr = param.Type.String()
		}

		rep += paramStr
	}

	returnType := "<none>"
	if t.ReturnType != nil {
		returnType = t.ReturnType.String()
	}
	rep = returnType + t.Name + "(" + rep + ")"
	return rep
}
