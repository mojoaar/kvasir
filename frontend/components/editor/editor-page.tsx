"use client"

import { useEffect, useRef, useState, useCallback } from "react"
import { KvasirEditor } from "./kvasir-editor"
import { Button } from "@/components/ui/button"
import { Eye, Pencil, Columns3 } from "lucide-react"
import { cn } from "@/lib/utils"

interface EditorPageProps {
  initialContent?: string
  onSave?: (html: string) => void
  readOnly?: boolean
}

type ViewMode = "edit" | "preview" | "split"

export function EditorPage({
  initialContent = "",
  onSave,
  readOnly = false,
}: EditorPageProps) {
  const [content, setContent] = useState(initialContent)
  const [viewMode, setViewMode] = useState<ViewMode>(readOnly ? "preview" : "edit")
  const previewRef = useRef<HTMLDivElement>(null)

  const updatePreview = useCallback((html: string) => {
    if (previewRef.current) {
      previewRef.current.innerHTML = html
      previewRef.current.querySelectorAll("pre code").forEach((block) => {
        if (block.parentElement?.tagName === "PRE") {
          block.parentElement.classList.add(
            "rounded-lg",
            "bg-muted",
            "p-4",
            "overflow-x-auto",
            "text-sm"
          )
        }
      })
    }
  }, [])

  useEffect(() => {
    if (content) updatePreview(content)
  }, [content, updatePreview])

  const handleSave = useCallback(
    (html: string) => {
      setContent(html)
      onSave?.(html)
    },
    [onSave]
  )

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center gap-1 px-4 py-1.5 border-b bg-background/50">
        <Button
          variant="ghost"
          size="sm"
          className={cn("h-7 text-xs gap-1.5", viewMode === "edit" && "bg-accent")}
          onClick={() => setViewMode("edit")}
        >
          <Pencil className="h-3.5 w-3.5" />
          Edit
        </Button>
        <Button
          variant="ghost"
          size="sm"
          className={cn("h-7 text-xs gap-1.5", viewMode === "preview" && "bg-accent")}
          onClick={() => setViewMode("preview")}
        >
          <Eye className="h-3.5 w-3.5" />
          Preview
        </Button>
        <Button
          variant="ghost"
          size="sm"
          className={cn("h-7 text-xs gap-1.5", viewMode === "split" && "bg-accent")}
          onClick={() => setViewMode("split")}
        >
          <Columns3 className="h-3.5 w-3.5" />
          Split
        </Button>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {(viewMode === "edit" || viewMode === "split") && (
          <div className={cn("flex-1 overflow-auto border-r", viewMode === "split" && "w-1/2 flex-none")}>
            <KvasirEditor
              content={content}
              onChange={setContent}
              onSave={handleSave}
              readOnly={readOnly}
              showToolbar={viewMode === "edit"}
            />
          </div>
        )}
        {(viewMode === "preview" || viewMode === "split") && (
          <div className={cn("flex-1 overflow-auto", viewMode === "split" && "w-1/2 flex-none")}>
            <div
              ref={previewRef}
              className="prose prose-sm dark:prose-invert max-w-none px-6 py-4"
            />
          </div>
        )}
      </div>
    </div>
  )
}
