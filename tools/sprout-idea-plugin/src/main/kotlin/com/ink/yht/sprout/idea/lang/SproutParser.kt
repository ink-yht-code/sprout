package com.ink.yht.sprout.idea.lang

import com.intellij.lang.ASTNode
import com.intellij.lang.PsiBuilder
import com.intellij.lang.PsiParser
import com.intellij.lang.parser.GeneratedParserUtilBase
import com.intellij.psi.tree.IElementType

class SproutParser : PsiParser, GeneratedParserUtilBase() {
    override fun parse(root: IElementType, builder: PsiBuilder): ASTNode {
        val marker = builder.mark()
        
        while (!builder.eof()) {
            val tokenType = builder.tokenType
            
            when {
                tokenType == SproutTokenTypes.TYPE -> parseTypeDefinition(builder)
                tokenType == SproutTokenTypes.SERVER -> parseServerDefinition(builder)
                tokenType == SproutTokenTypes.PUBLIC || tokenType == SproutTokenTypes.PRIVATE -> parseServiceBlock(builder)
                tokenType == SproutTokenTypes.COMMENT -> builder.advanceLexer()
                tokenType == SproutTokenTypes.WHITESPACE -> builder.advanceLexer()
                else -> builder.advanceLexer()
            }
        }
        
        marker.done(root)
        return builder.treeBuilt
    }
    
    private fun parseTypeDefinition(builder: PsiBuilder) {
        val marker = builder.mark()
        builder.advanceLexer() // 'type'
        
        if (builder.tokenType == SproutTokenTypes.IDENTIFIER) {
            builder.advanceLexer() // type name
        }
        
        if (builder.tokenType == SproutTokenTypes.LBRACE) {
            builder.advanceLexer() // '{'
            parseTypeFields(builder)
            if (builder.tokenType == SproutTokenTypes.RBRACE) {
                builder.advanceLexer() // '}'
            }
        }
        
        marker.done(SproutElementTypes.TYPE_DEFINITION)
    }
    
    private fun parseTypeFields(builder: PsiBuilder) {
        while (builder.tokenType != SproutTokenTypes.RBRACE && !builder.eof()) {
            if (builder.tokenType == SproutTokenTypes.IDENTIFIER) {
                val fieldMarker = builder.mark()
                
                // Field name
                val nameMarker = builder.mark()
                builder.advanceLexer() // field name
                nameMarker.done(SproutElementTypes.FIELD_NAME)
                
                if (builder.tokenType == SproutTokenTypes.COLON) {
                    builder.advanceLexer() // ':'
                }
                
                // Field type
                if (builder.tokenType != null && 
                    builder.tokenType != SproutTokenTypes.RBRACE &&
                    builder.tokenType != SproutTokenTypes.COMMA &&
                    builder.tokenType != SproutTokenTypes.RAW_STRING) {
                    val typeMarker = builder.mark()
                    builder.advanceLexer()
                    typeMarker.done(SproutElementTypes.FIELD_TYPE)
                }
                
                // Field tag (raw string)
                if (builder.tokenType == SproutTokenTypes.RAW_STRING) {
                    val tagMarker = builder.mark()
                    builder.advanceLexer()
                    tagMarker.done(SproutElementTypes.FIELD_TAG)
                }
                
                if (builder.tokenType == SproutTokenTypes.COMMA) {
                    builder.advanceLexer()
                }
                
                fieldMarker.done(SproutElementTypes.TYPE_FIELD)
            } else {
                builder.advanceLexer()
            }
        }
    }
    
    private fun parseServerDefinition(builder: PsiBuilder) {
        val marker = builder.mark()
        builder.advanceLexer() // 'server'
        
        if (builder.tokenType == SproutTokenTypes.LBRACE) {
            builder.advanceLexer() // '{'
            parseServerFields(builder)
            if (builder.tokenType == SproutTokenTypes.RBRACE) {
                builder.advanceLexer() // '}'
            }
        }
        
        marker.done(SproutElementTypes.SERVER_DEFINITION)
    }
    
