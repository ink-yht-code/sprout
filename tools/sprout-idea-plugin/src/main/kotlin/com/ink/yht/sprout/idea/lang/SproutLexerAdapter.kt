package com.ink.yht.sprout.idea.lang

import com.intellij.lexer.LexerBase

class SproutLexerAdapter : LexerBase() {
    private val lexer = SproutLexer()
    
    override fun start(buffer: CharSequence, startOffset: Int, endOffset: Int, initialState: Int) {
        lexer.start(buffer, startOffset, endOffset, initialState)
    }
    
    override fun getState() = lexer.state
    override fun getTokenType() = lexer.tokenType
    override fun getTokenStart() = lexer.tokenStart
    override fun getTokenEnd() = lexer.tokenEnd
    override fun advance() = lexer.advance()
    override fun getBufferSequence() = lexer.bufferSequence
    override fun getBufferEnd() = lexer.bufferEnd
}
