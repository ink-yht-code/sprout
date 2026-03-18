package com.ink.yht.gint.idea.lang;

import com.intellij.lexer.FlexLexer;
import com.intellij.psi.tree.IElementType;

%%

%class GintLexer
%implements FlexLexer
%unicode
%public
%final
%function advance
%type IElementType
%eof{ return null;
%eof}

// Whitespace
WHITE_SPACE=[ \t\r\n]+

// Comments
COMMENT="//"[^\n]*

// Literals
STRING=\"([^\"\n]|\\.)*\"
NUMBER=[0-9]+
IDENTIFIER=[a-zA-Z_][a-zA-Z0-9_]*

%%

<YYINITIAL> {
    // Keywords
    "type"                  { return GintTokenTypes.TYPE; }
    "server"                { return GintTokenTypes.SERVER; }
    "prefix"                { return GintTokenTypes.PREFIX; }
    "public"                { return GintTokenTypes.PUBLIC; }
    "private"               { return GintTokenTypes.PRIVATE; }
    "service"               { return GintTokenTypes.SERVICE; }
    "rpc"                   { return GintTokenTypes.RPC; }
    
    // HTTP Methods
    "GET"                   { return GintTokenTypes.GET; }
    "POST"                  { return GintTokenTypes.POST; }
    "PUT"                   { return GintTokenTypes.PUT; }
    "DELETE"                { return GintTokenTypes.DELETE; }
    "PATCH"                 { return GintTokenTypes.PATCH; }
    
    // Operators and punctuation
    "=>"                    { return GintTokenTypes.ARROW; }
    ":"                     { return GintTokenTypes.COLON; }
    "{"                     { return GintTokenTypes.LBRACE; }
    "}"                     { return GintTokenTypes.RBRACE; }
    "("                     { return GintTokenTypes.LPAREN; }
    ")"                     { return GintTokenTypes.RPAREN; }
    "["                     { return GintTokenTypes.LBRACKET; }
    "]"                     { return GintTokenTypes.RBRACKET; }
    ","                     { return GintTokenTypes.COMMA; }
    "/"                     { return GintTokenTypes.SLASH; }
    
    // Literals
    {STRING}                { return GintTokenTypes.STRING; }
    {NUMBER}                { return GintTokenTypes.NUMBER; }
    {IDENTIFIER}            { return GintTokenTypes.IDENTIFIER; }
    
    // Comments
    {COMMENT}               { return GintTokenTypes.COMMENT; }
    
    // Whitespace
    {WHITE_SPACE}           { return GintTokenTypes.WHITESPACE; }
}

. { return com.intellij.psi.TokenType.BAD_CHARACTER; }
