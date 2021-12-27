package main

import (
	"fmt"
	"io"

	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/symbols"
	"github.com/TassoKarkanis/minic/types"
)

type MainPass struct {
	*parser.BaseCListener
	Output                   io.Writer
	Symbols                  *symbols.Table
	LastDeclaratorIdentifier string
	LastType                 int // i.e. CParserVoid
	LastFunction             *types.FunctionType
	Codegen                  *Codegen
}

func NewMainPass(output io.Writer) *MainPass {
	return &MainPass{
		Output:  output,
		Symbols: symbols.NewTable(),
	}
}

func (c *MainPass) EnterCompilationUnit(ctx *parser.CompilationUnitContext) {
	fmt.Printf("EnterCompilationUnit(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitCompilationUnit(ctx *parser.CompilationUnitContext) {
	fmt.Printf("ExitCompilationUnit()\n")
}

func (c *MainPass) EnterExternalDeclaration(ctx *parser.ExternalDeclarationContext) {
	fmt.Printf("EnterExternalDeclaration(): %s\n", ctx.GetText())
}

func (c *MainPass) EnterFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	fmt.Printf("EnterFunctionDefinition(): %s\n", ctx.GetText())
	c.LastFunction = &types.FunctionType{}
	c.Symbols.PushScope()
	c.Codegen = NewCodegen()
}

func (c *MainPass) ExitFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	fmt.Printf("ExitFunctionDefinition(): %s\n", ctx.GetText())
	fmt.Fprintf(c.Output, "\tret\n\n")
	c.Symbols.PopScope()
	c.LastFunction = nil
}

func (c *MainPass) EnterDeclarator(ctx *parser.DeclaratorContext) {
	fmt.Printf("EnterDeclarator(): %s\n", ctx.GetText())
}

