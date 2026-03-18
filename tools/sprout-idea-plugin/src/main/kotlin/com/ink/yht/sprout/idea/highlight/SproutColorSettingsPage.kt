package com.ink.yht.sprout.idea.highlight

import com.intellij.openapi.editor.colors.TextAttributesKey
import com.intellij.openapi.options.colors.AttributesDescriptor
import com.intellij.openapi.options.colors.ColorDescriptor
import com.intellij.openapi.options.colors.ColorSettingsPage
import com.ink.yht.sprout.idea.lang.SproutIcons
import javax.swing.Icon

class SproutColorSettingsPage : ColorSettingsPage {
    companion object {
        private val DESCRIPTORS = arrayOf(
            AttributesDescriptor("Keyword", SproutSyntaxHighlighter.KEYWORD),
            AttributesDescriptor("HTTP Method", SproutSyntaxHighlighter.HTTP_METHOD),
            AttributesDescriptor("String", SproutSyntaxHighlighter.STRING),
            AttributesDescriptor("Number", SproutSyntaxHighlighter.NUMBER),
            AttributesDescriptor("Identifier", SproutSyntaxHighlighter.IDENTIFIER),
            AttributesDescriptor("Operator", SproutSyntaxHighlighter.OPERATOR),
            AttributesDescriptor("Comment", SproutSyntaxHighlighter.COMMENT),
            AttributesDescriptor("Braces", SproutSyntaxHighlighter.BRACES),
            AttributesDescriptor("Parentheses", SproutSyntaxHighlighter.PARENTHESES),
            AttributesDescriptor("Brackets", SproutSyntaxHighlighter.BRACKETS),
            AttributesDescriptor("Bad character", SproutSyntaxHighlighter.BAD_CHARACTER),
        )
    }

    override fun getIcon(): Icon? = SproutIcons.FILE
    override fun getHighlighter() = SproutSyntaxHighlighter()
    override fun getDemoText(): String = """
// Sprout API definition file
type HelloReq {
    Name string `json:"name"`
}

type HelloResp {
    Message string `json:"message"`
}

server {
    prefix: "/api"
}

public service HelloService {
    GET /hello(HelloReq) => HelloResp
    POST /greet(HelloReq) => HelloResp
}

private service InternalService {
    rpc Ping(PingReq) => PingResp
}
""".trimIndent()

    override fun getAdditionalHighlightingTagToDescriptorMap(): Map<String, TextAttributesKey>? = null
    override fun getAttributeDescriptors() = DESCRIPTORS
    override fun getColorDescriptors(): Array<ColorDescriptor> = ColorDescriptor.EMPTY_ARRAY
    override fun getDisplayName() = "Sprout"
}
