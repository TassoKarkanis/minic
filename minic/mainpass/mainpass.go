package mainpass

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
	Output             io.Writer
	Symbols            *symbols.Table
	Function           *types.FunctionType
	cgen               *codegen.Codegen
	Err                error
	Declarations       map[antlr.ParserRuleContext]Declaration
	EnterContinuations map[antlr.ParserRuleContext]func()
	ExitContinuations  map[antlr.ParserRuleContext]func()
	types              map[antlr.ParserRuleContext]types.Type
	enterStack         []string
}

var _ parser.CListener = &MainPass{}

func NewMainPass(output io.Writer) *MainPass {
	return &MainPass{
		Output:             output,
		Symbols:            symbols.NewTable(),
		Declarations:       make(map[antlr.ParserRuleContext]Declaration),
		EnterContinuations: make(map[antlr.ParserRuleContext]func()),
		ExitContinuations:  make(map[antlr.ParserRuleContext]func()),
		types:              make(map[antlr.ParserRuleContext]types.Type),
	}
}

func (c *MainPass) EnterEveryRule(ctx antlr.ParserRuleContext) {
}

func (c *MainPass) ExitEveryRule(ctx antlr.ParserRuleContext) {
}

func (c *MainPass) VisitErrorNode(node antlr.ErrorNode) {
}

func (c *MainPass) VisitTerminal(node antlr.TerminalNode) {
}

func (c *MainPass) EnterCompilationUnit(ctx *parser.CompilationUnitContext) {
	c.enterRule(ctx, "CompilationUnit")
}

func (c *MainPass) ExitCompilationUnit(ctx *parser.CompilationUnitContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterTypeSpecifier(ctx *parser.TypeSpecifierContext) {
	c.enterRule(ctx, "TypeSpecifier")
}

func (c *MainPass) ExitTypeSpecifier(ctx *parser.TypeSpecifierContext) {
	e := c.exitRule(ctx)
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
		c.setType(ctx, typ)
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

	c.debugf("result: %s\n", c.getType(ctx))
}

func (c *MainPass) EnterCompoundStatement(ctx *parser.CompoundStatementContext) {
	c.enterRule(ctx, "CompoundStatement")
}

func (c *MainPass) ExitCompoundStatement(ctx *parser.CompoundStatementContext) {
	e := c.exitRule(ctx)
	defer e()
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

func (c *MainPass) setType(ctx antlr.ParserRuleContext, typ types.Type) {
	_, found := c.types[ctx]
	if found {
		c.fail("setType(): type already set")
	}

	c.types[ctx] = typ
}

func (c *MainPass) getType(ctx antlr.ParserRuleContext) types.Type {
	typ, found := c.types[ctx]
	if !found {
		c.fail("getType(): type not set")
	}

	return typ
}

func (c *MainPass) copyType(dest, src antlr.ParserRuleContext) {
	typ := c.getType(src)
	c.setType(dest, typ)
}

func (c *MainPass) fail(format string, a ...interface{}) {
	c.Err = fmt.Errorf(format, a...)
	panic(c.Err.Error())
}

func (c *MainPass) enterRule(ctx antlr.ParserRuleContext, name string) {
	// log the message
	funcName := "Enter" + name + "()"
	c.debugf("%s: %s", funcName, ctx.GetText())

	// push the enter stack
	c.enterStack = append(c.enterStack, name)

	// run the enter continuation
	f, found := c.EnterContinuations[ctx]
	if found {
		f()
	}
}

func (c *MainPass) exitRule(ctx antlr.ParserRuleContext) func() {
	// pop the stack
	name := c.enterStack[len(c.enterStack)-1]
	c.enterStack = c.enterStack[0 : len(c.enterStack)-1]

	// print the message
	funcName := "Exit" + name + "()"
	c.debugf("%s: %s", funcName, ctx.GetText())

	// put it back on the stack (to indent log statements in the current function)
	c.enterStack = append(c.enterStack, name)

	// return a function that pops the stack
	return func() {
		// run the exit continuation
		f, found := c.ExitContinuations[ctx]
		if found {
			f()
		}

		// pop the stack
		c.enterStack = c.enterStack[0 : len(c.enterStack)-1]
	}
}

func (c *MainPass) debugf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	indent := strings.Repeat(" ", len(c.enterStack))
	log.Debugf("%s%s", indent, msg)
}
