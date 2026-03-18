package com.ink.yht.sprout.idea.highlight

import com.intellij.lexer.Lexer
import com.intellij.openapi.editor.DefaultLanguageHighlighterColors
import com.intellij.openapi.editor.HighlighterColors
import com.intellij.openapi.editor.colors.TextAttributesKey
import com.intellij.openapi.fileTypes.SyntaxHighlighter
import com.intellij.openapi.fileTypes.SyntaxHighlighterFactory
import com.intellij.openapi.project.Project
import com.intellij.openapi.vfs.VirtualFile
import com.intellij.psi.tree.IElementType
import com.ink.yht.sprout.idea.lang.SproutLexerAdapter
import com.ink.yht.sprout.idea.lang.SproutTokenTypes

class SproutSyntaxHighlighterFactory : SyntaxHighlighterFactory() {
    override fun getSyntaxHighlighter(project: Project?, virtualFile: VirtualFile?): SyntaxHighlighter {
        return SproutSyntaxHighlighter()
    }
}

class SproutSyntaxHighlighter : SyntaxHighlighter {
    companion object {
        // Define text attribute keys for different token types
        val KEYWORD = TextAttributesKey.createTextAttributesKey("SPROUT.KEYWORD", DefaultLanguageHighlighterColors.KEYWORD)
        val HTTP_METHOD = TextAttributesKey.createTextAttributesKey("SPROUT.HTTP_METHOD", DefaultLanguageHighlighterColors.KEYWORD)
        val STRING = TextAttributesKey.createTextAttributesKey("SPROUT.STRING", DefaultLanguageHighlighterColors.STRING)
        val NUMBER = TextAttributesKey.createTextAttributesKey("SPROUT.NUMBER", DefaultLanguageHighlighterColors.NUMBER)
        val IDENTIFIER = TextAttributesKey.createTextAttributesKey("SPROUT.IDENTIFIER", DefaultLanguageHighlighterColors.IDENTIFIER)
        val OPERATOR = TextAttributesKey.createTextAttributesKey("SPROUT.OPERATOR", DefaultLanguageHighlighterColors.OPERATION_SIGN)
        val COMMENT = TextAttributesKey.createTextAttributesKey("SPROUT.COMMENT", DefaultLanguageHighlighterColors.LINE_COMMENT)
        val BRACES = TextAttributesKey.createTextAttributesKey("SPROUT.BRACES", DefaultLanguageHighlighterColors.BRACES)
        val PARENTHESES = TextAttributesKey.createTextAttributesKey("SPROUT.PARENTHESES", DefaultLanguageHighlighterColors.PARENTHESES)
        val BRACKETS = TextAttributesKey.createTextAttributesKey("SPROUT.BRACKETS", DefaultLanguageHighlighterColors.BRACKETS)
        val BAD_CHARACTER = TextAttributesKey.createTextAttributesKey("SPROUT.BAD_CHARACTER", HighlighterColors.BAD_CHARACTER)

        val KEYWORD_KEYS = arrayOf(KEYWORD)
        val HTTP_METHOD_KEYS = arrayOf(HTTP_METHOD)
        val STRING_KEYS = arrayOf(STRING)
        val NUMBER_KEYS = arrayOf(NUMBER)
        val IDENTIFIER_KEYS = arrayOf(IDENTIFIER)
        val OPERATOR_KEYS = arrayOf(OPERATOR)
        val COMMENT_KEYS = arrayOf(COMMENT)
        val BRACES_KEYS = arrayOf(BRACES)
        val PARENTHESES_KEYS = arrayOf(PARENTHESES)
        val BRACKETS_KEYS = arrayOf(BRACKETS)
        val BAD_CHAR_KEYS = arrayOf(BAD_CHARACTER)
        val EMPTY_KEYS = arrayOf<TextAttributesKey>()
    }

    override fun getHighlightingLexer(): Lexer = SproutLexerAdapter()

    override fun getTokenHighlights(tokenType: IElementType): Array<TextAttributesKey> {
        return when (tokenType) {
            // Keywords
            SproutTokenTypes.TYPE, SproutTokenTypes.SERVER, SproutTokenTypes.PREFIX,
            SproutTokenTypes.PUBLIC, SproutTokenTypes.PRIVATE, SproutTokenTypes.SERVICE,
            SproutTokenTypes.RPC -> KEYWORD_KEYS
            
            // HTTP Methods
            SproutTokenTypes.GET, SproutTokenTypes.POST, SproutTokenTypes.PUT,
            SproutTokenTypes.DELETE, SproutTokenTypes.PATCH -> HTTP_METHOD_KEYS
            
            // Literals
            SproutTokenTypes.STRING, SproutTokenTypes.RAW_STRING -> STRING_KEYS
            SproutTokenTypes.NUMBER -> NUMBER_KEYS
            SproutTokenTypes.IDENTIFIER -> IDENTIFIER_KEYS
            
            // Operators
            SproutTokenTypes.ARROW, SproutTokenTypes.COLON, SproutTokenTypes.SLASH -> OPERATOR_KEYS
            
            // Comments
            SproutTokenTypes.COMMENT -> COMMENT_KEYS
            
            // Braces
            SproutTokenTypes.LBRACE, SproutTokenTypes.RBRACE -> BRACES_KEYS
            SproutTokenTypes.LPAREN, SproutTokenTypes.RPAREN -> PARENTHESES_KEYS
            SproutTokenTypes.LBRACKET, SproutTokenTypes.RBRACKET -> BRACKETS_KEYS
            
            // Bad character
            com.intellij.psi.TokenType.BAD_CHARACTER -> BAD_CHAR_KEYS
            
            else -> EMPTY_KEYS
        }
    }
}
