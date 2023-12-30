package mainpass

import (
	"github.com/TassoKarkanis/minic/parser"
)

func (c *MainPass) EnterCompoundStatement(ctx *parser.CompoundStatementContext) {
	c.enterRule(ctx, "CompoundStatement")

	c.Symbols.PushScope()
}

func (c *MainPass) ExitCompoundStatement(ctx *parser.CompoundStatementContext) {
	e := c.exitRule(ctx)
	defer e()

	c.Symbols.PopScope()
}

func (c *MainPass) EnterStatement(ctx *parser.StatementContext) {
	c.enterRule(ctx, "Statement")
}

func (c *MainPass) ExitStatement(ctx *parser.StatementContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterExpressionStatement(ctx *parser.ExpressionStatementContext) {
	c.enterRule(ctx, "ExpressionStatement")
}

func (c *MainPass) ExitExpressionStatement(ctx *parser.ExpressionStatementContext) {
	e := c.exitRule(ctx)
	defer e()

	e1 := ctx.Expression()
	if e1 != nil {
		c.cgen.ReleaseValue(e1)
	}
}

func (c *MainPass) EnterIterationStatement(ctx *parser.IterationStatementContext) {
	c.enterRule(ctx, "IterationStatement")
}

func (c *MainPass) ExitIterationStatement(ctx *parser.IterationStatementContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterLabeledStatement(ctx *parser.LabeledStatementContext) {
	c.enterRule(ctx, "LabeledStatement")
}

func (c *MainPass) ExitLabeledStatement(ctx *parser.LabeledStatementContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterSelectionStatement(ctx *parser.SelectionStatementContext) {
	c.enterRule(ctx, "SelectionStatement")
}

func (c *MainPass) ExitSelectionStatement(ctx *parser.SelectionStatementContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterJumpStatement(ctx *parser.JumpStatementContext) {
	c.enterRule(ctx, "JumpStatement")
}

func (c *MainPass) ExitJumpStatement(ctx *parser.JumpStatementContext) {
	e := c.exitRule(ctx)
	defer e()

	switch {
	case ctx.Return() != nil:
		expr := ctx.Expression()
		if expr != nil {
			c.cgen.ReturnValue(expr)
			c.cgen.ReleaseValue(expr)
		} else {
			c.cgen.Return()
		}
	}
}