    private fun parseServerFields(builder: PsiBuilder) {
        while (builder.tokenType != SproutTokenTypes.RBRACE && !builder.eof()) {
            if (builder.tokenType == SproutTokenTypes.PREFIX) {
                val fieldMarker = builder.mark()
                builder.advanceLexer() // 'prefix'
                
                if (builder.tokenType == SproutTokenTypes.COLON) {
                    builder.advanceLexer() // ':'
                }
                
                if (builder.tokenType == SproutTokenTypes.STRING || 
                    builder.tokenType == SproutTokenTypes.IDENTIFIER) {
                    builder.advanceLexer() // value
                }
                
                fieldMarker.done(SproutElementTypes.SERVER_FIELD)
            } else {
                builder.advanceLexer()
            }
        }
    }
    
    private fun parseServiceBlock(builder: PsiBuilder) {
        val marker = builder.mark()
        val visibility = builder.tokenType
        builder.advanceLexer() // 'public' or 'private'
        
        if (builder.tokenType == SproutTokenTypes.SERVICE) {
            builder.advanceLexer() // 'service'
        }
        
        if (builder.tokenType == SproutTokenTypes.IDENTIFIER) {
            builder.advanceLexer() // service name
        }
        
        if (builder.tokenType == SproutTokenTypes.LBRACE) {
            builder.advanceLexer() // '{'
            parseServiceMethods(builder)
            if (builder.tokenType == SproutTokenTypes.RBRACE) {
                builder.advanceLexer() // '}'
            }
        }
        
        marker.done(SproutElementTypes.SERVICE_BLOCK)
    }
    
    private fun parseServiceMethods(builder: PsiBuilder) {
        while (builder.tokenType != SproutTokenTypes.RBRACE && !builder.eof()) {
            val tokenType = builder.tokenType
            
            if (SproutTokenTypes.HTTP_METHODS.contains(tokenType)) {
                parseHttpEndpoint(builder)
            } else if (tokenType == SproutTokenTypes.RPC) {
                parseRpcMethod(builder)
            } else if (tokenType == SproutTokenTypes.COMMENT || tokenType == SproutTokenTypes.WHITESPACE) {
                builder.advanceLexer()
            } else {
                builder.advanceLexer()
            }
        }
    }
    
    private fun parseHttpEndpoint(builder: PsiBuilder) {
        val marker = builder.mark()
        builder.advanceLexer() // HTTP method
        
        // Parse path
        while (builder.tokenType == SproutTokenTypes.SLASH || 
               builder.tokenType == SproutTokenTypes.IDENTIFIER ||
               builder.tokenType == SproutTokenTypes.STRING) {
            builder.advanceLexer()
        }
        
        // Parse parameters
        if (builder.tokenType == SproutTokenTypes.LPAREN) {
            builder.advanceLexer() // '('
            if (builder.tokenType == SproutTokenTypes.IDENTIFIER) {
                builder.advanceLexer() // request type
            }
            if (builder.tokenType == SproutTokenTypes.RPAREN) {
                builder.advanceLexer() // ')'
            }
        }
        
        // Parse return type
        if (builder.tokenType == SproutTokenTypes.ARROW) {
            builder.advanceLexer() // '=>'
            if (builder.tokenType == SproutTokenTypes.IDENTIFIER) {
                builder.advanceLexer() // response type
            }
        }
        
        marker.done(SproutElementTypes.HTTP_ENDPOINT)
    }
    
    private fun parseRpcMethod(builder: PsiBuilder) {
        val marker = builder.mark()
        builder.advanceLexer() // 'rpc'
        
        if (builder.tokenType == SproutTokenTypes.IDENTIFIER) {
            builder.advanceLexer() // method name
        }
        
        if (builder.tokenType == SproutTokenTypes.LPAREN) {
            builder.advanceLexer() // '('
            if (builder.tokenType == SproutTokenTypes.IDENTIFIER) {
                builder.advanceLexer() // request type
            }
            if (builder.tokenType == SproutTokenTypes.RPAREN) {
                builder.advanceLexer() // ')'
            }
        }
        
        if (builder.tokenType == SproutTokenTypes.ARROW) {
            builder.advanceLexer() // '=>'
            if (builder.tokenType == SproutTokenTypes.IDENTIFIER) {
                builder.advanceLexer() // response type
            }
        }
        
        marker.done(SproutElementTypes.RPC_METHOD)
    }
}
