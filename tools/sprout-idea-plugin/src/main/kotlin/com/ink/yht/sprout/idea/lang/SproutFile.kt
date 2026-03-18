package com.ink.yht.sprout.idea.lang

import com.intellij.extapi.psi.PsiFileBase
import com.intellij.psi.FileViewProvider

class SproutFile(viewProvider: FileViewProvider) : PsiFileBase(viewProvider, SproutLanguage) {
    override fun getFileType() = SproutFileType
    override fun toString() = "Sprout File"
}
