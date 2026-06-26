"use client"

import { useEffect, useRef, useCallback } from "react"
import { useEditor, EditorContent } from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import { Table } from "@tiptap/extension-table"
import TableRow from "@tiptap/extension-table-row"
import TableHeader from "@tiptap/extension-table-header"
import TableCell from "@tiptap/extension-table-cell"
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight"
import Highlight from "@tiptap/extension-highlight"
import TaskList from "@tiptap/extension-task-list"
import TaskItem from "@tiptap/extension-task-item"
import Link from "@tiptap/extension-link"
import Underline from "@tiptap/extension-underline"
import Subscript from "@tiptap/extension-subscript"
import Superscript from "@tiptap/extension-superscript"
import Placeholder from "@tiptap/extension-placeholder"
import CharacterCount from "@tiptap/extension-character-count"
import { common, createLowlight } from "lowlight"
import { Toolbar } from "./toolbar"
import { MathExtension, DisplayMathExtension } from "./extensions/math"
import { MermaidExtension } from "./extensions/mermaid"

const lowlight = createLowlight(common)

interface KvasirEditorProps {
  content: string
  onChange?: (html: string) => void
  onSave?: (html: string) => void
  placeholder?: string
  readOnly?: boolean
  showToolbar?: boolean
}

export function KvasirEditor({
  content,
  onChange,
  onSave,
  placeholder = "Start writing...",
  readOnly = false,
  showToolbar = true,
}: KvasirEditorProps) {
  const saveTimerRef = useRef<ReturnType<typeof setTimeout>>(undefined)

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        codeBlock: false,
      }),
      Table.configure({ resizable: true }),
      TableRow,
      TableHeader,
      TableCell,
      CodeBlockLowlight.configure({ lowlight }),
      Highlight,
      TaskList,
      TaskItem.configure({ nested: true }),
      Link.configure({
        openOnClick: true,
        HTMLAttributes: { class: "text-primary underline" },
      }),
      Underline,
      Subscript,
      Superscript,
      Placeholder.configure({ placeholder }),
      CharacterCount,
      MathExtension,
      DisplayMathExtension,
      MermaidExtension,
    ],
    content,
    editable: !readOnly,
    editorProps: {
      attributes: {
        class:
          "prose prose-sm dark:prose-invert max-w-none focus:outline-none min-h-[300px] px-6 py-4",
      },
    },
    onUpdate: ({ editor }) => {
      const html = editor.getHTML()
      onChange?.(html)

      if (saveTimerRef.current) clearTimeout(saveTimerRef.current)
      saveTimerRef.current = setTimeout(() => {
        onSave?.(html)
      }, 2000)
    },
  })

  useEffect(() => {
    return () => {
      if (saveTimerRef.current) clearTimeout(saveTimerRef.current)
    }
  }, [])

  const handleSave = useCallback(() => {
    if (editor) {
      const html = editor.getHTML()
      if (saveTimerRef.current) clearTimeout(saveTimerRef.current)
      onSave?.(html)
    }
  }, [editor, onSave])

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "s") {
        e.preventDefault()
        handleSave()
      }
    }
    window.addEventListener("keydown", handler)
    return () => window.removeEventListener("keydown", handler)
  }, [handleSave])

  if (!editor) return null

  return (
    <div className="flex flex-col h-full">
      {showToolbar && !readOnly && <Toolbar editor={editor} />}
      <div className="flex-1 overflow-auto">
        <EditorContent editor={editor} />
      </div>
      <div className="px-6 py-1.5 border-t text-xs text-muted-foreground">
        {editor.storage.characterCount.characters()} characters
      </div>
    </div>
  )
}
