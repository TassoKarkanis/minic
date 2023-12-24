package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/TassoKarkanis/minic/codegen"
	"github.com/TassoKarkanis/minic/parser"
	"github.com/TassoKarkanis/minic/symbols"
	"github.com/TassoKarkanis/minic/types"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	log "github.com/sirupsen/logrus"
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
	enterStack         []string
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
	c.enterf("CompilationUnit", ctx.GetText())
}

func (c *MainPass) ExitCompilationUnit(ctx *parser.CompilationUnitContext) {
	e := c.exitf("")
	defer e()
}

func (c *MainPass) EnterExternalDeclaration(ctx *parser.ExternalDeclarationContext) {
	c.enterf("ExternalDeclaration", ctx.GetText())
}

func (c *MainPass) ExitExternalDeclaration(ctx *parser.ExternalDeclarationContext) {
	e := c.exitf("")
	defer e()
}

func (c *MainPass) EnterFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	c.enterf("FunctionDefinition", ctx.GetText())

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
		c.debugf("  starting function: %s\n", c.Function)

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
	e := c.exitf("%s", ctx.GetText())
	defer e()

	c.Symbols.PopScope()
	c.cgen.EndStackFrame()
	c.cgen.Close()

	c.Function = nil
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

func (c *MainPass) EnterTypeSpecifier(ctx *parser.TypeSpecifierContext) {
	c.enterf("TypeSpecifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitTypeSpecifier(ctx *parser.TypeSpecifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

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
		c.debugf("setBasicType(): %s\n", typ)
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

	c.debugf("result: %s\n", c.Types[ctx])
}

func (c *MainPass) EnterParameterTypeList(ctx *parser.ParameterTypeListContext) {
	c.enterf("ParameterTypeList", "%s", ctx.GetText())
}

func (c *MainPass) ExitParameterTypeList(ctx *parser.ParameterTypeListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

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

	c.debugf("result: %+v\n", c.Types[ctx])
}

func (c *MainPass) EnterParameterList(ctx *parser.ParameterListContext) {
	c.enterf("ParameterList", "%s", ctx.GetText())
}

func (c *MainPass) ExitParameterList(ctx *parser.ParameterListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	paramDecl := ctx.ParameterDeclaration()
	paramList := ctx.ParameterList()

	switch {
	case paramDecl != nil && paramList == nil:
		param := c.Declarations[paramDecl]
		c.Types[ctx] = &types.FunctionType{
			Params: []types.Param{
				{
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

	c.debugf("result: %+v\n", c.Types[ctx])
}

func (c *MainPass) EnterParameterDeclaration(ctx *parser.ParameterDeclarationContext) {
	c.enterf("ParameterDeclaration", "%s", ctx.GetText())
}

func (c *MainPass) ExitParameterDeclaration(ctx *parser.ParameterDeclarationContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	declSpec := ctx.DeclarationSpecifiers()
	declarator := ctx.Declarator()
	cp := ctx.GetParser().(*parser.CParser)

	switch {
	case declSpec != nil && declarator != nil:
		typ := c.Types[declSpec]
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

func (c *MainPass) EnterCompoundStatement(ctx *parser.CompoundStatementContext) {
	c.enterf("CompoundStatement", "%s", ctx.GetText())
	c.runEnterContinuation(ctx)
}

func (c *MainPass) ExitCompoundStatement(ctx *parser.CompoundStatementContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	c.runExitContinuation(ctx)
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
		panic("unhandled case!")
	}
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
		panic("unhandled case!")
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
		panic("unhandled case!")
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
		panic("unhandled case!")
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
		panic("unhandled case!")
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
		panic("unhandled case!")
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
		panic("unhandled case!")
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
		panic("unhandled case!")
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
		panic("unhandled case!")
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
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	c.enterf("AdditiveExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitAdditiveExpression(ctx *parser.AdditiveExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	e1 := ctx.MultiplicativeExpression()
	e2 := ctx.AdditiveExpression()
	plus := ctx.Plus()
	minus := ctx.Minus()

	c.debugf("e1: %v\n", e1)
	c.debugf("e2: %v\n", e2)
	c.debugf("plus: %v\n", plus)
	c.debugf("minus: %v\n", minus)

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
	c.enterf("MultiplicativeExpression", "%s", ctx.GetText())
}

func (c *MainPass) ExitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

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
		panic("unhandled case!")
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
		panic("unhandled case!")
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
		panic("unhandled case!")
	}
}

func (c *MainPass) EnterPrimaryExpression(ctx *parser.PrimaryExpressionContext) {
	c.enterf("PrimaryExpression", "%s", ctx.GetText())
	tok := ctx.GetStart()
	c.debugf("start token index: %d\n", tok.GetTokenIndex())
	c.debugf("start token type: %d\n", tok.GetTokenType())
	c.debugf("start token channel: %d\n", tok.GetChannel())
}

func (c *MainPass) ExitPrimaryExpression(ctx *parser.PrimaryExpressionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	cp := ctx.GetParser().(*parser.CParser)
	tok := ctx.GetStart()

	switch {
	case tok.GetTokenType() == parser.CLexerConstant:
		typ := types.NewBasicType(parser.CParserInt, cp) // TODO: type
		c.cgen.CreateIntLiteralValue(ctx, typ, tok.GetText())

	case ctx.Identifier() != nil:
		// An identifier being evaluated.  Look up the symbol.
		name := ctx.Identifier().GetText()
		_, value, ok := c.Symbols.FindSymbol(name)
		if !ok {
			c.fail("unknown identifier: %s", name)
		}

		// it should have a value
		if value == nil {
			c.fail("symbol has no value: %s", name)
		}

		// forward the value
		c.cgen.CreateValue(ctx, value)

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

func (c *MainPass) enterf(stackStr string, format string, a ...interface{}) {
	// log the message
	msg := "Enter" + stackStr + "()"
	if format != "" {
		msg += ": " + fmt.Sprintf(format, a...)
	}
	c.debugf(msg)

	// push the enter stack
	c.enterStack = append(c.enterStack, stackStr)
}

func (c *MainPass) exitf(format string, a ...interface{}) func() {
	// pop the stack
	stackStr := c.enterStack[len(c.enterStack)-1]
	c.enterStack = c.enterStack[0 : len(c.enterStack)-1]

	// print the message
	msg := "Exit" + stackStr + "()"
	if format != "" {
		msg += ": " + fmt.Sprintf(format, a...)
	}
	c.debugf(msg)

	// put it back on the stack (to indent log statements in the current function)
	c.enterStack = append(c.enterStack, stackStr)

	// return a function that pops the stack
	return func() {
		c.enterStack = c.enterStack[0 : len(c.enterStack)-1]
	}
}

func (c *MainPass) debugf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	indent := strings.Repeat(" ", len(c.enterStack))
	log.Debugf("%s%s", indent, msg)
}
