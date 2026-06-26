"use client"

import { useState } from "react"
import { X, Plus } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import type { Tag } from "@/lib/api/schemas"
import { getTagColor, getTagContrastColor } from "@/lib/tag-colors"

interface TagChipsProps {
  noteId?: number
  tags: Tag[]
  allTags: Tag[]
  onAdd: (tagId: number) => void
  onRemove: (tagId: number) => void
  onCreate: (name: string, color: string) => void
}

export function TagChips({
  tags,
  allTags,
  onAdd,
  onRemove,
  onCreate,
}: TagChipsProps) {
  const [showDropdown, setShowDropdown] = useState(false)
  const [newTagName, setNewTagName] = useState("")
  const [isCreating, setIsCreating] = useState(false)

  const availableTags = allTags.filter(
    (t) => !tags.some((nt) => nt.id === t.id),
  )

  function handleCreate() {
    const name = newTagName.trim()
    if (!name) return
    const colorIndex = allTags.length
    onCreate(name, getTagColor(colorIndex))
    setNewTagName("")
    setIsCreating(false)
  }

  return (
    <div className="flex flex-wrap items-center gap-1.5">
      {tags.map((tag) => (
        <span
          key={tag.id}
          className="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium"
          style={{
            backgroundColor: tag.color,
            color: getTagContrastColor(tag.color),
          }}
        >
          {tag.name}
          <button
            type="button"
            onClick={() => onRemove(tag.id)}
            className="ml-0.5 rounded-full p-0.5 hover:bg-black/10"
          >
            <X className="h-3 w-3" />
          </button>
        </span>
      ))}
      <div className="relative">
        <button
          type="button"
          onClick={() => setShowDropdown(!showDropdown)}
          className="inline-flex items-center gap-1 rounded-full border border-dashed px-2.5 py-0.5 text-xs text-muted-foreground hover:border-foreground hover:text-foreground"
        >
          <Plus className="h-3 w-3" />
          Add tag
        </button>
        {showDropdown && (
          <div className="absolute left-0 top-full z-50 mt-1 w-52 rounded-md border bg-popover p-2 shadow-md">
            {availableTags.length > 0 && (
              <div className="mb-1 max-h-32 overflow-auto">
                {availableTags.map((tag) => (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() => {
                      onAdd(tag.id)
                      setShowDropdown(false)
                    }}
                    className="flex w-full items-center gap-2 rounded px-2 py-1 text-xs hover:bg-muted"
                  >
                    <span
                      className="h-2.5 w-2.5 rounded-full"
                      style={{ backgroundColor: tag.color }}
                    />
                    {tag.name}
                  </button>
                ))}
              </div>
            )}
            {isCreating ? (
              <div className="flex gap-1">
                <Input
                  autoFocus
                  value={newTagName}
                  onChange={(e) => setNewTagName(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") handleCreate()
                    if (e.key === "Escape") {
                      setIsCreating(false)
                      setNewTagName("")
                    }
                  }}
                  placeholder="New tag..."
                  className="h-7 text-xs"
                />
                <Button
                  type="button"
                  size="sm"
                  className="h-7 text-xs"
                  onClick={handleCreate}
                >
                  Add
                </Button>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => setIsCreating(true)}
                className="flex w-full items-center gap-2 rounded px-2 py-1 text-xs text-muted-foreground hover:bg-muted"
              >
                <Plus className="h-3 w-3" />
                Create new tag
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
