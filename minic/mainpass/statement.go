package mainpass

import (
	"github.com/TassoKarkanis/minic/parser"
)

func (c *MainPass) EnterStatement(ctx *parser.StatementContext) {
	c.enterf("Statement", "%s", ctx.GetText())
}

func (c *MainPass) ExitStatement(ctx *parser.StatementContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterExpressionStatement(ctx *parser.ExpressionStatementContext) {
	c.enterf("ExpressionStatement", "%s", ctx.GetText())
}

func (c *MainPass) ExitExpressionStatement(ctx *parser.ExpressionStatementContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterIterationStatement(ctx *parser.IterationStatementContext) {
	c.enterf("IterationStatement", "%s", ctx.GetText())
}

func (c *MainPass) ExitIterationStatement(ctx *parser.IterationStatementContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterLabeledStatement(ctx *parser.LabeledStatementContext) {
	c.enterf("LabeledStatement", "%s", ctx.GetText())
}

func (c *MainPass) ExitLabeledStatement(ctx *parser.LabeledStatementContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterSelectionStatement(ctx *parser.SelectionStatementContext) {
	c.enterf("SelectionStatement", "%s", ctx.GetText())
}

func (c *MainPass) ExitSelectionStatement(ctx *parser.SelectionStatementContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterJumpStatement(ctx *parser.JumpStatementContext) {
	c.enterf("JumpStatement", "%s", ctx.GetText())
}

func (c *MainPass) ExitJumpStatement(ctx *parser.JumpStatementContext) {
	e := c.exitf("%s", ctx.GetText())
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
