package mainpass

import (
	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/types"
)

func (c *MainPass) EnterExpression(ctx *parser.ExpressionContext) {
	c.enterf("Expression", "%s", ctx.GetText())
}

func (c *MainPass) ExitExpression(ctx *parser.ExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.AssignmentExpression()
	exp2 := ctx.Expression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterConstantExpression(ctx *parser.ConstantExpressionContext) {
	c.enterf("ConstantExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitConstantExpression(ctx *parser.ConstantExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterAssignmentExpression(ctx *parser.AssignmentExpressionContext) {
	c.enterf("AssignmentExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitAssignmentExpression(ctx *parser.AssignmentExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.ConditionalExpression()

	switch {
	case exp1 != nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterConditionalExpression(ctx *parser.ConditionalExpressionContext) {
	c.enterf("ConditionalExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitConditionalExpression(ctx *parser.ConditionalExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	c.debugf("LogicalOrExpression: %v\n", ctx.LogicalOrExpression() != nil)

	exp1 := ctx.LogicalOrExpression()
	c.cgen.MoveValue(ctx, exp1)
}

func (c *MainPass) EnterLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) {
	c.enterf("LogicalOrExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.LogicalAndExpression()
	exp2 := ctx.LogicalOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) {
	c.enterf("LogicalAndExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.InclusiveOrExpression()
	exp2 := ctx.LogicalAndExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterInclusiveOrExpression(ctx *parser.InclusiveOrExpressionContext) {
	c.enterf("InclusiveOrExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitInclusiveOrExpression(ctx *parser.InclusiveOrExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.ExclusiveOrExpression()
	exp2 := ctx.InclusiveOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterExclusiveOrExpression(ctx *parser.ExclusiveOrExpressionContext) {
	c.enterf("ExclusiveOrExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitExclusiveOrExpression(ctx *parser.ExclusiveOrExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.AndExpression()
	exp2 := ctx.ExclusiveOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterAndExpression(ctx *parser.AndExpressionContext) {
	c.enterf("AndExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitAndExpression(ctx *parser.AndExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.EqualityExpression()
	exp2 := ctx.AndExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterEqualityExpression(ctx *parser.EqualityExpressionContext) {
	c.enterf("EqualityExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitEqualityExpression(ctx *parser.EqualityExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.RelationalExpression()
	exp2 := ctx.EqualityExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterRelationalExpression(ctx *parser.RelationalExpressionContext) {
	c.enterf("RelationalExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitRelationalExpression(ctx *parser.RelationalExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.ShiftExpression()
	exp2 := ctx.RelationalExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterShiftExpression(ctx *parser.ShiftExpressionContext) {
	c.enterf("ShiftExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitShiftExpression(ctx *parser.ShiftExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.AdditiveExpression()
	exp2 := ctx.ShiftExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	c.enterf("AdditiveExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
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
		c.cgen.MoveValue(ctx, e2)

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
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) {
	c.enterf("MultiplicativeExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	e1 := ctx.CastExpression()
	e2 := ctx.MultiplicativeExpression()
	mult := ctx.Star() != nil
	div := ctx.Div() != nil
	mod := ctx.Mod() != nil

	switch {
	case e1 != nil && e2 == nil:
		c.cgen.MoveValue(ctx, e1)

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
			c.fail("unhandled case!")
		}
		c.cgen.ReleaseValue(e1)
		c.cgen.ReleaseValue(e2)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterCastExpression(ctx *parser.CastExpressionContext) {
	c.enterf("CastExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitCastExpression(ctx *parser.CastExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	c.debugf("TypeName: %v\n", ctx.TypeName() != nil)
	c.debugf("CastExpression: %v\n", ctx.CastExpression() != nil)
	c.debugf("UnaryExpression: %v\n", ctx.UnaryExpression() != nil)

	exp1 := ctx.TypeName()
	exp2 := ctx.CastExpression()
	exp3 := ctx.UnaryExpression()

	switch {
	case exp1 == nil && exp2 == nil && exp3 != nil:
		c.cgen.MoveValue(ctx, exp3)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterUnaryExpression(ctx *parser.UnaryExpressionContext) {
	c.enterf("UnaryExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitUnaryExpression(ctx *parser.UnaryExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.PostfixExpression()

	switch {
	case exp1 != nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterPostfixExpression(ctx *parser.PostfixExpressionContext) {
	c.enterf("PostfixExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitPostfixExpression(ctx *parser.PostfixExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	exp1 := ctx.PrimaryExpression()

	switch {
	case exp1 != nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		c.fail("unhandled case!")
	}
}

func (c *MainPass) EnterPrimaryExpressionIdentifier(ctx *parser.PrimaryExpressionIdentifierContext) {
	c.enterf("PrimaryExpressionIdentifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitPrimaryExpressionIdentifier(ctx *parser.PrimaryExpressionIdentifierContext) {
	e := c.exitf("%s", ctx.GetText())
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
	c.enterf("PrimaryExpressionConstant", "%s", ctx.GetText())
}

func (c *MainPass) ExitPrimaryExpressionConstant(ctx *parser.PrimaryExpressionConstantContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	cp := ctx.GetParser().(*parser.CParser)
	tok := ctx.GetStart()
	typ := types.NewBasicType(parser.CParserInt, cp) // TODO: type
	c.cgen.CreateIntLiteralValue(ctx, typ, tok.GetText())
}

func (c *MainPass) EnterPrimaryExpressionStringLiteral(ctx *parser.PrimaryExpressionStringLiteralContext) {
	c.enterf("PrimaryExpressionStringLiteral", "%s", ctx.GetText())
}

func (c *MainPass) ExitPrimaryExpressionStringLiteral(ctx *parser.PrimaryExpressionStringLiteralContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterPrimaryExpressionParens(ctx *parser.PrimaryExpressionParensContext) {
	c.enterf("PrimaryExpressionParens", "%s", ctx.GetText())
}

func (c *MainPass) ExitPrimaryExpressionParens(ctx *parser.PrimaryExpressionParensContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	// forward the value
	exp := ctx.Expression()
	c.cgen.MoveValue(ctx, exp)
}

func (c *MainPass) EnterPrimaryExpressionExtension(ctx *parser.PrimaryExpressionExtensionContext) {
	c.enterf("PrimaryExpressionExtension", "%s", ctx.GetText())
}

func (c *MainPass) ExitPrimaryExpressionExtension(ctx *parser.PrimaryExpressionExtensionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterPrimaryExpressionGeneric(ctx *parser.PrimaryExpressionGenericContext) {
	c.enterf("PrimaryExpressionGeneric", "%s", ctx.GetText())
}

func (c *MainPass) ExitPrimaryExpressionGeneric(ctx *parser.PrimaryExpressionGenericContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterPrimaryExpressionOffsetOf(ctx *parser.PrimaryExpressionOffsetOfContext) {
	c.enterf("PrimaryExpressionOffsetOf", "%s", ctx.GetText())
}

func (c *MainPass) ExitPrimaryExpressionOffsetOf(ctx *parser.PrimaryExpressionOffsetOfContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterPrimaryExpressionVaArg(ctx *parser.PrimaryExpressionVaArgContext) {
	c.enterf("PrimaryExpressionVaArg", "%s", ctx.GetText())
}

func (c *MainPass) ExitPrimaryExpressionVaArg(ctx *parser.PrimaryExpressionVaArgContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterAssignmentOperator(ctx *parser.AssignmentOperatorContext) {
	c.enterf("AssignmentOperator", "%s", ctx.GetText())
}

func (c *MainPass) ExitAssignmentOperator(ctx *parser.AssignmentOperatorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterForExpression(ctx *parser.ForExpressionContext) {
	c.enterf("ForExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitForExpression(ctx *parser.ForExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
