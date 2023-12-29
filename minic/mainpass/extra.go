package mainpass

import (
	"github.com/TassoKarkanis/minic/parser"
)

//
// Statement
//

//
// Other
//

func (c *MainPass) EnterBlockItemList(ctx *parser.BlockItemListContext) {
	c.enterf("BlockItemList", "%s", ctx.GetText())
}

func (c *MainPass) ExitBlockItemList(ctx *parser.BlockItemListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterBlockItem(ctx *parser.BlockItemContext) {
	c.enterf("BlockItem", "%s", ctx.GetText())
}

func (c *MainPass) ExitBlockItem(ctx *parser.BlockItemContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterAtomicTypeSpecifier(ctx *parser.AtomicTypeSpecifierContext) {
	c.enterf("AtomicTypeSpecifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitAtomicTypeSpecifier(ctx *parser.AtomicTypeSpecifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterArgumentExpressionList(ctx *parser.ArgumentExpressionListContext) {
	c.enterf("ArgumentExpressionList", "%s", ctx.GetText())
}

func (c *MainPass) ExitArgumentExpressionList(ctx *parser.ArgumentExpressionListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterAlignmentSpecifier(ctx *parser.AlignmentSpecifierContext) {
	c.enterf("AlignmentSpecifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitAlignmentSpecifier(ctx *parser.AlignmentSpecifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterDeclarationList(ctx *parser.DeclarationListContext) {
	c.enterf("DeclarationList", "%s", ctx.GetText())
}

func (c *MainPass) ExitDeclarationList(ctx *parser.DeclarationListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterDeclarationSpecifiers2(ctx *parser.DeclarationSpecifiers2Context) {
	c.enterf("DeclarationSpecifiers2", "%s", ctx.GetText())
}

func (c *MainPass) ExitDeclarationSpecifiers2(ctx *parser.DeclarationSpecifiers2Context) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterDesignation(ctx *parser.DesignationContext) {
	c.enterf("Designation", "%s", ctx.GetText())
}

func (c *MainPass) ExitDesignation(ctx *parser.DesignationContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterDesignator(ctx *parser.DesignatorContext) {
	c.enterf("Designator", "%s", ctx.GetText())
}

func (c *MainPass) ExitDesignator(ctx *parser.DesignatorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterDesignatorList(ctx *parser.DesignatorListContext) {
	c.enterf("DesignatorList", "%s", ctx.GetText())
}

func (c *MainPass) ExitDesignatorList(ctx *parser.DesignatorListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterDirectAbstractDeclarator(ctx *parser.DirectAbstractDeclaratorContext) {
	c.enterf("DirectAbstractDeclarator", "%s", ctx.GetText())
}

func (c *MainPass) ExitDirectAbstractDeclarator(ctx *parser.DirectAbstractDeclaratorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterEnumSpecifier(ctx *parser.EnumSpecifierContext) {
	c.enterf("EnumSpecifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitEnumSpecifier(ctx *parser.EnumSpecifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterEnumerationConstant(ctx *parser.EnumerationConstantContext) {
	c.enterf("EnumerationConstant", "%s", ctx.GetText())
}

func (c *MainPass) ExitEnumerationConstant(ctx *parser.EnumerationConstantContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterEnumerator(ctx *parser.EnumeratorContext) {
	c.enterf("Enumerator", "%s", ctx.GetText())
}

func (c *MainPass) ExitEnumerator(ctx *parser.EnumeratorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterEnumeratorList(ctx *parser.EnumeratorListContext) {
	c.enterf("EnumeratorList", "%s", ctx.GetText())
}

func (c *MainPass) ExitEnumeratorList(ctx *parser.EnumeratorListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterForCondition(ctx *parser.ForConditionContext) {
	c.enterf("ForCondition", "%s", ctx.GetText())
}

func (c *MainPass) ExitForCondition(ctx *parser.ForConditionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterForDeclaration(ctx *parser.ForDeclarationContext) {
	c.enterf("ForDeclaration", "%s", ctx.GetText())
}

func (c *MainPass) ExitForDeclaration(ctx *parser.ForDeclarationContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterFunctionSpecifier(ctx *parser.FunctionSpecifierContext) {
	c.enterf("FunctionSpecifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitFunctionSpecifier(ctx *parser.FunctionSpecifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterGccAttribute(ctx *parser.GccAttributeContext) {
	c.enterf("GccAttribute", "%s", ctx.GetText())
}

func (c *MainPass) ExitGccAttribute(ctx *parser.GccAttributeContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterGccAttributeList(ctx *parser.GccAttributeListContext) {
	c.enterf("GccAttributeList", "%s", ctx.GetText())
}

func (c *MainPass) ExitGccAttributeList(ctx *parser.GccAttributeListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterGccAttributeSpecifier(ctx *parser.GccAttributeSpecifierContext) {
	c.enterf("GccAttributeSpecifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitGccAttributeSpecifier(ctx *parser.GccAttributeSpecifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterGccDeclaratorExtension(ctx *parser.GccDeclaratorExtensionContext) {
	c.enterf("GccDeclaratorExtension", "%s", ctx.GetText())
}

func (c *MainPass) ExitGccDeclaratorExtension(ctx *parser.GccDeclaratorExtensionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterGenericAssocList(ctx *parser.GenericAssocListContext) {
	c.enterf("GenericAssocList", "%s", ctx.GetText())
}

func (c *MainPass) ExitGenericAssocList(ctx *parser.GenericAssocListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterGenericAssociation(ctx *parser.GenericAssociationContext) {
	c.enterf("GenericAssociation", "%s", ctx.GetText())
}

func (c *MainPass) ExitGenericAssociation(ctx *parser.GenericAssociationContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterGenericSelection(ctx *parser.GenericSelectionContext) {
	c.enterf("GenericSelection", "%s", ctx.GetText())
}

func (c *MainPass) ExitGenericSelection(ctx *parser.GenericSelectionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterIdentifierList(ctx *parser.IdentifierListContext) {
	c.enterf("IdentifierList", "%s", ctx.GetText())
}

func (c *MainPass) ExitIdentifierList(ctx *parser.IdentifierListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterInitDeclarator(ctx *parser.InitDeclaratorContext) {
	c.enterf("InitDeclarator", "%s", ctx.GetText())
}

func (c *MainPass) ExitInitDeclarator(ctx *parser.InitDeclaratorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterInitDeclaratorList(ctx *parser.InitDeclaratorListContext) {
	c.enterf("InitDeclaratorList", "%s", ctx.GetText())
}

func (c *MainPass) ExitInitDeclaratorList(ctx *parser.InitDeclaratorListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterInitializer(ctx *parser.InitializerContext) {
	c.enterf("Initializer", "%s", ctx.GetText())
}

func (c *MainPass) ExitInitializer(ctx *parser.InitializerContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterInitializerList(ctx *parser.InitializerListContext) {
	c.enterf("InitializerList", "%s", ctx.GetText())
}

func (c *MainPass) ExitInitializerList(ctx *parser.InitializerListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterNestedParenthesesBlock(ctx *parser.NestedParenthesesBlockContext) {
	c.enterf("NestedParenthesesBlock", "%s", ctx.GetText())
}

func (c *MainPass) ExitNestedParenthesesBlock(ctx *parser.NestedParenthesesBlockContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterPointer(ctx *parser.PointerContext) {
	c.enterf("Pointer", "%s", ctx.GetText())
}

func (c *MainPass) ExitPointer(ctx *parser.PointerContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterSpecifierQualifierList(ctx *parser.SpecifierQualifierListContext) {
	c.enterf("SpecifierQualifierList", "%s", ctx.GetText())
}

func (c *MainPass) ExitSpecifierQualifierList(ctx *parser.SpecifierQualifierListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterStaticAssertDeclaration(ctx *parser.StaticAssertDeclarationContext) {
	c.enterf("StaticAssertDeclaration", "%s", ctx.GetText())
}

func (c *MainPass) ExitStaticAssertDeclaration(ctx *parser.StaticAssertDeclarationContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterStorageClassSpecifier(ctx *parser.StorageClassSpecifierContext) {
	c.enterf("StorageClassSpecifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitStorageClassSpecifier(ctx *parser.StorageClassSpecifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterStructDeclaration(ctx *parser.StructDeclarationContext) {
	c.enterf("StructDeclaration", "%s", ctx.GetText())
}

func (c *MainPass) ExitStructDeclaration(ctx *parser.StructDeclarationContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterStructDeclarationList(ctx *parser.StructDeclarationListContext) {
	c.enterf("StructDeclarationList", "%s", ctx.GetText())
}

func (c *MainPass) ExitStructDeclarationList(ctx *parser.StructDeclarationListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterStructDeclarator(ctx *parser.StructDeclaratorContext) {
	c.enterf("StructDeclarator", "%s", ctx.GetText())
}

func (c *MainPass) ExitStructDeclarator(ctx *parser.StructDeclaratorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterStructDeclaratorList(ctx *parser.StructDeclaratorListContext) {
	c.enterf("StructDeclaratorList", "%s", ctx.GetText())
}

func (c *MainPass) ExitStructDeclaratorList(ctx *parser.StructDeclaratorListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterStructOrUnion(ctx *parser.StructOrUnionContext) {
	c.enterf("StructOrUnion", "%s", ctx.GetText())
}

func (c *MainPass) ExitStructOrUnion(ctx *parser.StructOrUnionContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterStructOrUnionSpecifier(ctx *parser.StructOrUnionSpecifierContext) {
	c.enterf("StructOrUnionSpecifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitStructOrUnionSpecifier(ctx *parser.StructOrUnionSpecifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterTranslationUnit(ctx *parser.TranslationUnitContext) {
	c.enterf("TranslationUnit", "%s", ctx.GetText())
}

func (c *MainPass) ExitTranslationUnit(ctx *parser.TranslationUnitContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterTypeName(ctx *parser.TypeNameContext) {
	c.enterf("TypeName", "%s", ctx.GetText())
}

func (c *MainPass) ExitTypeName(ctx *parser.TypeNameContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterTypeQualifier(ctx *parser.TypeQualifierContext) {
	c.enterf("TypeQualifier", "%s", ctx.GetText())
}

func (c *MainPass) ExitTypeQualifier(ctx *parser.TypeQualifierContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterTypeQualifierList(ctx *parser.TypeQualifierListContext) {
	c.enterf("TypeQualifierList", "%s", ctx.GetText())
}

func (c *MainPass) ExitTypeQualifierList(ctx *parser.TypeQualifierListContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
func (c *MainPass) EnterTypedefName(ctx *parser.TypedefNameContext) {
	c.enterf("ypedefName", "%s", ctx.GetText())
}

func (c *MainPass) ExitTypedefName(ctx *parser.TypedefNameContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}

func (c *MainPass) EnterUnaryOperator(ctx *parser.UnaryOperatorContext) {
	c.enterf("UnaryOperator", "%s", ctx.GetText())
}

func (c *MainPass) ExitUnaryOperator(ctx *parser.UnaryOperatorContext) {
	e := c.exitf("%s", ctx.GetText())
	defer e()
}
