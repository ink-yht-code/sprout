package com.ink.yht.sprout.idea.lang

import com.intellij.psi.tree.IElementType

object SproutElementTypes {
    val TYPE_DEFINITION = IElementType("TYPE_DEFINITION", SproutLanguage)
    val TYPE_FIELD = IElementType("TYPE_FIELD", SproutLanguage)
    val FIELD_NAME = IElementType("FIELD_NAME", SproutLanguage)
    val FIELD_TYPE = IElementType("FIELD_TYPE", SproutLanguage)
    val FIELD_TAG = IElementType("FIELD_TAG", SproutLanguage)
    val SERVER_DEFINITION = IElementType("SERVER_DEFINITION", SproutLanguage)
    val SERVER_FIELD = IElementType("SERVER_FIELD", SproutLanguage)
    val SERVICE_BLOCK = IElementType("SERVICE_BLOCK", SproutLanguage)
    val HTTP_ENDPOINT = IElementType("HTTP_ENDPOINT", SproutLanguage)
    val RPC_METHOD = IElementType("RPC_METHOD", SproutLanguage)
}
