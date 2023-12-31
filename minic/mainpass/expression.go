package mainpass

import (
	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/types"
)

func (c *MainPass) EnterExpression(ctx *parser.ExpressionContext) {
	c.enterRule(ctx, "Expression")
}

func (c *MainPass) ExitExpression(ctx *parser.ExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	exp1 := ctx.AssignmentExpression()
	exp2 := ctx.Expression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.TransferValue(ctx, exp1)

	default:
		c.fail("ExitExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterConstantExpression(ctx *parser.ConstantExpressionContext) {
	c.enterRule(ctx, "ConstantExpression")
}

func (c *MainPass) ExitConstantExpression(ctx *parser.ConstantExpressionContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterAssignmentExpression(ctx *parser.AssignmentExpressionContext) {
	c.enterRule(ctx, "AssignmentExpression")
}

func (c *MainPass) ExitAssignmentExpression(ctx *parser.AssignmentExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	e1 := ctx.ConditionalExpression()
	e2 := ctx.UnaryExpression()
	e3 := ctx.AssignmentExpression()
	op := ctx.AssignmentOperator()

	switch {
	case e1 != nil:
		c.cgen.TransferValue(ctx, e1)

	case e2 != nil && e3 != nil && op.GetText() == "=":
		c.cgen.MoveValue(e2, e3)
		c.cgen.ReleaseValue(e3)
		c.cgen.TransferValue(ctx, e2)

	default:
		c.fail("AssignmentExpression: unhandled case")
	}
}

func (c *MainPass) EnterConditionalExpression(ctx *parser.ConditionalExpressionContext) {
	c.enterRule(ctx, "ConditionalExpression")
}

func (c *MainPass) ExitConditionalExpression(ctx *parser.ConditionalExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	c.debugf("LogicalOrExpression: %v\n", ctx.LogicalOrExpression() != nil)

	exp1 := ctx.LogicalOrExpression()
	c.cgen.TransferValue(ctx, exp1)
}

func (c *MainPass) EnterLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) {
	c.enterRule(ctx, "LogicalOrExpression")
}

func (c *MainPass) ExitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	exp1 := ctx.LogicalAndExpression()
	exp2 := ctx.LogicalOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.TransferValue(ctx, exp1)

	default:
		c.fail("ExitLogicalOrExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) {
	c.enterRule(ctx, "LogicalAndExpression")
}

func (c *MainPass) ExitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	exp1 := ctx.InclusiveOrExpression()
	exp2 := ctx.LogicalAndExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.TransferValue(ctx, exp1)

	default:
		c.fail("ExitLogicalAndExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterInclusiveOrExpression(ctx *parser.InclusiveOrExpressionContext) {
	c.enterRule(ctx, "InclusiveOrExpression")
}

func (c *MainPass) ExitInclusiveOrExpression(ctx *parser.InclusiveOrExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	exp1 := ctx.ExclusiveOrExpression()
	exp2 := ctx.InclusiveOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.TransferValue(ctx, exp1)

	default:
		c.fail("ExitInclusiveOrExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterExclusiveOrExpression(ctx *parser.ExclusiveOrExpressionContext) {
	c.enterRule(ctx, "ExclusiveOrExpression")
}

func (c *MainPass) ExitExclusiveOrExpression(ctx *parser.ExclusiveOrExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	exp1 := ctx.AndExpression()
	exp2 := ctx.ExclusiveOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.TransferValue(ctx, exp1)

	default:
		c.fail("ExitExclusiveOrExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterAndExpression(ctx *parser.AndExpressionContext) {
	c.enterRule(ctx, "AndExpression")
}

func (c *MainPass) ExitAndExpression(ctx *parser.AndExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	exp1 := ctx.EqualityExpression()
	exp2 := ctx.AndExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.TransferValue(ctx, exp1)

	default:
		c.fail("ExitAndExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterEqualityExpression(ctx *parser.EqualityExpressionContext) {
	c.enterRule(ctx, "EqualityExpression")
}

func (c *MainPass) ExitEqualityExpression(ctx *parser.EqualityExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	e1 := ctx.EqualityExpression()
	e2 := ctx.RelationalExpression()
	eq := ctx.Equal()
	neq := ctx.NotEqual()

	switch {
	case e1 == nil && e2 != nil:
		c.cgen.TransferValue(ctx, e2)

	case eq != nil:
		v1 := c.cgen.GetValue(e1)
		v2 := c.cgen.GetValue(e2)
		c.cgen.Equal(ctx, v1, v2)
		c.cgen.ReleaseValue(e1)
		c.cgen.ReleaseValue(e2)

	case neq != nil:
		v1 := c.cgen.GetValue(e1)
		v2 := c.cgen.GetValue(e2)
		c.cgen.NotEqual(ctx, v1, v2)
		c.cgen.ReleaseValue(e1)
		c.cgen.ReleaseValue(e2)

	default:
		c.fail("ExitEqualityExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterRelationalExpression(ctx *parser.RelationalExpressionContext) {
	c.enterRule(ctx, "RelationalExpression")
}

func (c *MainPass) ExitRelationalExpression(ctx *parser.RelationalExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	exp1 := ctx.ShiftExpression()
	exp2 := ctx.RelationalExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.TransferValue(ctx, exp1)

	default:
		c.fail("ExitRelationalExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterShiftExpression(ctx *parser.ShiftExpressionContext) {
	c.enterRule(ctx, "ShiftExpression")
}

func (c *MainPass) ExitShiftExpression(ctx *parser.ShiftExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	exp1 := ctx.AdditiveExpression()
	exp2 := ctx.ShiftExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.TransferValue(ctx, exp1)

	default:
		c.fail("ExitShiftExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	c.enterRule(ctx, "AdditiveExpression")
}

func (c *MainPass) ExitAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	e1 := ctx.AdditiveExpression()
	e2 := ctx.MultiplicativeExpression()
	plus := ctx.Plus()
	minus := ctx.Minus()

	c.debugf("e1: %v\n", e1)
	c.debugf("e2: %v\n", e2)
	c.debugf("plus: %v\n", plus)
	c.debugf("minus: %v\n", minus)

	switch {
	case e1 == nil && e2 != nil:
		c.cgen.TransferValue(ctx, e2)

	case e1 != nil && e2 != nil && plus != nil:
		v1 := c.cgen.GetValue(e1)
		v2 := c.cgen.GetValue(e2)
		c.cgen.Add(ctx, v1, v2)
		c.cgen.ReleaseValue(e1)
		c.cgen.ReleaseValue(e2)

	case e1 != nil && e2 != nil && minus != nil:
		v1 := c.cgen.GetValue(e1)
		v2 := c.cgen.GetValue(e2)
		c.cgen.Subtract(ctx, v1, v2)
		c.cgen.ReleaseValue(e1)
		c.cgen.ReleaseValue(e2)

	default:
		c.fail("ExitAdditiveExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) {
	c.enterRule(ctx, "MultiplicativeExpression")
}

func (c *MainPass) ExitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	e1 := ctx.MultiplicativeExpression()
	e2 := ctx.CastExpression()
	mult := ctx.Star() != nil
	div := ctx.Div() != nil
	mod := ctx.Mod() != nil

	switch {
	case e1 == nil && e2 != nil:
		c.cgen.TransferValue(ctx, e2)

	case e1 != nil && e2 != nil:
		v1 := c.cgen.GetValue(e1)
		v2 := c.cgen.GetValue(e2)

		switch {
		case mult:
			c.cgen.Multiply(ctx, v1, v2)

		case div:
			c.cgen.Divide(ctx, v1, v2)

		case mod:
			c.cgen.Modulus(ctx, v1, v2)

		default:
			c.fail("ExitMultiplicativeExpression(): unhandled sub-case!")
		}
		c.cgen.ReleaseValue(e1)
		c.cgen.ReleaseValue(e2)

	default:
		c.fail("ExitMultiplicativeExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterCastExpression(ctx *parser.CastExpressionContext) {
	c.enterRule(ctx, "CastExpression")
}

func (c *MainPass) ExitCastExpression(ctx *parser.CastExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	c.debugf("TypeName: %v\n", ctx.TypeName() != nil)
	c.debugf("CastExpression: %v\n", ctx.CastExpression() != nil)
	c.debugf("UnaryExpression: %v\n", ctx.UnaryExpression() != nil)

	exp1 := ctx.TypeName()
	exp2 := ctx.CastExpression()
	exp3 := ctx.UnaryExpression()

	switch {
	case exp1 == nil && exp2 == nil && exp3 != nil:
		c.cgen.TransferValue(ctx, exp3)

	default:
		c.fail("ExitCastExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterUnaryExpression(ctx *parser.UnaryExpressionContext) {
	c.enterRule(ctx, "UnaryExpression")
}

func (c *MainPass) ExitUnaryExpression(ctx *parser.UnaryExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	e1 := ctx.PostfixExpression()
	e2 := ctx.UnaryOperator()
	e3 := ctx.CastExpression()

	switch {
	case e1 != nil:
		c.cgen.TransferValue(ctx, e1)

	case e2 != nil && e3 != nil:
		op := e2.(*parser.UnaryOperatorContext)
		switch {
		case op.Minus() != nil:
			val := c.cgen.GetValue(e3)
			c.cgen.UnaryMinus(ctx, val)
			c.cgen.ReleaseValue(e3)

		case op.Plus() != nil:
			// nothing to do
			c.cgen.TransferValue(ctx, e3)

		default:
			c.fail("ExitUnaryExpression(): unhandled sub-case")
		}

	default:
		c.fail("ExitUnaryExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterPostfixExpression(ctx *parser.PostfixExpressionContext) {
	c.enterRule(ctx, "PostfixExpression")
}

func (c *MainPass) ExitPostfixExpression(ctx *parser.PostfixExpressionContext) {
	e := c.exitRule(ctx)
	defer e()

	exp1 := ctx.PrimaryExpression()

	switch {
	case exp1 != nil:
		c.cgen.TransferValue(ctx, exp1)

	default:
		c.fail("ExitPostfixExpression(): unhandled case!")
	}
}

func (c *MainPass) EnterPrimaryExpressionIdentifier(ctx *parser.PrimaryExpressionIdentifierContext) {
	c.enterRule(ctx, "PrimaryExpressionIdentifier")
}

func (c *MainPass) ExitPrimaryExpressionIdentifier(ctx *parser.PrimaryExpressionIdentifierContext) {
	e := c.exitRule(ctx)
	defer e()

	// An identifier being evaluated.  Look up the symbol.
	name := ctx.Identifier().GetText()
	_, value, ok := c.Symbols.FindSymbol(name)
	if !ok {
		c.fail("unknown identifier: %s", name)
	}

	// forward the value
	c.cgen.CreateValue(ctx, value)
}

func (c *MainPass) EnterPrimaryExpressionConstant(ctx *parser.PrimaryExpressionConstantContext) {
	c.enterRule(ctx, "PrimaryExpressionConstant")
}

func (c *MainPass) ExitPrimaryExpressionConstant(ctx *parser.PrimaryExpressionConstantContext) {
	e := c.exitRule(ctx)
	defer e()

	cp := ctx.GetParser().(*parser.CParser)
	tok := ctx.GetStart()
	typ := types.NewBasicType(parser.CParserInt, cp) // TODO: type
	c.cgen.CreateIntLiteralValue(ctx, typ, tok.GetText())
}

func (c *MainPass) EnterPrimaryExpressionStringLiteral(ctx *parser.PrimaryExpressionStringLiteralContext) {
	c.enterRule(ctx, "PrimaryExpressionStringLiteral")
}

func (c *MainPass) ExitPrimaryExpressionStringLiteral(ctx *parser.PrimaryExpressionStringLiteralContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterPrimaryExpressionParens(ctx *parser.PrimaryExpressionParensContext) {
	c.enterRule(ctx, "PrimaryExpressionParens")
}

func (c *MainPass) ExitPrimaryExpressionParens(ctx *parser.PrimaryExpressionParensContext) {
	e := c.exitRule(ctx)
	defer e()

	// forward the value
	exp := ctx.Expression()
	c.cgen.TransferValue(ctx, exp)
}

func (c *MainPass) EnterPrimaryExpressionExtension(ctx *parser.PrimaryExpressionExtensionContext) {
	c.enterRule(ctx, "PrimaryExpressionExtension")
}

func (c *MainPass) ExitPrimaryExpressionExtension(ctx *parser.PrimaryExpressionExtensionContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterPrimaryExpressionGeneric(ctx *parser.PrimaryExpressionGenericContext) {
	c.enterRule(ctx, "PrimaryExpressionGeneric")
}

func (c *MainPass) ExitPrimaryExpressionGeneric(ctx *parser.PrimaryExpressionGenericContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterPrimaryExpressionOffsetOf(ctx *parser.PrimaryExpressionOffsetOfContext) {
	c.enterRule(ctx, "PrimaryExpressionOffsetOf")
}

func (c *MainPass) ExitPrimaryExpressionOffsetOf(ctx *parser.PrimaryExpressionOffsetOfContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterPrimaryExpressionVaArg(ctx *parser.PrimaryExpressionVaArgContext) {
	c.enterRule(ctx, "PrimaryExpressionVaArg")
}

func (c *MainPass) ExitPrimaryExpressionVaArg(ctx *parser.PrimaryExpressionVaArgContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterAssignmentOperator(ctx *parser.AssignmentOperatorContext) {
	c.enterRule(ctx, "AssignmentOperator")
}

func (c *MainPass) ExitAssignmentOperator(ctx *parser.AssignmentOperatorContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterForExpression(ctx *parser.ForExpressionContext) {
	c.enterRule(ctx, "ForExpression")
}

func (c *MainPass) ExitForExpression(ctx *parser.ForExpressionContext) {
	e := c.exitRule(ctx)
	defer e()
}
