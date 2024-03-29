package mainpass

import (
	"github.com/TassoKarkanis/minic/codegen"
	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/types"
)

func (c *MainPass) EnterParameterTypeList(ctx *parser.ParameterTypeListContext) {
	c.enterRule(ctx, "ParameterTypeList")
}

func (c *MainPass) ExitParameterTypeList(ctx *parser.ParameterTypeListContext) {
	e := c.exitRule(ctx)
	defer e()

	paramList := ctx.ParameterList()
	ellipsis := ctx.Ellipsis()

	switch {
	case paramList != nil && ellipsis == nil:
		params := c.getType(paramList).(*types.FunctionType)
		c.setType(ctx, &types.FunctionType{
			Params: params.Params,
		})

	default:
		c.fail("varargs not supported")
	}

	c.debugf("result: %+v\n", c.getType(ctx))
}

func (c *MainPass) EnterParameterList(ctx *parser.ParameterListContext) {
	c.enterRule(ctx, "ParameterList")
}

func (c *MainPass) ExitParameterList(ctx *parser.ParameterListContext) {
	e := c.exitRule(ctx)
	defer e()

	paramDecl := ctx.ParameterDeclaration()
	paramList := ctx.ParameterList()

	switch {
	case paramDecl != nil && paramList == nil:
		param := c.Declarations[paramDecl]
		c.setType(ctx, &types.FunctionType{
			Params: []types.Param{
				{
					Name: param.Name,
					Type: param.Type,
				},
			},
		})

	case paramList != nil && paramDecl != nil:
		params := c.getType(paramList).(*types.FunctionType)
		param := c.Declarations[paramDecl]
		c.setType(ctx, &types.FunctionType{
			Params: append(params.Params, types.Param{
				Name: param.Name,
				Type: param.Type,
			}),
		})

	default:
		c.fail("ExitParameterList: invalid case")
	}

	c.debugf("result: %+v\n", c.getType(ctx))
}

func (c *MainPass) EnterParameterDeclaration(ctx *parser.ParameterDeclarationContext) {
	c.enterRule(ctx, "ParameterDeclaration")
}

func (c *MainPass) ExitParameterDeclaration(ctx *parser.ParameterDeclarationContext) {
	e := c.exitRule(ctx)
	defer e()

	declSpec := ctx.DeclarationSpecifiers()
	declarator := ctx.Declarator()
	cp := ctx.GetParser().(*parser.CParser)

	switch {
	case declSpec != nil && declarator != nil:
		typ := c.getType(declSpec)
		decl := c.Declarations[declarator]
		c.debugf("typ: %s\n", typ)
		c.debugf("decl: %s\n", decl)

		if decl.Type == nil {
			c.Declarations[ctx] = Declaration{
				Name: decl.Name,
				Type: typ,
			}
		} else {
			f := decl.Type.(*types.FunctionType)
			c.Declarations[ctx] = Declaration{
				Name: decl.Name,
				Type: &types.FunctionType{
					Name:       f.Name,
					ReturnType: typ,
					Params:     f.Params,
				},
			}
		}

	case declSpec == nil && declarator == nil:
		// "void", I guess
		c.Declarations[ctx] = Declaration{
			Name: "",
			Type: types.NewBasicType(parser.CParserVoid, cp),
		}

	default:
		c.fail("ExitParameterDeclaration: invalid case")
	}

	c.debugf("result: %+v\n", c.Declarations[ctx])
}

func (c *MainPass) EnterFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	c.enterRule(ctx, "FunctionDefinition")

	if c.Function != nil {
		c.fail("nested function not implemented")
	}

	c.Function = &types.FunctionType{}

	// make sure the return type is set
	{
		if ctx.DeclarationSpecifiers() == nil {
			c.fail("return type not specified")
		}
		declSpecs := ctx.DeclarationSpecifiers().(*parser.DeclarationSpecifiersContext)
		if declSpecs.DeclarationSpecifier(1) != nil {
			c.fail("invalid return types")
		}

		declSpec := declSpecs.DeclarationSpecifier(0).(*parser.DeclarationSpecifierContext)
		c.setExitContinuation(declSpec, func() {
			c.Function.ReturnType = c.getType(declSpec)
		})
	}

	// make sure the function name and parameters are set
	{
		declarator := ctx.Declarator()
		c.setExitContinuation(declarator, func() {
			decl := c.Declarations[declarator]
			c.Function.Name = decl.Name
			c.Function.Params = decl.Type.(*types.FunctionType).Params
		})
	}

	// create the code generator
	c.cgen = codegen.NewCodegen(c.Output)

	// push the scope for the parameters
	c.Symbols.PushScope()

	// generate the function prologue on entry of the statements
	block := ctx.CompoundStatement()
	c.setEnterContinuation(block, func() {
		c.debugf("  starting function: %s\n", c.Function)

		// fix functions with "void" for argument
		if len(c.Function.Params) == 1 {
			typ := c.Function.Params[0].Type
			basicType, ok := typ.(*types.BasicType)
			if ok && basicType.Type == parser.CParserVoid {
				c.Function.Params = nil
			}
		}

		// get the parameter names and types
		var names []string
		var types []types.Type
		for _, param := range c.Function.Params {
			names = append(names, param.Name)
			types = append(types, param.Type)
		}

		values := c.cgen.StartStackFrame(c.Function.Name, names, types)

		// add symbols for the parameters
		for i, param := range c.Function.Params {
			c.Symbols.AddSymbol(param.Name, param.Type, values[i])
		}

		// add the global symbol for the function
		funcValue := c.cgen.CreateFunction(c.Function.Name, c.Function)
		c.Symbols.AddSymbol(c.Function.Name, c.Function, funcValue)
	})
}

func (c *MainPass) ExitFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	e := c.exitRule(ctx)
	defer e()

	c.cgen.EndStackFrame()
	c.Symbols.PopScope()
	c.cgen.Close()
	c.cgen = nil

	c.Function = nil
}
