"use client"

import { useState, useCallback, useMemo } from "react"
import { useNotesStore, type Note } from "@/lib/store/note-store"
import { cn } from "@/lib/utils"
import {
  ChevronRight,
  FileText,
  Folder,
  FolderOpen,
  Plus,
  MoreHorizontal,
  Pencil,
  Trash2,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"

interface TreeNode {
  note: Note
  children: TreeNode[]
}

function buildTree(notes: Note[]): TreeNode[] {
  const map = new Map<number, TreeNode>()
  const roots: TreeNode[] = []

  for (const note of notes) {
    map.set(note.id, { note, children: [] })
  }

  for (const note of notes) {
    const node = map.get(note.id)!
    if (note.parentId != null && map.has(note.parentId)) {
      map.get(note.parentId)!.children.push(node)
    } else {
      roots.push(node)
    }
  }

  const sorter = (a: TreeNode, b: TreeNode) => a.note.sortOrder - b.note.sortOrder
  for (const node of map.values()) {
    node.children.sort(sorter)
  }
  roots.sort(sorter)

  return roots
}

function countNotes(node: TreeNode): number {
  let count = node.note.isFolder ? 0 : 1
  for (const child of node.children) {
    count += countNotes(child)
  }
  return count
}

export function NoteTree() {
  const {
    notes,
    activeNoteId,
    setActiveNoteId,
    addNote,
    updateNote,
    removeNote,
    moveNote,
  } = useNotesStore()

  const [editingId, setEditingId] = useState<number | null>(null)
  const [editValue, setEditValue] = useState("")
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set())
  const [creatingIn, setCreatingIn] = useState<number | null>(null)
  const [creatingType, setCreatingType] = useState<"note" | "folder">("note")
  const [createValue, setCreateValue] = useState("")
  const [dragOverId, setDragOverId] = useState<number | null>(null)

  const tree = useMemo(() => buildTree(notes), [notes])

  const toggleExpand = useCallback((id: number) => {
    setExpandedIds((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }, [])

  const startRename = useCallback((note: Note) => {
    setEditingId(note.id)
    setEditValue(note.title)
  }, [])

  const commitRename = useCallback(() => {
    if (editingId != null && editValue.trim()) {
      updateNote(editingId, { title: editValue.trim() })
    }
    setEditingId(null)
    setEditValue("")
  }, [editingId, editValue, updateNote])

  const startCreate = useCallback((parentId: number | null, type: "note" | "folder") => {
    setCreatingIn(parentId)
    setCreatingType(type)
    setCreateValue("")
  }, [])

  const commitCreate = useCallback(() => {
    if (createValue.trim()) {
      const now = new Date().toISOString()
      addNote({
        id: Date.now(),
        title: createValue.trim(),
        content: "",
        parentId: creatingIn,
        isFolder: creatingType === "folder",
        sortOrder: 0,
        createdAt: now,
        updatedAt: now,
      })
      if (creatingIn != null) {
        setExpandedIds((prev) => new Set(prev).add(creatingIn))
      }
    }
    setCreatingIn(null)
    setCreateValue("")
  }, [createValue, creatingIn, creatingType, addNote])

  const handleDelete = useCallback((id: number) => {
    removeNote(id)
  }, [removeNote])

  const handleDragStart = useCallback((e: React.DragEvent, id: number) => {
    e.dataTransfer.setData("application/kvasir-note", id.toString())
    e.dataTransfer.effectAllowed = "move"
  }, [])

  const handleDragOver = useCallback((e: React.DragEvent, id: number | null) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = "move"
    setDragOverId(id)
  }, [])

  const handleDragLeave = useCallback(() => {
    setDragOverId(null)
  }, [])

  const handleDrop = useCallback(
    (e: React.DragEvent, targetId: number | null) => {
      e.preventDefault()
      setDragOverId(null)
      const draggedId = parseInt(e.dataTransfer.getData("application/kvasir-note"))
      if (draggedId === targetId) return
      if (targetId != null && !notes.find((n) => n.id === targetId)?.isFolder) return

      moveNote(draggedId, targetId, 0)
      if (targetId != null) {
        setExpandedIds((prev) => new Set(prev).add(targetId))
      }
    },
    [notes, moveNote]
  )

  const renderNode = (node: TreeNode, depth: number) => {
    const { note } = node
    const isExpanded = expandedIds.has(note.id)
    const isActive = activeNoteId === note.id
    const isEditing = editingId === note.id
    const noteCount = note.isFolder ? countNotes(node) : 0
    const isDragOver = dragOverId === note.id

    return (
      <div key={note.id}>
        <div
          className={cn(
            "group flex items-center gap-1 px-1 py-0.5 text-sm rounded-md cursor-pointer transition-colors",
            isActive && "bg-accent text-accent-foreground",
            !isActive && "hover:bg-accent/50",
            isDragOver && "ring-2 ring-primary/50 bg-accent/30"
          )}
          style={{ paddingLeft: `${8 + depth * 16}px` }}
          draggable={!isEditing}
          onDragStart={(e) => handleDragStart(e, note.id)}
          onDragOver={(e) => note.isFolder && handleDragOver(e, note.id)}
          onDragLeave={handleDragLeave}
          onDrop={(e) => note.isFolder && handleDrop(e, note.id)}
          onClick={() => {
            if (note.isFolder) {
              toggleExpand(note.id)
            } else {
              setActiveNoteId(note.id)
            }
          }}
          onDoubleClick={() => {
            if (!note.isFolder) return
            startRename(note)
          }}
        >
          {note.isFolder && (
            <ChevronRight
              className={cn(
                "h-3.5 w-3.5 shrink-0 transition-transform text-muted-foreground",
                isExpanded && "rotate-90"
              )}
            />
          )}
          {!note.isFolder && <div className="w-3.5 shrink-0" />}

          {note.isFolder ? (
            isExpanded ? (
              <FolderOpen className="h-4 w-4 shrink-0 text-muted-foreground" />
            ) : (
              <Folder className="h-4 w-4 shrink-0 text-muted-foreground" />
            )
          ) : (
            <FileText className="h-4 w-4 shrink-0 text-muted-foreground" />
          )}

          {isEditing ? (
            <Input
              value={editValue}
              onChange={(e) => setEditValue(e.target.value)}
              onBlur={commitRename}
              onKeyDown={(e) => {
                if (e.key === "Enter") commitRename()
                if (e.key === "Escape") {
                  setEditingId(null)
                  setEditValue("")
                }
              }}
              className="h-5 flex-1 text-xs px-1 py-0"
              autoFocus
              onClick={(e) => e.stopPropagation()}
            />
          ) : (
            <span className="flex-1 truncate">{note.title}</span>
          )}

          {note.isFolder && noteCount > 0 && (
            <span className="text-[10px] text-muted-foreground tabular-nums">
              {noteCount}
            </span>
          )}

          <div className="flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
            <Button
              variant="ghost"
              size="icon"
              className="h-5 w-5"
              onClick={(e) => {
                e.stopPropagation()
                startCreate(note.id, "note")
              }}
            >
              <Plus className="h-3 w-3" />
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-5 w-5"
                  onClick={(e) => e.stopPropagation()}
                >
                  <MoreHorizontal className="h-3 w-3" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="start" className="w-36">
                <DropdownMenuItem
                  onClick={(e) => {
                    e.stopPropagation()
                    startRename(note)
                  }}
                >
                  <Pencil className="h-3.5 w-3.5 mr-2" />
                  Rename
                </DropdownMenuItem>
                {note.isFolder && (
                  <DropdownMenuItem
                    onClick={(e) => {
                      e.stopPropagation()
                      startCreate(note.id, "folder")
                    }}
                  >
                    <Folder className="h-3.5 w-3.5 mr-2" />
                    New folder
                  </DropdownMenuItem>
                )}
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={(e) => {
                    e.stopPropagation()
                    handleDelete(note.id)
                  }}
                >
                  <Trash2 className="h-3.5 w-3.5 mr-2" />
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>

        {note.isFolder && isExpanded && node.children.length > 0 && (
          <div>
            {node.children.map((child) => renderNode(child, depth + 1))}
          </div>
        )}

        {note.isFolder && isExpanded && node.children.length === 0 && (
          <div
            className="text-[11px] text-muted-foreground/60 italic py-1"
            style={{ paddingLeft: `${16 + (depth + 1) * 16}px` }}
          >
            Empty folder
          </div>
        )}

        {creatingIn === note.id && (
          <div
            className="flex items-center gap-1 px-1 py-0.5"
            style={{ paddingLeft: `${8 + (depth + 1) * 16}px` }}
          >
            {createValue === "" ? (
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 text-xs gap-1"
                  onClick={() => {
                    setCreateValue(creatingType === "note" ? "Untitled" : "New Folder")
                  }}
                >
                  <FileText className="h-3 w-3" />
                  Note
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 text-xs gap-1"
                  onClick={() => {
                    setCreateValue("New Folder")
                  }}
                >
                  <Folder className="h-3 w-3" />
                  Folder
                </Button>
              </div>
            ) : (
              <Input
                value={createValue}
                onChange={(e) => setCreateValue(e.target.value)}
                onBlur={commitCreate}
                onKeyDown={(e) => {
                  if (e.key === "Enter") commitCreate()
                  if (e.key === "Escape") {
                    setCreatingIn(null)
                    setCreateValue("")
                  }
                }}
                className="h-5 flex-1 text-xs px-1 py-0"
                autoFocus
              />
            )}
          </div>
        )}
      </div>
    )
  }

  return (
    <div
      className="flex flex-col h-full"
      onDragOver={(e) => handleDragOver(e, null)}
      onDragLeave={handleDragLeave}
      onDrop={(e) => handleDrop(e, null)}
    >
      <div className="flex items-center justify-between px-2 py-1.5">
        <span className="text-xs font-semibold text-muted-foreground tracking-wide uppercase">
          Notes
        </span>
        <div className="flex items-center gap-0.5">
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6"
            title="New note"
            onClick={() => startCreate(null, "note")}
          >
            <FileText className="h-3.5 w-3.5" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6"
            title="New folder"
            onClick={() => startCreate(null, "folder")}
          >
            <Folder className="h-3.5 w-3.5" />
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto px-1">
        {tree.length === 0 && creatingIn == null ? (
          <p className="text-xs text-muted-foreground px-3 py-4 text-center">
            No notes yet. Create one to get started.
          </p>
        ) : (
          tree.map((node) => renderNode(node, 0))
        )}
      </div>

      {creatingIn == null && (
        <div className="px-2 py-2">
          {createValue === "" ? (
            <div className="flex items-center gap-1">
              <Button
                variant="outline"
                size="sm"
                className="flex-1 h-7 text-xs gap-1"
                onClick={() => {
                  setCreatingType("note")
                  setCreateValue("Untitled")
                }}
              >
                <FileText className="h-3 w-3" />
                New Note
              </Button>
              <Button
                variant="outline"
                size="sm"
                className="flex-1 h-7 text-xs gap-1"
                onClick={() => {
                  setCreatingType("folder")
                  setCreateValue("New Folder")
                }}
              >
                <Folder className="h-3 w-3" />
                New Folder
              </Button>
            </div>
          ) : (
            <Input
              value={createValue}
              onChange={(e) => setCreateValue(e.target.value)}
              onBlur={commitCreate}
              onKeyDown={(e) => {
                if (e.key === "Enter") commitCreate()
                if (e.key === "Escape") {
                  setCreatingIn(null)
                  setCreateValue("")
                }
              }}
              placeholder={creatingType === "note" ? "Note title..." : "Folder name..."}
              className="h-7 text-xs"
              autoFocus
            />
          )}
        </div>
      )}
    </div>
  )
}
