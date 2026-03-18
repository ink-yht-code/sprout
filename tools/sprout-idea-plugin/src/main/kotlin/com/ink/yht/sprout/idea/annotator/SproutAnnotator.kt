package com.ink.yht.sprout.idea.annotator

import com.intellij.lang.annotation.AnnotationHolder
import com.intellij.lang.annotation.Annotator
import com.intellij.lang.annotation.HighlightSeverity
import com.intellij.psi.PsiElement
import com.intellij.psi.impl.source.tree.LeafPsiElement
import com.ink.yht.sprout.idea.lang.SproutTokenTypes

class SproutAnnotator : Annotator {
    override fun annotate(element: PsiElement, holder: AnnotationHolder) {
        // Only process leaf elements (tokens)
        if (element !is LeafPsiElement) return
        
        val tokenType = element.node.elementType
        
        // Check for common syntax issues
        when (tokenType) {
            SproutTokenTypes.IDENTIFIER -> {
                // Check if identifier looks like a type reference (capitalized)
                val text = element.text
                if (text.isNotEmpty() && text[0].isUpperCase()) {
                    // Mark as potential type reference - could add reference checking here
                    // For now, just highlight it
                }
            }
            
            SproutTokenTypes.ARROW -> {
                // Verify arrow has valid context (should be followed by identifier)
                val nextSibling = element.node.treeNext
                if (nextSibling == null || nextSibling.elementType != SproutTokenTypes.IDENTIFIER) {
                    holder.newAnnotation(
                        HighlightSeverity.WARNING,
                        "Expected return type after '=>'"
                    ).range(element.textRange).create()
                }
            }
            
            SproutTokenTypes.RBRACE -> {
                // Check for unmatched closing brace
            }
            
            SproutTokenTypes.LBRACE -> {
                // Check for unclosed opening brace
            }
        }
    }
}
