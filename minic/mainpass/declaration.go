package mainpass

import (
	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/types"
)

func (c *MainPass) EnterDeclaration(ctx *parser.DeclarationContext) {
	c.enterf("Declaration", "%s", ctx.GetText())
}

func (c *MainPass) ExitDeclaration(ctx *parser.DeclarationContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterExternalDeclaration(ctx *parser.ExternalDeclarationContext) {
	c.enterf("ExternalDeclaration", ctx.GetText())
}

func (c *MainPass) ExitExternalDeclaration(ctx *parser.ExternalDeclarationContext) {
	e := c.exitf("")
	defer e()
}

func (c *MainPass) EnterDeclarationSpecifiers(ctx *parser.DeclarationSpecifiersContext) {
	c.enterf("DeclarationSpecifiers", ctx.GetText())
	c.runEnterContinuation(ctx)
}

func (c *MainPass) ExitDeclarationSpecifiers(ctx *parser.DeclarationSpecifiersContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	if ctx.DeclarationSpecifier(1) != nil {
		c.fail("multiple declaration specifiers!")
	}

	// forward
	c.Types[ctx] = c.Types[ctx.DeclarationSpecifier(0)]

	c.debugf("  result: %s\n", c.Types[ctx])

	c.runExitContinuation(ctx)
}

func (c *MainPass) EnterDeclarationSpecifier(ctx *parser.DeclarationSpecifierContext) {
	c.enterf("DeclarationSpecifier", ctx.GetText())
	c.runEnterContinuation(ctx)
}

func (c *MainPass) ExitDeclarationSpecifier(ctx *parser.DeclarationSpecifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	typeSpec := ctx.TypeSpecifier()
	switch {
	case typeSpec != nil:
		// forward the type
		c.Types[ctx] = c.Types[typeSpec]

	default:
		c.fail("unsupported type")
	}

	c.debugf("  result: %s\n", c.Types[ctx])

	c.runExitContinuation(ctx)
}

func (c *MainPass) EnterDeclarator(ctx *parser.DeclaratorContext) {
	c.enterf("Declarator", "%s", ctx.GetText())
	c.runEnterContinuation(ctx)
}

func (c *MainPass) ExitDeclarator(ctx *parser.DeclaratorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	if ctx.Pointer() != nil {
		c.fail("pointers not yet supported")
	}

	// forward
	c.Declarations[ctx] = c.Declarations[ctx.DirectDeclarator()]

	c.runExitContinuation(ctx)
}

func (c *MainPass) EnterAbstractDeclarator(ctx *parser.AbstractDeclaratorContext) {
	c.enterf("AbstractDeclarator", "%s", ctx.GetText())
}

func (c *MainPass) ExitAbstractDeclarator(ctx *parser.AbstractDeclaratorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterDirectDeclarator(ctx *parser.DirectDeclaratorContext) {
	c.enterf("DirectDeclarator", "%s", ctx.GetText())
}

func (c *MainPass) ExitDirectDeclarator(ctx *parser.DirectDeclaratorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	ident := ctx.Identifier()
	directDecl := ctx.DirectDeclarator()
	colon := ctx.Colon()
	params := ctx.ParameterTypeList()

	switch {
	case ident != nil && colon == nil:
		c.Declarations[ctx] = Declaration{
			Name: ident.GetSymbol().GetText(),
		}

	case directDecl != nil && params != nil:
		decl := c.Declarations[directDecl]
		f := c.Types[params].(*types.FunctionType)
		c.Declarations[ctx] = Declaration{
			Name: decl.Name,
			Type: &types.FunctionType{
				Name:   decl.Name,
				Params: f.Params,
			},
		}

	default:
		c.fail("unsupported direct declarator")
	}
}
