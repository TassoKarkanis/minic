package main

import (
	"fmt"
	"io"

	"github.com/TassoKarkanis/minic/codegen"
	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/symbols"
	"github.com/TassoKarkanis/minic/types"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type Declaration struct {
	Name string
	Type types.Type
}

func (d Declaration) String() string {
	name := "<none>"
	if d.Name != "" {
		name = d.Name
	}

	typeStr := "<none>"
	if d.Type != nil {
		typeStr = d.Type.String()
	}

	return fmt.Sprintf("%s %s", name, typeStr)
}

type MainPass struct {
	*parser.BaseCListener
	Output             io.Writer
	Symbols            *symbols.Table
	Function           *types.FunctionType
	cgen               *codegen.Codegen
	Err                error
	Types              map[antlr.ParserRuleContext]types.Type
	Declarations       map[antlr.ParserRuleContext]Declaration
	EnterContinuations map[antlr.ParserRuleContext]func()
	ExitContinuations  map[antlr.ParserRuleContext]func()
}

func NewMainPass(output io.Writer) *MainPass {
	return &MainPass{
		Output:             output,
		Symbols:            symbols.NewTable(),
		Types:              make(map[antlr.ParserRuleContext]types.Type),
		Declarations:       make(map[antlr.ParserRuleContext]Declaration),
		EnterContinuations: make(map[antlr.ParserRuleContext]func()),
		ExitContinuations:  make(map[antlr.ParserRuleContext]func()),
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

func (c *MainPass) ExitExternalDeclaration(ctx *parser.ExternalDeclarationContext) {
	fmt.Printf("ExitExternalDeclaration(): %s\n", ctx.GetText())
}

func (c *MainPass) EnterFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	fmt.Printf("EnterFunctionDefinition(): %s\n", ctx.GetText())

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
			c.Function.ReturnType = c.Types[declSpec]
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

	// generate the function prologue on entry of the statements
	c.setEnterContinuation(ctx.CompoundStatement(), func() {
		fmt.Printf("  starting function: %s\n", c.Function)

		// get the parameter types
		var types []types.Type
		for _, param := range c.Function.Params {
			types = append(types, param.Type)
		}

		values := c.cgen.StartStackFrame(c.Function.Name, types)

		// add symbols for the parameters
		for i, param := range c.Function.Params {
			c.Symbols.AddSymbol(param.Name, param.Type, values[i])
		}

		// add the global symbol for the function
		funcValue := codegen.NewGlobalValue(c.Function.Name, c.Function)
		c.Symbols.AddSymbol(c.Function.Name, c.Function, funcValue)

		// push the scope for the function
		c.Symbols.PushScope()
	})

	c.cgen = codegen.NewCodegen(c.Output)
}

func (c *MainPass) ExitFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	fmt.Printf("ExitFunctionDefinition(): %s\n", ctx.GetText())
	c.Symbols.PopScope()
	c.cgen.EndStackFrame()
	c.cgen.Close()

	c.Function = nil
}

func (c *MainPass) EnterDeclarationSpecifiers(ctx *parser.DeclarationSpecifiersContext) {
	fmt.Printf("EnterDeclarationsSpecifier(): %s\n", ctx.GetText())
	c.runEnterContinuation(ctx)
}

func (c *MainPass) ExitDeclarationSpecifiers(ctx *parser.DeclarationSpecifiersContext) {
	fmt.Printf("ExitDeclarationSpecifiers(): %s\n", ctx.GetText())

	if ctx.DeclarationSpecifier(1) != nil {
		c.fail("multiple declaration specifiers!")
	}

	// forward
	c.Types[ctx] = c.Types[ctx.DeclarationSpecifier(0)]

	fmt.Printf("  result: %s\n", c.Types[ctx])

	c.runExitContinuation(ctx)
}

func (c *MainPass) EnterDeclarationSpecifier(ctx *parser.DeclarationSpecifierContext) {
	fmt.Printf("EnterDeclarationSpecifier(): %s\n", ctx.GetText())
	c.runEnterContinuation(ctx)
}

func (c *MainPass) ExitDeclarationSpecifier(ctx *parser.DeclarationSpecifierContext) {
	fmt.Printf("ExitDeclarationSpecifier(): %s\n", ctx.GetText())

	typeSpec := ctx.TypeSpecifier()
	switch {
	case typeSpec != nil:
		// forward the type
		c.Types[ctx] = c.Types[typeSpec]

	default:
		c.fail("unsupported type")
	}

	fmt.Printf("  result: %s\n", c.Types[ctx])

	c.runExitContinuation(ctx)
}

func (c *MainPass) EnterDeclarator(ctx *parser.DeclaratorContext) {
	fmt.Printf("EnterDeclarator(): %s\n", ctx.GetText())
	c.runEnterContinuation(ctx)
}

func (c *MainPass) ExitDeclarator(ctx *parser.DeclaratorContext) {
	fmt.Printf("ExitDeclarator(): %s\n", ctx.GetText())

	if ctx.Pointer() != nil {
		c.fail("pointers not yet supported")
	}

	// forward
	c.Declarations[ctx] = c.Declarations[ctx.DirectDeclarator()]

	c.runExitContinuation(ctx)
}

func (c *MainPass) EnterDirectDeclarator(ctx *parser.DirectDeclaratorContext) {
	fmt.Printf("EnterDirectDeclarator(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitDirectDeclarator(ctx *parser.DirectDeclaratorContext) {
	fmt.Printf("ExitDirectDeclarator(): %s\n", ctx.GetText())

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

func (c *MainPass) EnterTypeSpecifier(ctx *parser.TypeSpecifierContext) {
	fmt.Printf("EnterTypeSpecifier(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitTypeSpecifier(ctx *parser.TypeSpecifierContext) {
	fmt.Printf("ExitTypeSpecifier(): %s\n", ctx.GetText())

	cp := ctx.GetParser().(*parser.CParser)

	typ_void := ctx.Void()

	typ_char := ctx.Char()
	typ_short := ctx.Short()
	typ_int := ctx.Int()
	typ_long := ctx.Long()

	typ_float := ctx.Float()
	typ_double := ctx.Double()

	setBasicType := func(parserType int) {
		typ := types.NewBasicType(parserType, cp)
		fmt.Printf("  setBasicType(): %s\n", typ)
		c.Types[ctx] = typ
	}

	switch {
	case typ_void != nil:
		setBasicType(parser.CParserVoid)

	case typ_char != nil:
		setBasicType(parser.CParserChar)

	case typ_short != nil:
		setBasicType(parser.CParserShort)

	case typ_int != nil:
		setBasicType(parser.CParserInt)

	case typ_long != nil:
		setBasicType(parser.CParserLong)

	case typ_float != nil:
		setBasicType(parser.CParserFloat)

	case typ_double != nil:
		setBasicType(parser.CParserDouble)

	default:
		c.fail("type specifier not supported")
	}

	fmt.Printf("  result: %s\n", c.Types[ctx])
}

func (c *MainPass) EnterParameterTypeList(ctx *parser.ParameterTypeListContext) {
	fmt.Printf("EnterParameterTypeList(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitParameterTypeList(ctx *parser.ParameterTypeListContext) {
	fmt.Printf("ExitParameterTypeList(): %s\n", ctx.GetText())

	paramList := ctx.ParameterList()
	ellipsis := ctx.Ellipsis()

	switch {
	case paramList != nil && ellipsis == nil:
		params := c.Types[paramList].(*types.FunctionType)
		c.Types[ctx] = &types.FunctionType{
			Params: params.Params,
		}

	default:
		c.fail("varargs not supported")
	}

	fmt.Printf("  result: %+v\n", c.Types[ctx])
}

func (c *MainPass) EnterParameterList(ctx *parser.ParameterListContext) {
	fmt.Printf("EnterParameterList(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitParameterList(ctx *parser.ParameterListContext) {
	fmt.Printf("ExitParameterList(): %s\n", ctx.GetText())

	paramDecl := ctx.ParameterDeclaration()
	paramList := ctx.ParameterList()

	switch {
	case paramDecl != nil && paramList == nil:
		param := c.Declarations[paramDecl]
		c.Types[ctx] = &types.FunctionType{
			Params: []types.Param{
				types.Param{
					Name: param.Name,
					Type: param.Type,
				},
			},
		}

	case paramList != nil && paramDecl != nil:
		params := c.Types[paramList].(*types.FunctionType)
		param := c.Declarations[paramDecl]
		c.Types[ctx] = &types.FunctionType{
			Params: append(params.Params, types.Param{
				Name: param.Name,
				Type: param.Type,
			}),
		}

	default:
		c.fail("ExitParameterList: invalid case")
	}

	fmt.Printf("  result: %+v\n", c.Types[ctx])
}

func (c *MainPass) EnterParameterDeclaration(ctx *parser.ParameterDeclarationContext) {
	fmt.Printf("EnterParameterDeclaration(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitParameterDeclaration(ctx *parser.ParameterDeclarationContext) {
	fmt.Printf("ExitParameterDeclaration(): %s\n", ctx.GetText())

	declSpec := ctx.DeclarationSpecifiers()
	declarator := ctx.Declarator()

	switch {
	case declSpec != nil && declarator != nil:
		typ := c.Types[declSpec]
		decl := c.Declarations[declarator]
		fmt.Printf("  typ: %s\n", typ)
		fmt.Printf("  decl: %s\n", decl)

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

	default:
		c.fail("ExitParameterDeclaration: invalid case")
	}

	fmt.Printf("  result: %+v\n", c.Declarations[ctx])
}

func (c *MainPass) EnterCompoundStatement(ctx *parser.CompoundStatementContext) {
	fmt.Printf("EnterCompoundStatement(): %s\n", ctx.GetText())
	c.runEnterContinuation(ctx)
}

func (c *MainPass) ExitCompoundStatement(ctx *parser.CompoundStatementContext) {
	fmt.Printf("ExitCompoundStatement(): %s\n", ctx.GetText())
	c.runExitContinuation(ctx)
}

func (c *MainPass) EnterJumpStatement(ctx *parser.JumpStatementContext) {
	fmt.Printf("EnterJumpStatement(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitJumpStatement(ctx *parser.JumpStatementContext) {
	fmt.Printf("ExitJumpStatement(): %s\n", ctx.GetText())

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

func (c *MainPass) EnterExpression(ctx *parser.ExpressionContext) {
	fmt.Printf("EnterExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitExpression(ctx *parser.ExpressionContext) {
	fmt.Printf("ExitExpression(): %s\n", ctx.GetText())

	exp1 := ctx.AssignmentExpression()
	exp2 := ctx.Expression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterAssignmentExpression(ctx *parser.AssignmentExpressionContext) {
	fmt.Printf("EnterAssignmentExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitAssignmentExpression(ctx *parser.AssignmentExpressionContext) {
	fmt.Printf("ExitAssignmentExpression(): %s\n", ctx.GetText())

	exp1 := ctx.ConditionalExpression()

	switch {
	case exp1 != nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) ExitConditionalExpression(ctx *parser.ConditionalExpressionContext) {
	fmt.Printf("ExitConditionalExpression(): %s\n", ctx.GetText())
	fmt.Printf("  LogicalOrExpression: %v\n", ctx.LogicalOrExpression() != nil)

	exp1 := ctx.LogicalOrExpression()
	c.cgen.MoveValue(ctx, exp1)
}

func (c *MainPass) EnterLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) {
	fmt.Printf("EnterLogicalOrExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) {
	fmt.Printf("ExitLogicalOrExpression(): %s\n", ctx.GetText())

	exp1 := ctx.LogicalAndExpression()
	exp2 := ctx.LogicalOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) {
	fmt.Printf("EnterLogicalAndExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) {
	fmt.Printf("ExitLogicalAndExpression(): %s\n", ctx.GetText())

	exp1 := ctx.InclusiveOrExpression()
	exp2 := ctx.LogicalAndExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterInclusiveOrExpression(ctx *parser.InclusiveOrExpressionContext) {
	fmt.Printf("EnterInclusiveOrExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitInclusiveOrExpression(ctx *parser.InclusiveOrExpressionContext) {
	fmt.Printf("ExitInclusiveOrExpression(): %s\n", ctx.GetText())

	exp1 := ctx.ExclusiveOrExpression()
	exp2 := ctx.InclusiveOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterExclusiveOrExpression(ctx *parser.ExclusiveOrExpressionContext) {
	fmt.Printf("EnterExclusiveOrExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitExclusiveOrExpression(ctx *parser.ExclusiveOrExpressionContext) {
	fmt.Printf("ExitExclusiveOrExpression(): %s\n", ctx.GetText())

	exp1 := ctx.AndExpression()
	exp2 := ctx.ExclusiveOrExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterAndExpression(ctx *parser.AndExpressionContext) {
	fmt.Printf("EnterAndExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitAndExpression(ctx *parser.AndExpressionContext) {
	fmt.Printf("ExitAndExpression(): %s\n", ctx.GetText())

	exp1 := ctx.EqualityExpression()
	exp2 := ctx.AndExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterEqualityExpression(ctx *parser.EqualityExpressionContext) {
	fmt.Printf("EnterEqualityExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitEqualityExpression(ctx *parser.EqualityExpressionContext) {
	fmt.Printf("ExitEqualityExpression(): %s\n", ctx.GetText())

	exp1 := ctx.RelationalExpression()
	exp2 := ctx.EqualityExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterRelationalExpression(ctx *parser.RelationalExpressionContext) {
	fmt.Printf("EnterRelationalExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitRelationalExpression(ctx *parser.RelationalExpressionContext) {
	fmt.Printf("ExitRelationalExpression(): %s\n", ctx.GetText())

	exp1 := ctx.ShiftExpression()
	exp2 := ctx.RelationalExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterShiftExpression(ctx *parser.ShiftExpressionContext) {
	fmt.Printf("EnterShiftExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitShiftExpression(ctx *parser.ShiftExpressionContext) {
	fmt.Printf("ExitShiftExpression(): %s\n", ctx.GetText())

	exp1 := ctx.AdditiveExpression()
	exp2 := ctx.ShiftExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	fmt.Printf("EnterAdditiveExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	fmt.Printf("ExitAdditiveExpression(): %s\n", ctx.GetText())

	e1 := ctx.MultiplicativeExpression()
	e2 := ctx.AdditiveExpression()
	plus := ctx.Plus()
	minus := ctx.Minus()

	fmt.Printf("e1: %v\n", e1)
	fmt.Printf("e2: %v\n", e2)
	fmt.Printf("plus: %v\n", plus)
	fmt.Printf("minus: %v\n", minus)

	switch {
	case e1 != nil && e2 == nil:
		c.cgen.MoveValue(ctx, e1)

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
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) {
	fmt.Printf("EnterMultiplicativeExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) {
	fmt.Printf("ExitMultiplicativeExpression(): %s\n", ctx.GetText())

	exp1 := ctx.CastExpression()
	exp2 := ctx.MultiplicativeExpression()

	switch {
	case exp1 != nil && exp2 == nil:
		c.cgen.MoveValue(ctx, exp1)

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

	exp1 := ctx.TypeName()
	exp2 := ctx.CastExpression()
	exp3 := ctx.UnaryExpression()

	switch {
	case exp1 == nil && exp2 == nil && exp3 != nil:
		c.cgen.MoveValue(ctx, exp3)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterUnaryExpression(ctx *parser.UnaryExpressionContext) {
	fmt.Printf("EnterUnaryExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitUnaryExpression(ctx *parser.UnaryExpressionContext) {
	fmt.Printf("ExitUnaryExpression(): %s\n", ctx.GetText())

	exp1 := ctx.PostfixExpression()

	switch {
	case exp1 != nil:
		c.cgen.MoveValue(ctx, exp1)

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterPostfixExpression(ctx *parser.PostfixExpressionContext) {
	fmt.Printf("EnterPostfixExpression(): %s\n", ctx.GetText())
}

func (c *MainPass) ExitPostfixExpression(ctx *parser.PostfixExpressionContext) {
	fmt.Printf("ExitPostfixExpression(): %s\n", ctx.GetText())

	exp1 := ctx.PrimaryExpression()

	switch {
	case exp1 != nil:
		c.cgen.MoveValue(ctx, exp1)

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

	cp := ctx.GetParser().(*parser.CParser)
	tok := ctx.GetStart()

	switch {
	case tok.GetTokenType() == parser.CLexerConstant:
		typ := types.NewBasicType(parser.CParserInt, cp)
		c.cgen.CreateIntValue(ctx, typ, tok.GetText())

	case ctx.Identifier() != nil:
		fmt.Printf("ExitPrimaryExpression(): identifier(%s)\n", ctx.Identifier().GetText())

	default:
		panic("unhandled case!")
	}
}

func (c *MainPass) setEnterContinuation(ctx antlr.ParserRuleContext, f func()) {
	_, found := c.EnterContinuations[ctx]
	if found {
		c.fail("continuation already set!")
	}

	c.EnterContinuations[ctx] = f
}

func (c *MainPass) setExitContinuation(ctx antlr.ParserRuleContext, f func()) {
	_, found := c.ExitContinuations[ctx]
	if found {
		c.fail("continuation already set!")
	}

	c.ExitContinuations[ctx] = f
}

func (c *MainPass) runEnterContinuation(ctx antlr.ParserRuleContext) {
	f, found := c.EnterContinuations[ctx]
	if found {
		f()
	}
}

func (c *MainPass) runExitContinuation(ctx antlr.ParserRuleContext) {
	f, found := c.ExitContinuations[ctx]
	if found {
		f()
	}
}

func (c *MainPass) fail(format string, a ...interface{}) {
	c.Err = fmt.Errorf(format, a...)
	panic(c.Err.Error())
}
