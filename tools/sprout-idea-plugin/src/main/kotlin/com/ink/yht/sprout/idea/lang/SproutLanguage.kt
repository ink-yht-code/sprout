package com.ink.yht.sprout.idea.lang

import com.intellij.lang.Language

object SproutLanguage : Language("Sprout") {
    override fun isCaseSensitive() = true
    override fun getDisplayName() = "Sprout"
}
