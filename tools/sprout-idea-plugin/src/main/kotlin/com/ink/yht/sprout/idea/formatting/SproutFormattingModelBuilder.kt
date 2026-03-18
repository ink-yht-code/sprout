package com.ink.yht.sprout.idea.formatting

import com.intellij.formatting.FormattingModel
import com.intellij.formatting.FormattingModelBuilder
import com.intellij.formatting.FormattingModelProvider
import com.intellij.formatting.SpacingBuilder
import com.intellij.lang.ASTNode
import com.intellij.openapi.util.TextRange
import com.intellij.psi.PsiElement
import com.intellij.psi.codeStyle.CodeStyleSettings
import com.ink.yht.sprout.idea.lang.SproutLanguage

class SproutFormattingModelBuilder : FormattingModelBuilder {
    override fun createModel(element: PsiElement, settings: CodeStyleSettings): FormattingModel {
        val spacingBuilder = createSpacingBuilder(settings)
        val block = SproutBlock(
            element.node,
            settings,
            spacingBuilder,
            null,
            null
        )

        return FormattingModelProvider.createFormattingModelForPsiFile(
            element.containingFile,
            block,
            settings
        )
    }

    override fun getRangeAffectingIndent(file: com.intellij.psi.PsiFile, offset: Int, elementAtOffset: ASTNode): TextRange? {
        return null
    }

    private fun createSpacingBuilder(settings: CodeStyleSettings): SpacingBuilder {
        return SpacingBuilder(settings, SproutLanguage)
            .around(com.ink.yht.sprout.idea.lang.SproutTokenTypes.ARROW).spaces(1)
            .after(com.ink.yht.sprout.idea.lang.SproutTokenTypes.COMMA).spaces(1)
            .after(com.ink.yht.sprout.idea.lang.SproutTokenTypes.TYPE).spaces(1)
            .after(com.ink.yht.sprout.idea.lang.SproutTokenTypes.SERVER).spaces(1)
            .after(com.ink.yht.sprout.idea.lang.SproutTokenTypes.PREFIX).spaces(1)
            .after(com.ink.yht.sprout.idea.lang.SproutTokenTypes.PUBLIC).spaces(1)
            .after(com.ink.yht.sprout.idea.lang.SproutTokenTypes.PRIVATE).spaces(1)
            .after(com.ink.yht.sprout.idea.lang.SproutTokenTypes.SERVICE).spaces(1)
            .after(com.ink.yht.sprout.idea.lang.SproutTokenTypes.RPC).spaces(1)
            .after(com.ink.yht.sprout.idea.lang.SproutTokenTypes.LBRACE).lineBreakInCode()
            .before(com.ink.yht.sprout.idea.lang.SproutTokenTypes.RBRACE).lineBreakInCode()
    }
}