func (c *MainPass) EnterDirectDeclaratorFunction(ctx *parser.DirectDeclaratorFunctionContext) {
	fmt.Printf("EnterDirectDeclaratorFunction(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitDirectDeclaratorFunction(ctx *parser.DirectDeclaratorFunctionContext) {
	fmt.Printf("ExitDirectDeclaratorFunction(): %s\n", ctx.GetText())
}

func (c *MainPass) EnterDirectDeclaratorIdentifier(ctx *parser.DirectDeclaratorIdentifierContext) {
	fmt.Printf("EnterDirectDeclaratorIdentifier(): %s\n", ctx.GetText())
	c.LastDeclaratorIdentifier = ctx.Identifier().GetText()
	fmt.Printf("LastDeclaratorIdentifier: %s\n", c.LastDeclaratorIdentifier)
}

func (c *MainPass) EnterTypeSpecifierSimple(ctx *parser.TypeSpecifierSimpleContext) {
	fmt.Printf("EnterTypeSpecifierSimple(): %s\n", ctx.GetText())
	fmt.Printf("  GetStart(): %s\n", ctx.GetStart().GetText())
	c.LastType = ctx.GetStart().GetTokenType()
	fmt.Printf("LastType: %v\n", c.LastType)
}

func (c *MainPass) EnterParameterTypeList(ctx *parser.ParameterTypeListContext) {
	fmt.Printf("EnterParameterTypeList(): %s\n", ctx.GetText())
	cp := ctx.GetParser().(*parser.CParser)
	name := c.LastDeclaratorIdentifier
	c.LastFunction.ReturnType = types.NewBasicType(c.LastType, cp)
	c.LastFunction.Name = name
	c.Symbols.AddSymbol(name, c.LastFunction)
	fmt.Fprintf(c.Output, "%s:\n", name)
}

func (c *MainPass) EnterCompoundStatement(ctx *parser.CompoundStatementContext) {
	fmt.Printf("EnterCompoundStatement(): %s\n", ctx.GetText())
}

func (c *MainPass) EnterJumpStatement(ctx *parser.JumpStatementContext) {
	fmt.Printf("EnterJumpStatement(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitJumpStatement(ctx *parser.JumpStatementContext) {
	fmt.Printf("ExitJumpStatement(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(true)

	switch {
	case ctx.Return() != nil:
		expr := ctx.Expression()
		if expr != nil {
			fmt.Printf(" got value!\n")
			value := cgen.GetValue(expr)
			reg1 := cgen.GetReturnRegister()
			load := cgen.LoadValue(value)
			fmt.Fprintf(c.Output, "\tmov %s, %s\n", reg1.Name(4), load)
		}
		fmt.Fprintf(c.Output, "\tret\n")
	}
}

func (c *MainPass) EnterExpression(ctx *parser.ExpressionContext) {
	fmt.Printf("EnterExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitExpression(ctx *parser.ExpressionContext) {
	fmt.Printf("ExitExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(true)

	exp1 := ctx.AssignmentExpression()
	exp2 := ctx.Expression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterAssignmentExpression(ctx *parser.AssignmentExpressionContext) {
	fmt.Printf("EnterAssignmentExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitAssignmentExpression(ctx *parser.AssignmentExpressionContext) {
	fmt.Printf("ExitAssignmentExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(true)

	exp1 := ctx.ConditionalExpression()

	switch {
	case exp1 != nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) ExitConditionalExpression(ctx *parser.ConditionalExpressionContext) {
	fmt.Printf("ExitConditionalExpression(): %s\n", ctx.GetText())
	fmt.Printf("  LogicalOrExpression: %v\n", ctx.LogicalOrExpression() != nil)

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.LogicalOrExpression()
	cgen.CopyValue(ctx, exp1)
}

func (c *MainPass) EnterLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) {
	fmt.Printf("EnterLogicalOrExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) {
	fmt.Printf("ExitLogicalOrExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.LogicalAndExpression()
	exp2 := ctx.LogicalOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) {
	fmt.Printf("EnterLogicalAndExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) {
	fmt.Printf("ExitLogicalAndExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.InclusiveOrExpression()
	exp2 := ctx.LogicalAndExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterInclusiveOrExpression(ctx *parser.InclusiveOrExpressionContext) {
	fmt.Printf("EnterInclusiveOrExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitInclusiveOrExpression(ctx *parser.InclusiveOrExpressionContext) {
	fmt.Printf("ExitInclusiveOrExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.ExclusiveOrExpression()
	exp2 := ctx.InclusiveOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterExclusiveOrExpression(ctx *parser.ExclusiveOrExpressionContext) {
	fmt.Printf("EnterExclusiveOrExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitExclusiveOrExpression(ctx *parser.ExclusiveOrExpressionContext) {
	fmt.Printf("ExitExclusiveOrExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.AndExpression()
	exp2 := ctx.ExclusiveOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterAndExpression(ctx *parser.AndExpressionContext) {
	fmt.Printf("EnterAndExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitAndExpression(ctx *parser.AndExpressionContext) {
	fmt.Printf("ExitAndExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.EqualityExpression()
	exp2 := ctx.AndExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterEqualityExpression(ctx *parser.EqualityExpressionContext) {
	fmt.Printf("EnterEqualityExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitEqualityExpression(ctx *parser.EqualityExpressionContext) {
	fmt.Printf("ExitEqualityExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.RelationalExpression()
	exp2 := ctx.EqualityExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterRelationalExpression(ctx *parser.RelationalExpressionContext) {
	fmt.Printf("EnterRelationalExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitRelationalExpression(ctx *parser.RelationalExpressionContext) {
	fmt.Printf("ExitRelationalExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.ShiftExpression()
	exp2 := ctx.RelationalExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterShiftExpression(ctx *parser.ShiftExpressionContext) {
	fmt.Printf("EnterShiftExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitShiftExpression(ctx *parser.ShiftExpressionContext) {
	fmt.Printf("ExitShiftExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.AdditiveExpression()
	exp2 := ctx.ShiftExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	fmt.Printf("EnterAdditiveExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	fmt.Printf("ExitAdditiveExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.MultiplicativeExpression()
	exp2 := ctx.AdditiveExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) {
	fmt.Printf("EnterMultiplicativeExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) {
	fmt.Printf("ExitMultiplicativeExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.CastExpression()
	exp2 := ctx.MultiplicativeExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterCastExpression(ctx *parser.CastExpressionContext) {
	fmt.Printf("EnterCastExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitCastExpression(ctx *parser.CastExpressionContext) {
	fmt.Printf("ExitCastExpression(): %s\n", ctx.GetText())
	fmt.Printf("  TypeName: %v\n", ctx.TypeName() != nil)
	fmt.Printf("  CastExpression: %v\n", ctx.CastExpression() != nil)
	fmt.Printf("  UnaryExpression: %v\n", ctx.UnaryExpression() != nil)

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.TypeName()
	exp2 := ctx.CastExpression()
	exp3 := ctx.UnaryExpression()

	switch {
	case exp1 == nil && exp2 == nil && exp3 != nil:
		cgen.CopyValue(ctx, exp3)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterUnaryExpression(ctx *parser.UnaryExpressionContext) {
	fmt.Printf("EnterUnaryExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitUnaryExpression(ctx *parser.UnaryExpressionContext) {
	fmt.Printf("ExitUnaryExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.PostfixExpression()

	switch {
	case exp1 != nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterPostfixExpression(ctx *parser.PostfixExpressionContext) {
	fmt.Printf("EnterPostfixExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitPostfixExpression(ctx *parser.PostfixExpressionContext) {
	fmt.Printf("ExitPostfixExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	exp1 := ctx.PrimaryExpression()

	switch {
	case exp1 != nil:
		cgen.CopyValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterPrimaryExpression(ctx *parser.PrimaryExpressionContext) {
	fmt.Printf("EnterPrimaryExpression(): %s\n", ctx.GetText())
	tok := ctx.GetStart()
	fmt.Printf("  start token index: %d\n", tok.GetTokenIndex())
	fmt.Printf("  start token type: %d\n", tok.GetTokenType())
	fmt.Printf("  start token channel: %d\n", tok.GetChannel())
}

func (c *MainPass) ExitPrimaryExpression(ctx *parser.PrimaryExpressionContext) {
	fmt.Printf("ExitPrimaryExpression(): %s\n", ctx.GetText())

	cgen := NewCodegenContext(c.Codegen)
	defer cgen.Close(false)

	cp := ctx.GetParser().(*parser.CParser)
	tok := ctx.GetStart()

	switch {
	case tok.GetTokenType() == parser.CLexerConstant:
		typ := types.NewBasicType(parser.CParserInt, cp)
		val := cgen.NewConstant(tok.GetText(), typ)
		cgen.SetValue(ctx, val)

	default:
		panic("unhandled case!")
	}
}
