package com.ink.yht.sprout.idea.completion

import com.intellij.codeInsight.completion.CompletionContributor
import com.intellij.codeInsight.completion.CompletionParameters
import com.intellij.codeInsight.completion.CompletionProvider
import com.intellij.codeInsight.completion.CompletionResultSet
import com.intellij.codeInsight.completion.CompletionType
import com.intellij.codeInsight.lookup.LookupElementBuilder
import com.intellij.patterns.PlatformPatterns
import com.intellij.util.ProcessingContext
import com.ink.yht.sprout.idea.lang.SproutLanguage

class SproutCompletionContributor : CompletionContributor() {
    init {
        // Keyword completion
        extend(
            CompletionType.BASIC,
            PlatformPatterns.psiElement().withLanguage(SproutLanguage),
            KeywordCompletionProvider
        )
        
        // HTTP method completion
        extend(
            CompletionType.BASIC,
            PlatformPatterns.psiElement().withLanguage(SproutLanguage),
            HttpMethodCompletionProvider
        )

        // Struct tag template completion (Go-like), triggered by backtick
        extend(
            CompletionType.BASIC,
            PlatformPatterns.psiElement().withLanguage(SproutLanguage),
            StructTagTemplateCompletionProvider
        )
    }
}

object StructTagTemplateCompletionProvider : CompletionProvider<CompletionParameters>() {
    private val TEMPLATES = listOf(
        "`json:\"name\"`",
        "`json:\"name\" validate:\"required\"`",
        "`json:\"name\" validate:\"required,email\"`",
    )

    override fun addCompletions(
        parameters: CompletionParameters,
        context: ProcessingContext,
        result: CompletionResultSet
    ) {
        val offset = parameters.offset
        val fileText = parameters.originalFile.text
        if (offset <= 0 || offset > fileText.length) return

        if (fileText[offset - 1] != '`') return

        TEMPLATES.forEach { template ->
            result.addElement(
                LookupElementBuilder.create(template)
                    .withTypeText("struct tag")
            )
        }
    }
}

object KeywordCompletionProvider : CompletionProvider<CompletionParameters>() {
    private val KEYWORDS = listOf(
        "type", "server", "prefix", "public", "private", "service", "rpc"
    )
    
    override fun addCompletions(
        parameters: CompletionParameters,
        context: ProcessingContext,
        result: CompletionResultSet
    ) {
        val prefix = result.prefixMatcher.prefix
        val r = result.withPrefixMatcher(prefix)

        KEYWORDS
            .asSequence()
            .filter { prefix.isEmpty() || it.startsWith(prefix, ignoreCase = true) }
            .forEach { keyword ->
                r.addElement(
                    LookupElementBuilder.create(keyword)
                        .withBoldness(true)
                        .withTypeText("keyword")
                )
            }
    }
}

object HttpMethodCompletionProvider : CompletionProvider<CompletionParameters>() {
    private val HTTP_METHODS = listOf("GET", "POST", "PUT", "DELETE", "PATCH")
    
    override fun addCompletions(
        parameters: CompletionParameters,
        context: ProcessingContext,
        result: CompletionResultSet
    ) {
        val position = parameters.position

        val prefix = result.prefixMatcher.prefix
        val r = result.withPrefixMatcher(prefix)
        
        // Check if we're in a service block context
        // This is a simplified check - in production you'd want more precise context detection
        val file = position.containingFile
        val text = file.text
        
        // Simple heuristic: if we're inside a service block, offer HTTP methods
        val cursorOffset = parameters.offset
        val beforeCursor = text.substring(0, cursorOffset)
        
        if (beforeCursor.contains("service") && beforeCursor.contains("{")) {
            val afterLastBrace = beforeCursor.lastIndexOf("{")
            val afterLastClose = beforeCursor.lastIndexOf("}")
            
            if (afterLastBrace > afterLastClose) {
                // We're inside a service block
                HTTP_METHODS
                    .asSequence()
                    .filter { prefix.isEmpty() || it.startsWith(prefix, ignoreCase = true) }
                    .forEach { method ->
                        r.addElement(
                            LookupElementBuilder.create(method)
                                .withBoldness(true)
                                .withTypeText("HTTP method")
                        )
                    }
            }
        }
    }
}
