package com.ink.yht.sprout.idea.formatting

import com.intellij.formatting.Alignment
import com.intellij.formatting.Block
import com.intellij.formatting.ChildAttributes
import com.intellij.formatting.Indent
import com.intellij.formatting.Spacing
import com.intellij.formatting.SpacingBuilder
import com.intellij.formatting.Wrap
import com.intellij.lang.ASTNode
import com.intellij.psi.TokenType
import com.intellij.psi.formatter.common.AbstractBlock
import com.ink.yht.sprout.idea.lang.SproutElementTypes
import com.ink.yht.sprout.idea.lang.SproutTokenTypes

class SproutBlock(
    node: ASTNode,
    private val settings: com.intellij.psi.codeStyle.CodeStyleSettings,
    private val spacingBuilder: SpacingBuilder,
    wrap: Wrap?,
    alignment: Alignment?,
    private val indent: Indent? = null,
    private val insideBraces: Boolean = false,
) : AbstractBlock(node, wrap, alignment) {

    override fun buildChildren(): List<Block> {
        val blocks = ArrayList<Block>()
        
        // For TYPE_DEFINITION, calculate alignment for fields
        if (myNode.elementType == SproutElementTypes.TYPE_DEFINITION) {
            return buildTypeDefinitionChildren()
        }
        
        var child = myNode.firstChildNode
        while (child != null) {
            val type = child.elementType
            if (type != TokenType.WHITE_SPACE) {
                val childIndent = if (insideBraces && type != SproutTokenTypes.RBRACE) {
                    Indent.getNormalIndent()
                } else {
                    Indent.getNoneIndent()
                }
                val childInsideBraces = insideBraces || type == SproutTokenTypes.LBRACE
                blocks.add(
                    SproutBlock(
                        child,
                        settings,
                        spacingBuilder,
                        null,
                        null,
                        childIndent,
                        childInsideBraces,
                    )
                )
            }
            child = child.treeNext
        }
        return blocks
    }
    
    private fun buildTypeDefinitionChildren(): List<Block> {
        val blocks = ArrayList<Block>()
        
        // Collect all TYPE_FIELD nodes to calculate alignment
        val typeFields = mutableListOf<ASTNode>()
        var child = myNode.firstChildNode
        while (child != null) {
            if (child.elementType == SproutElementTypes.TYPE_FIELD) {
                typeFields.add(child)
            }
            child = child.treeNext
        }
        
        // Calculate max widths for alignment
        var maxNameWidth = 0
        var maxTypeWidth = 0
        for (field in typeFields) {
            val nameNode = field.findChildByType(SproutElementTypes.FIELD_NAME)
            val typeNode = field.findChildByType(SproutElementTypes.FIELD_TYPE)
            if (nameNode != null) {
                maxNameWidth = maxOf(maxNameWidth, nameNode.textLength)
            }
            if (typeNode != null) {
                maxTypeWidth = maxOf(maxTypeWidth, typeNode.textLength)
            }
        }
        
        // Create alignments
        val nameAlignment = Alignment.createAlignment()
        val typeAlignment = Alignment.createAlignment()
        val tagAlignment = Alignment.createAlignment()
        
        // Build children with alignment
        child = myNode.firstChildNode
        while (child != null) {
            val type = child.elementType
            if (type != TokenType.WHITE_SPACE) {
                val childIndent = if (type != SproutTokenTypes.RBRACE && type != SproutTokenTypes.TYPE) {
                    Indent.getNormalIndent()
                } else {
                    Indent.getNoneIndent()
                }
                
                if (type == SproutElementTypes.TYPE_FIELD) {
                    blocks.add(buildTypeFieldBlock(child, nameAlignment, typeAlignment, tagAlignment, childIndent))
                } else {
                    blocks.add(
                        SproutBlock(
                            child,
                            settings,
                            spacingBuilder,
                            null,
                            null,
                            childIndent,
                            true,
                        )
                    )
                }
            }
            child = child.treeNext
        }
        return blocks
    }
    
    private fun buildTypeFieldBlock(
        node: ASTNode,
        nameAlignment: Alignment,
        typeAlignment: Alignment,
        tagAlignment: Alignment,
        fieldIndent: Indent?
    ): Block {
        val blocks = ArrayList<Block>()
        var child = node.firstChildNode
        while (child != null) {
            val type = child.elementType
            if (type != TokenType.WHITE_SPACE) {
                val alignment = when (type) {
                    SproutElementTypes.FIELD_NAME -> nameAlignment
                    SproutElementTypes.FIELD_TYPE -> typeAlignment
                    SproutElementTypes.FIELD_TAG -> tagAlignment
                    else -> null
                }
                blocks.add(
                    SproutBlock(
                        child,
                        settings,
                        spacingBuilder,
                        null,
                        alignment,
                        fieldIndent,
                        true,
                    )
                )
            }
            child = child.treeNext
        }
        return object : AbstractBlock(node, null, null) {
            override fun buildChildren(): List<Block> = blocks
            override fun getSpacing(child1: Block?, child2: Block): Spacing? = spacingBuilder.getSpacing(this, child1, child2)
            override fun getChildAttributes(newChildIndex: Int): ChildAttributes = ChildAttributes(Indent.getNormalIndent(), null)
            override fun isLeaf(): Boolean = false
            override fun getIndent(): Indent? = fieldIndent
        }
    }

    override fun getIndent(): Indent? = indent

    override fun getSpacing(child1: Block?, child2: Block): Spacing? {
        return spacingBuilder.getSpacing(this, child1, child2)
    }

    override fun getChildAttributes(newChildIndex: Int): ChildAttributes {
        return ChildAttributes(Indent.getNormalIndent(), null)
    }

    override fun isLeaf(): Boolean = myNode.firstChildNode == null
}
