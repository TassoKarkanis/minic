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
	Types              map[antlr.ParserRuleContext]types.Type
	Declarations       map[antlr.ParserRuleContext]Declaration
	EnterContinuations map[antlr.ParserRuleContext]func()
	ExitContinuations  map[antlr.ParserRuleContext]func()
	enterStack         []string
}

var _ parser.CListener = &MainPass{}

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

func (c *MainPass) EnterEveryRule(ctx antlr.ParserRuleContext) {
}

func (c *MainPass) ExitEveryRule(ctx antlr.ParserRuleContext) {
}

func (c *MainPass) VisitErrorNode(node antlr.ErrorNode) {
}

func (c *MainPass) VisitTerminal(node antlr.TerminalNode) {
}

func (c *MainPass) EnterCompilationUnit(ctx *parser.CompilationUnitContext) {
	c.enterf("CompilationUnit", ctx.GetText())
}

func (c *MainPass) ExitCompilationUnit(ctx *parser.CompilationUnitContext) {
	e := c.exitf("")
	defer e()
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

func (c *MainPass) EnterCompoundStatement(ctx *parser.CompoundStatementContext) {
	c.enterf("CompoundStatement", "%s", ctx.GetText())
	c.runEnterContinuation(ctx)
}

func (c *MainPass) ExitCompoundStatement(ctx *parser.CompoundStatementContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()

	c.runExitContinuation(ctx)
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
