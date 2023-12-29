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
	c.enterRule(ctx, "BlockItemList")
}

func (c *MainPass) ExitBlockItemList(ctx *parser.BlockItemListContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterBlockItem(ctx *parser.BlockItemContext) {
	c.enterRule(ctx, "BlockItem")
}

func (c *MainPass) ExitBlockItem(ctx *parser.BlockItemContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterAtomicTypeSpecifier(ctx *parser.AtomicTypeSpecifierContext) {
	c.enterRule(ctx, "AtomicTypeSpecifier")
}

func (c *MainPass) ExitAtomicTypeSpecifier(ctx *parser.AtomicTypeSpecifierContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterArgumentExpressionList(ctx *parser.ArgumentExpressionListContext) {
	c.enterRule(ctx, "ArgumentExpressionList")
}

func (c *MainPass) ExitArgumentExpressionList(ctx *parser.ArgumentExpressionListContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterAlignmentSpecifier(ctx *parser.AlignmentSpecifierContext) {
	c.enterRule(ctx, "AlignmentSpecifier")
}

func (c *MainPass) ExitAlignmentSpecifier(ctx *parser.AlignmentSpecifierContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterDesignation(ctx *parser.DesignationContext) {
	c.enterRule(ctx, "Designation")
}

func (c *MainPass) ExitDesignation(ctx *parser.DesignationContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterDesignator(ctx *parser.DesignatorContext) {
	c.enterRule(ctx, "Designator")
}

func (c *MainPass) ExitDesignator(ctx *parser.DesignatorContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterDesignatorList(ctx *parser.DesignatorListContext) {
	c.enterRule(ctx, "DesignatorList")
}

func (c *MainPass) ExitDesignatorList(ctx *parser.DesignatorListContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterDirectAbstractDeclarator(ctx *parser.DirectAbstractDeclaratorContext) {
	c.enterRule(ctx, "DirectAbstractDeclarator")
}

func (c *MainPass) ExitDirectAbstractDeclarator(ctx *parser.DirectAbstractDeclaratorContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterEnumSpecifier(ctx *parser.EnumSpecifierContext) {
	c.enterRule(ctx, "EnumSpecifier")
}

func (c *MainPass) ExitEnumSpecifier(ctx *parser.EnumSpecifierContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterEnumerationConstant(ctx *parser.EnumerationConstantContext) {
	c.enterRule(ctx, "EnumerationConstant")
}

func (c *MainPass) ExitEnumerationConstant(ctx *parser.EnumerationConstantContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterEnumerator(ctx *parser.EnumeratorContext) {
	c.enterRule(ctx, "Enumerator")
}

func (c *MainPass) ExitEnumerator(ctx *parser.EnumeratorContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterEnumeratorList(ctx *parser.EnumeratorListContext) {
	c.enterRule(ctx, "EnumeratorList")
}

func (c *MainPass) ExitEnumeratorList(ctx *parser.EnumeratorListContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterForCondition(ctx *parser.ForConditionContext) {
	c.enterRule(ctx, "ForCondition")
}

func (c *MainPass) ExitForCondition(ctx *parser.ForConditionContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterForDeclaration(ctx *parser.ForDeclarationContext) {
	c.enterRule(ctx, "ForDeclaration")
}

func (c *MainPass) ExitForDeclaration(ctx *parser.ForDeclarationContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterFunctionSpecifier(ctx *parser.FunctionSpecifierContext) {
	c.enterRule(ctx, "FunctionSpecifier")
}

func (c *MainPass) ExitFunctionSpecifier(ctx *parser.FunctionSpecifierContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterGccAttribute(ctx *parser.GccAttributeContext) {
	c.enterRule(ctx, "GccAttribute")
}

func (c *MainPass) ExitGccAttribute(ctx *parser.GccAttributeContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterGccAttributeList(ctx *parser.GccAttributeListContext) {
	c.enterRule(ctx, "GccAttributeList")
}

func (c *MainPass) ExitGccAttributeList(ctx *parser.GccAttributeListContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterGccAttributeSpecifier(ctx *parser.GccAttributeSpecifierContext) {
	c.enterRule(ctx, "GccAttributeSpecifier")
}

func (c *MainPass) ExitGccAttributeSpecifier(ctx *parser.GccAttributeSpecifierContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterGccDeclaratorExtension(ctx *parser.GccDeclaratorExtensionContext) {
	c.enterRule(ctx, "GccDeclaratorExtension")
}

func (c *MainPass) ExitGccDeclaratorExtension(ctx *parser.GccDeclaratorExtensionContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterGenericAssocList(ctx *parser.GenericAssocListContext) {
	c.enterRule(ctx, "GenericAssocList")
}

func (c *MainPass) ExitGenericAssocList(ctx *parser.GenericAssocListContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterGenericAssociation(ctx *parser.GenericAssociationContext) {
	c.enterRule(ctx, "GenericAssociation")
}

func (c *MainPass) ExitGenericAssociation(ctx *parser.GenericAssociationContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterGenericSelection(ctx *parser.GenericSelectionContext) {
	c.enterRule(ctx, "GenericSelection")
}

func (c *MainPass) ExitGenericSelection(ctx *parser.GenericSelectionContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterIdentifierList(ctx *parser.IdentifierListContext) {
	c.enterRule(ctx, "IdentifierList")
}

func (c *MainPass) ExitIdentifierList(ctx *parser.IdentifierListContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterNestedParenthesesBlock(ctx *parser.NestedParenthesesBlockContext) {
	c.enterRule(ctx, "NestedParenthesesBlock")
}

func (c *MainPass) ExitNestedParenthesesBlock(ctx *parser.NestedParenthesesBlockContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterPointer(ctx *parser.PointerContext) {
	c.enterRule(ctx, "Pointer")
}

func (c *MainPass) ExitPointer(ctx *parser.PointerContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterSpecifierQualifierList(ctx *parser.SpecifierQualifierListContext) {
	c.enterRule(ctx, "SpecifierQualifierList")
}

func (c *MainPass) ExitSpecifierQualifierList(ctx *parser.SpecifierQualifierListContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterStaticAssertDeclaration(ctx *parser.StaticAssertDeclarationContext) {
	c.enterRule(ctx, "StaticAssertDeclaration")
}

func (c *MainPass) ExitStaticAssertDeclaration(ctx *parser.StaticAssertDeclarationContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterStorageClassSpecifier(ctx *parser.StorageClassSpecifierContext) {
	c.enterRule(ctx, "StorageClassSpecifier")
}

func (c *MainPass) ExitStorageClassSpecifier(ctx *parser.StorageClassSpecifierContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterStructDeclaration(ctx *parser.StructDeclarationContext) {
	c.enterRule(ctx, "StructDeclaration")
}

func (c *MainPass) ExitStructDeclaration(ctx *parser.StructDeclarationContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterStructDeclarationList(ctx *parser.StructDeclarationListContext) {
	c.enterRule(ctx, "StructDeclarationList")
}

func (c *MainPass) ExitStructDeclarationList(ctx *parser.StructDeclarationListContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterStructDeclarator(ctx *parser.StructDeclaratorContext) {
	c.enterRule(ctx, "StructDeclarator")
}

func (c *MainPass) ExitStructDeclarator(ctx *parser.StructDeclaratorContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterStructDeclaratorList(ctx *parser.StructDeclaratorListContext) {
	c.enterRule(ctx, "StructDeclaratorList")
}

func (c *MainPass) ExitStructDeclaratorList(ctx *parser.StructDeclaratorListContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterStructOrUnion(ctx *parser.StructOrUnionContext) {
	c.enterRule(ctx, "StructOrUnion")
}

func (c *MainPass) ExitStructOrUnion(ctx *parser.StructOrUnionContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterStructOrUnionSpecifier(ctx *parser.StructOrUnionSpecifierContext) {
	c.enterRule(ctx, "StructOrUnionSpecifier")
}

func (c *MainPass) ExitStructOrUnionSpecifier(ctx *parser.StructOrUnionSpecifierContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterTranslationUnit(ctx *parser.TranslationUnitContext) {
	c.enterRule(ctx, "TranslationUnit")
}

func (c *MainPass) ExitTranslationUnit(ctx *parser.TranslationUnitContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterTypeName(ctx *parser.TypeNameContext) {
	c.enterRule(ctx, "TypeName")
}

func (c *MainPass) ExitTypeName(ctx *parser.TypeNameContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterTypeQualifier(ctx *parser.TypeQualifierContext) {
	c.enterRule(ctx, "TypeQualifier")
}

func (c *MainPass) ExitTypeQualifier(ctx *parser.TypeQualifierContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterTypeQualifierList(ctx *parser.TypeQualifierListContext) {
	c.enterRule(ctx, "TypeQualifierList")
}

func (c *MainPass) ExitTypeQualifierList(ctx *parser.TypeQualifierListContext) {
	e := c.exitRule(ctx)
	defer e()
}
func (c *MainPass) EnterTypedefName(ctx *parser.TypedefNameContext) {
	c.enterRule(ctx, "ypedefName")
}

func (c *MainPass) ExitTypedefName(ctx *parser.TypedefNameContext) {
	e := c.exitRule(ctx)
	defer e()
}

func (c *MainPass) EnterUnaryOperator(ctx *parser.UnaryOperatorContext) {
	c.enterRule(ctx, "UnaryOperator")
}

func (c *MainPass) ExitUnaryOperator(ctx *parser.UnaryOperatorContext) {
	e := c.exitRule(ctx)
	defer e()
}
