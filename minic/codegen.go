package main

import (
	"github.com/TassoKarkanis/minic/types"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type Register struct {
	FullName  string
	DWordName string
	WordName  string
	ByteName  string
}

func (r Register) Name(byteSize int) string {
	switch byteSize {
	case 8:
		return r.FullName

	case 4:
		return r.DWordName

	case 2:
		return r.WordName

	case 1:
		return r.ByteName

	default:
		panic("invalid bytesize!")
	}
}

type Value interface {
	GetType() types.Type
}

type constValue struct {
	typ   types.Type
	value string
}

func (v *constValue) GetType() types.Type {
	return v.typ
}

type Codegen struct {
	rax    *Register
	rdi    *Register
	rsi    *Register
	rdx    *Register
	rcx    *Register
	r8     *Register
	r9     *Register
	values map[antlr.ParserRuleContext]Value
}

func NewCodegen() *Codegen {
	return &Codegen{
		rax: &Register{
			FullName:  "rax",
			DWordName: "eax",
			WordName:  "ax",
			ByteName:  "al",
		},
		rdi: &Register{
			FullName:  "rdi",
			DWordName: "edi",
			WordName:  "di",
			ByteName:  "dl",
		},
		rsi: &Register{
			FullName:  "rsi",
			DWordName: "esx",
			WordName:  "si",
			ByteName:  "sl",
		},
		rdx: &Register{
			FullName:  "rdx",
			DWordName: "edx",
			WordName:  "dx",
			ByteName:  "dl",
		},
		rcx: &Register{
			FullName:  "rcx",
			DWordName: "ecx",
			WordName:  "cx",
			ByteName:  "cl",
		},
		r8: &Register{
			FullName:  "r8",
			DWordName: "r8d",
			WordName:  "r8w",
			ByteName:  "r8b",
		},
		r9: &Register{
			FullName:  "r9",
			DWordName: "r9d",
			WordName:  "r9w",
			ByteName:  "r9b",
		},
		values: make(map[antlr.ParserRuleContext]Value),
	}
}

func (c *Codegen) closeContext(flushValues bool) {
}

type CodegenContext struct {
	cgen *Codegen
}

func NewCodegenContext(cgen *Codegen) *CodegenContext {
	return &CodegenContext{
		cgen: cgen,
	}
}

func (c *CodegenContext) Close(flushValues bool) {
	c.cgen.closeContext(flushValues)
}

func (c *CodegenContext) NewConstant(val string, typ types.Type) Value {
	return &constValue{
		typ:   typ,
		value: val,
	}
}

func (c *CodegenContext) GetValue(ruleCtx antlr.ParserRuleContext) Value {
	val := c.cgen.values[ruleCtx]
	if val == nil {
		panic("value not found!")
	}
	return val
}

func (c *CodegenContext) SetValue(ruleCtx antlr.ParserRuleContext, val Value) {
	if ruleCtx == nil {
		panic("assignment of value to null ruleCtx")
	}

	if c.cgen.values[ruleCtx] != nil {
		panic("value already assigned!")
	}
	c.cgen.values[ruleCtx] = val
}

func (c *CodegenContext) CopyValue(dest, src antlr.ParserRuleContext) {
	val := c.GetValue(src)
	c.SetValue(dest, val)
}

func (c *CodegenContext) GetReturnRegister() *Register {
	return c.cgen.rax
}

func (c *CodegenContext) LoadValue(val Value) string {
	cval, ok := val.(*constValue)
	if !ok {
		panic("LoadValue() only handles constants!")
	}
	return cval.value
}
