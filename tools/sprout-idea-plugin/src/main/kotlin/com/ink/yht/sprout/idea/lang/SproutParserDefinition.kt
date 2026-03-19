package com.ink.yht.sprout.idea.lang

import com.intellij.lang.ASTNode
import com.intellij.lang.ParserDefinition
import com.intellij.lang.PsiParser
import com.intellij.lexer.Lexer
import com.intellij.openapi.project.Project
import com.intellij.psi.FileViewProvider
import com.intellij.psi.PsiFile
import com.intellij.psi.tree.IFileElementType
import com.intellij.psi.tree.TokenSet

class SproutParserDefinition : ParserDefinition {
    companion object {
        val FILE = IFileElementType(SproutLanguage)
        val WHITESPACES = TokenSet.create(SproutTokenTypes.WHITESPACE)
        val COMMENTS = TokenSet.create(SproutTokenTypes.COMMENT)
        val STRINGS = TokenSet.create(SproutTokenTypes.STRING, SproutTokenTypes.RAW_STRING)
    }

    override fun createLexer(project: Project?): Lexer = SproutLexerAdapter()
    override fun createParser(project: Project?): PsiParser = SproutParser()
    override fun getFileNodeType() = FILE
    override fun getWhitespaceTokens() = WHITESPACES
    override fun getCommentTokens() = COMMENTS
    override fun getStringLiteralElements() = STRINGS
    override fun createElement(node: ASTNode) = SproutTokenTypes.Factory.createElement(node)
    override fun createFile(viewProvider: FileViewProvider): PsiFile = SproutFile(viewProvider)
}
