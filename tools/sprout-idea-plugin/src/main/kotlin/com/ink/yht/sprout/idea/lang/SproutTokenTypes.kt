package com.ink.yht.sprout.idea.lang

import com.intellij.extapi.psi.ASTWrapperPsiElement
import com.intellij.lang.ASTNode
import com.intellij.psi.PsiElement
import com.intellij.psi.tree.IElementType
import com.intellij.psi.tree.TokenSet

interface SproutTokenTypes {
    companion object {
        // Keywords
        val TYPE: IElementType = IElementType("TYPE", SproutLanguage)
        val SERVER: IElementType = IElementType("SERVER", SproutLanguage)
        val PREFIX: IElementType = IElementType("PREFIX", SproutLanguage)
        val PUBLIC: IElementType = IElementType("PUBLIC", SproutLanguage)
        val PRIVATE: IElementType = IElementType("PRIVATE", SproutLanguage)
        val SERVICE: IElementType = IElementType("SERVICE", SproutLanguage)
        val RPC: IElementType = IElementType("RPC", SproutLanguage)
        
        // HTTP Methods
        val GET: IElementType = IElementType("GET", SproutLanguage)
        val POST: IElementType = IElementType("POST", SproutLanguage)
        val PUT: IElementType = IElementType("PUT", SproutLanguage)
        val DELETE: IElementType = IElementType("DELETE", SproutLanguage)
        val PATCH: IElementType = IElementType("PATCH", SproutLanguage)
        
        // Literals
        val IDENTIFIER: IElementType = IElementType("IDENTIFIER", SproutLanguage)
        val STRING: IElementType = IElementType("STRING", SproutLanguage)
        val RAW_STRING: IElementType = IElementType("RAW_STRING", SproutLanguage)
        val NUMBER: IElementType = IElementType("NUMBER", SproutLanguage)
        
        // Operators and punctuation
        val ARROW: IElementType = IElementType("ARROW", SproutLanguage)  // =>
        val COLON: IElementType = IElementType("COLON", SproutLanguage)  // :
        val LBRACE: IElementType = IElementType("LBRACE", SproutLanguage)  // {
        val RBRACE: IElementType = IElementType("RBRACE", SproutLanguage)  // }
        val LPAREN: IElementType = IElementType("LPAREN", SproutLanguage)  // (
        val RPAREN: IElementType = IElementType("RPAREN", SproutLanguage)  // )
        val LBRACKET: IElementType = IElementType("LBRACKET", SproutLanguage)  // [
        val RBRACKET: IElementType = IElementType("RBRACKET", SproutLanguage)  // ]
        val COMMA: IElementType = IElementType("COMMA", SproutLanguage)
        val SLASH: IElementType = IElementType("SLASH", SproutLanguage)
        
        // Comments
        val COMMENT: IElementType = IElementType("COMMENT", SproutLanguage)
        
        // Whitespace
        val WHITESPACE: IElementType = IElementType("WHITESPACE", SproutLanguage)
        
        // Token sets
        val KEYWORDS = TokenSet.create(TYPE, SERVER, PREFIX, PUBLIC, PRIVATE, SERVICE, RPC)
        val HTTP_METHODS = TokenSet.create(GET, POST, PUT, DELETE, PATCH)
        val OPERATORS = TokenSet.create(ARROW, COLON, SLASH)
        val BRACES = TokenSet.create(LBRACE, RBRACE)
        val PARENTHESES = TokenSet.create(LPAREN, RPAREN)
        val BRACKETS = TokenSet.create(LBRACKET, RBRACKET)
    }

    object Factory {
        fun createElement(node: ASTNode): PsiElement {
            val type = node.elementType
            // Check if this is a leaf token (from SproutTokenTypes) or composite (from SproutElementTypes)
            return if (type == SproutTokenTypes.TYPE || type == SproutTokenTypes.SERVER || 
                       type == SproutTokenTypes.PREFIX || type == SproutTokenTypes.PUBLIC ||
                       type == SproutTokenTypes.PRIVATE || type == SproutTokenTypes.SERVICE ||
                       type == SproutTokenTypes.RPC || type == SproutTokenTypes.GET ||
                       type == SproutTokenTypes.POST || type == SproutTokenTypes.PUT ||
                       type == SproutTokenTypes.DELETE || type == SproutTokenTypes.PATCH ||
                       type == SproutTokenTypes.IDENTIFIER || type == SproutTokenTypes.STRING ||
                       type == SproutTokenTypes.RAW_STRING ||
                       type == SproutTokenTypes.NUMBER || type == SproutTokenTypes.ARROW ||
                       type == SproutTokenTypes.COLON || type == SproutTokenTypes.LBRACE ||
                       type == SproutTokenTypes.RBRACE || type == SproutTokenTypes.LPAREN ||
                       type == SproutTokenTypes.RPAREN || type == SproutTokenTypes.LBRACKET ||
                       type == SproutTokenTypes.RBRACKET || type == SproutTokenTypes.COMMA ||
                       type == SproutTokenTypes.SLASH || type == SproutTokenTypes.COMMENT ||
                       type == SproutTokenTypes.WHITESPACE) {
                com.intellij.psi.impl.source.tree.LeafPsiElement(type, node.text)
            } else {
                ASTWrapperPsiElement(node)
            }
        }
    }
}
