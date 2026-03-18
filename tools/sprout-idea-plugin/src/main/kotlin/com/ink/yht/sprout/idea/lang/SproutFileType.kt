package com.ink.yht.sprout.idea.lang

import com.intellij.openapi.fileTypes.LanguageFileType
import javax.swing.Icon

object SproutFileType : LanguageFileType(SproutLanguage) {
    override fun getName() = "Sprout"
    override fun getDescription() = "Sprout API definition file"
    override fun getDefaultExtension() = "sprout"
    override fun getIcon(): Icon = SproutIcons.FILE
}
