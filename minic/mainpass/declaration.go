package mainpass

import (
	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/types"
)

func (c *MainPass) EnterInitDeclarator(ctx *parser.InitDeclaratorContext) {
	c.enterRule(ctx, "InitDeclarator")
}

func (c *MainPass) ExitInitDeclarator(ctx *parser.InitDeclaratorContext) {
	e := c.exitRule(ctx)
	defer e()

	// get the declaration
	e1 := ctx.Declarator()
	decl := c.Declarations[e1]

	// set the type
	decl.Type = c.getType(ctx)

	// create a symbol
	val := c.cgen.CreateLocalValue(decl.Name, decl.Type)
	c.Symbols.AddSymbol(decl.Name, decl.Type, val)

	// assign the initializer, if any
	e2 := ctx.Initializer()
	if e2 != nil {
		c.cgen.Move(val, e2)
		c.cgen.ReleaseValue(e2)
	}
}

func (c *MainPass) EnterInitDeclaratorList(ctx *parser.InitDeclaratorListContext) {
	c.enterRule(ctx, "InitDeclaratorList")

	e1 := ctx.InitDeclarator()
	e2 := ctx.InitDeclaratorList()

	// push the type downward
	if e1 != nil {
		c.copyType(e1, ctx)
	}
	if e2 != nil {
		c.copyType(e2, ctx)
	}
}

func (c *MainPass) ExitInitDeclaratorList(ctx *parser.InitDeclaratorListContext) {
	e := c.exitRule(ctx)
	defer e()

	if ctx.InitDeclaratorList() != nil {
		c.fail("multi-declarations not yet supported")
	}
}

func (c *MainPass) EnterInitializer(ctx *parser.InitializerContext) {
	c.enterRule(ctx, "Initializer")
}

func (c *MainPass) ExitInitializer(ctx *parser.InitializerContext) {
	e := c.exitRule(ctx)
	defer e()

	e1 := ctx.AssignmentExpression()
	if e1 == nil {
		c.fail("complex initializers not supported")
	}

	c.cgen.TransferValue(ctx, e1)
}

func (c *MainPass) EnterInitializerList(ctx *parser.InitializerListContext) {
	c.enterRule(ctx, "InitializerList")
}

func (c *MainPass) ExitInitializerList(ctx *parser.InitializerListContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterDeclaration(ctx *parser.DeclarationContext) {
	c.enterRule(ctx, "Declaration")

	e1 := ctx.DeclarationSpecifiers()
	e2 := ctx.InitDeclaratorList()

	if e1 == nil || e2 == nil {
		c.fail("EnterDeclaration(): unsupported")
	}

	f := func() {
		// push the type downard
		c.copyType(e2, e1)
	}
	c.setEnterContinuation(e2, f)
}

func (c *MainPass) ExitDeclaration(ctx *parser.DeclarationContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterExternalDeclaration(ctx *parser.ExternalDeclarationContext) {
	c.enterRule(ctx, "ExternalDeclaration")
}

func (c *MainPass) ExitExternalDeclaration(ctx *parser.ExternalDeclarationContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterDeclarationSpecifiers(ctx *parser.DeclarationSpecifiersContext) {
	c.enterRule(ctx, "DeclarationSpecifiers")
}

func (c *MainPass) ExitDeclarationSpecifiers(ctx *parser.DeclarationSpecifiersContext) {
	e := c.exitRule(ctx)
	defer e()

	if ctx.DeclarationSpecifier(1) != nil {
		c.fail("multiple declaration specifiers!")
	}

	// forward
	c.copyType(ctx, ctx.DeclarationSpecifier(0))

	c.debugf("  result: %s\n", c.getType(ctx))
}

func (c *MainPass) EnterDeclarationSpecifier(ctx *parser.DeclarationSpecifierContext) {
	c.enterRule(ctx, "DeclarationSpecifier")
}

func (c *MainPass) ExitDeclarationSpecifier(ctx *parser.DeclarationSpecifierContext) {
	e := c.exitRule(ctx)
	defer e()

	typeSpec := ctx.TypeSpecifier()
	switch {
	case typeSpec != nil:
		// forward the type
		c.copyType(ctx, typeSpec)

	default:
		c.fail("unsupported type")
	}

	c.debugf("  result: %s", c.getType(ctx))
}

func (c *MainPass) EnterDeclarator(ctx *parser.DeclaratorContext) {
	c.enterRule(ctx, "Declarator")
}

func (c *MainPass) ExitDeclarator(ctx *parser.DeclaratorContext) {
	e := c.exitRule(ctx)
	defer e()

	if ctx.Pointer() != nil {
		c.fail("pointers not yet supported")
	}

	// forward
	c.Declarations[ctx] = c.Declarations[ctx.DirectDeclarator()]
}

func (c *MainPass) EnterAbstractDeclarator(ctx *parser.AbstractDeclaratorContext) {
	c.enterRule(ctx, "AbstractDeclarator")
}

func (c *MainPass) ExitAbstractDeclarator(ctx *parser.AbstractDeclaratorContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterDirectDeclarator(ctx *parser.DirectDeclaratorContext) {
	c.enterRule(ctx, "DirectDeclarator")
}

func (c *MainPass) ExitDirectDeclarator(ctx *parser.DirectDeclaratorContext) {
	e := c.exitRule(ctx)
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
		f := c.getType(params).(*types.FunctionType)
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

func (c *MainPass) EnterDeclarationList(ctx *parser.DeclarationListContext) {
	c.enterRule(ctx, "DeclarationList")
}

func (c *MainPass) ExitDeclarationList(ctx *parser.DeclarationListContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterDeclarationSpecifiers2(ctx *parser.DeclarationSpecifiers2Context) {
	c.enterRule(ctx, "DeclarationSpecifiers2")
}

func (c *MainPass) ExitDeclarationSpecifiers2(ctx *parser.DeclarationSpecifiers2Context) {
	e := c.exitRule(ctx)
	defer e()
}
