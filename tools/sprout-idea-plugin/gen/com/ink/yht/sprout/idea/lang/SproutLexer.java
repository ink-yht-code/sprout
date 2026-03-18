package com.ink.yht.sprout.idea.lang;

import com.intellij.lexer.LexerBase;
import com.intellij.psi.tree.IElementType;
import org.jetbrains.annotations.NotNull;

/**
 * Hand-written lexer for Sprout language.
 * This is a simple implementation that avoids the need for JFlex.
 */
public class SproutLexer extends LexerBase {
    private CharSequence buffer;
    private int startOffset;
    private int endOffset;
    private int position;
    private IElementType tokenType;
    private int tokenStart;
    private int tokenEnd;

    @Override
    public void start(@NotNull CharSequence buffer, int startOffset, int endOffset, int initialState) {
        this.buffer = buffer;
        this.startOffset = startOffset;
        this.endOffset = endOffset;
        this.position = startOffset;
        this.tokenType = null;
        this.tokenStart = startOffset;
        this.tokenEnd = startOffset;
        advance();
    }

    @Override
    public int getState() {
        return 0;
    }

    @Override
    public IElementType getTokenType() {
        return tokenType;
    }

    @Override
    public int getTokenStart() {
        return tokenStart;
    }

    @Override
    public int getTokenEnd() {
        return tokenEnd;
    }

    @Override
    public void advance() {
        if (position >= endOffset) {
            tokenType = null;
            return;
        }

        tokenStart = position;
        char c = buffer.charAt(position);

        // Whitespace
        if (Character.isWhitespace(c)) {
            while (position < endOffset && Character.isWhitespace(buffer.charAt(position))) {
                position++;
            }
            tokenEnd = position;
            tokenType = SproutTokenTypes.Companion.getWHITESPACE();
            return;
        }

        // Comment
        if (c == '/' && position + 1 < endOffset && buffer.charAt(position + 1) == '/') {
            while (position < endOffset && buffer.charAt(position) != '\n') {
                position++;
            }
            tokenEnd = position;
            tokenType = SproutTokenTypes.Companion.getCOMMENT();
            return;
        }

        // String literal
        if (c == '"') {
            position++;
            while (position < endOffset) {
                char ch = buffer.charAt(position);
                if (ch == '"') {
                    position++;
                    break;
                }
                if (ch == '\\' && position + 1 < endOffset) {
                    position += 2;
                } else {
                    position++;
                }
            }
            tokenEnd = position;
            tokenType = SproutTokenTypes.Companion.getSTRING();
            return;
        }

        // Raw string literal (Go-like backtick string), used for struct tags
        if (c == '`') {
            position++;
            while (position < endOffset) {
                char ch = buffer.charAt(position);
                if (ch == '`') {
                    position++;
                    break;
                }
                position++;
            }
            tokenEnd = position;
            tokenType = SproutTokenTypes.Companion.getRAW_STRING();
            return;
        }

        // Arrow operator
        if (c == '=' && position + 1 < endOffset && buffer.charAt(position + 1) == '>') {
            position += 2;
            tokenEnd = position;
            tokenType = SproutTokenTypes.Companion.getARROW();
            return;
        }

        // Single character tokens
        switch (c) {
            case ':':
                position++;
                tokenEnd = position;
                tokenType = SproutTokenTypes.Companion.getCOLON();
                return;
            case '{':
                position++;
                tokenEnd = position;
                tokenType = SproutTokenTypes.Companion.getLBRACE();
                return;
            case '}':
                position++;
                tokenEnd = position;
                tokenType = SproutTokenTypes.Companion.getRBRACE();
                return;
            case '(':
                position++;
                tokenEnd = position;
                tokenType = SproutTokenTypes.Companion.getLPAREN();
                return;
            case ')':
                position++;
                tokenEnd = position;
                tokenType = SproutTokenTypes.Companion.getRPAREN();
                return;
            case '[':
                position++;
                tokenEnd = position;
                tokenType = SproutTokenTypes.Companion.getLBRACKET();
                return;
            case ']':
                position++;
                tokenEnd = position;
                tokenType = SproutTokenTypes.Companion.getRBRACKET();
                return;
            case ',':
                position++;
                tokenEnd = position;
                tokenType = SproutTokenTypes.Companion.getCOMMA();
                return;
            case '/':
                position++;
                tokenEnd = position;
                tokenType = SproutTokenTypes.Companion.getSLASH();
                return;
        }

        // Number
        if (Character.isDigit(c)) {
            while (position < endOffset && Character.isDigit(buffer.charAt(position))) {
                position++;
            }
            tokenEnd = position;
            tokenType = SproutTokenTypes.Companion.getNUMBER();
            return;
        }

        // Identifier or keyword
        if (Character.isLetter(c) || c == '_') {
            int start = position;
            while (position < endOffset) {
                char ch = buffer.charAt(position);
                if (Character.isLetterOrDigit(ch) || ch == '_') {
                    position++;
                } else {
                    break;
                }
            }
            tokenEnd = position;
            String text = buffer.subSequence(start, position).toString();
            tokenType = getKeywordTokenType(text);
            return;
        }

        // Bad character
        position++;
        tokenEnd = position;
        tokenType = com.intellij.psi.TokenType.BAD_CHARACTER;
    }

    private IElementType getKeywordTokenType(String text) {
        switch (text) {
            case "type": return SproutTokenTypes.Companion.getTYPE();
            case "server": return SproutTokenTypes.Companion.getSERVER();
            case "prefix": return SproutTokenTypes.Companion.getPREFIX();
            case "public": return SproutTokenTypes.Companion.getPUBLIC();
            case "private": return SproutTokenTypes.Companion.getPRIVATE();
            case "service": return SproutTokenTypes.Companion.getSERVICE();
            case "rpc": return SproutTokenTypes.Companion.getRPC();
            case "GET": return SproutTokenTypes.Companion.getGET();
            case "POST": return SproutTokenTypes.Companion.getPOST();
            case "PUT": return SproutTokenTypes.Companion.getPUT();
            case "DELETE": return SproutTokenTypes.Companion.getDELETE();
            case "PATCH": return SproutTokenTypes.Companion.getPATCH();
            default: return SproutTokenTypes.Companion.getIDENTIFIER();
        }
    }

    @Override
    public @NotNull CharSequence getBufferSequence() {
        return buffer;
    }

    @Override
    public int getBufferEnd() {
        return endOffset;
    }
}
