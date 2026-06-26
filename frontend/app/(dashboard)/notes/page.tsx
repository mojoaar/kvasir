"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { EditorPage } from "@/components/editor/editor-page"
import { TagChips } from "@/components/tags/tag-chips"
import { tagsApi } from "@/lib/api/tags"
import { useNotesStore } from "@/lib/store/note-store"

export default function NotesPage() {
  const { activeNoteId } = useNotesStore()
  const queryClient = useQueryClient()

  const { data: allTags = [] } = useQuery({
    queryKey: ["tags"],
    queryFn: tagsApi.list,
  })

  const { data: noteTags = [] } = useQuery({
    queryKey: ["noteTags", activeNoteId],
    queryFn: () =>
      activeNoteId ? tagsApi.getNoteTags(activeNoteId) : Promise.resolve([]),
    enabled: !!activeNoteId,
  })

  const addMutation = useMutation({
    mutationFn: async (tagId: number) => {
      if (!activeNoteId) return
      await tagsApi.addToNote(activeNoteId, tagId)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["noteTags", activeNoteId] })
    },
  })

  const removeMutation = useMutation({
    mutationFn: async (tagId: number) => {
      if (!activeNoteId) return
      await tagsApi.removeFromNote(activeNoteId, tagId)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["noteTags", activeNoteId] })
    },
  })

  const createMutation = useMutation({
    mutationFn: async ({
      name,
      color,
    }: {
      name: string
      color: string
    }) => {
      const tag = await tagsApi.create({ name, color })
      if (activeNoteId) {
        await tagsApi.addToNote(activeNoteId, tag.id)
      }
      return tag
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tags"] })
      queryClient.invalidateQueries({ queryKey: ["noteTags", activeNoteId] })
    },
  })

  return (
    <div className="flex h-full flex-col">
      {activeNoteId && (
        <div className="border-b px-6 py-2">
          <TagChips
            noteId={activeNoteId}
            tags={noteTags}
            allTags={allTags}
            onAdd={(tagId) => addMutation.mutate(tagId)}
            onRemove={(tagId) => removeMutation.mutate(tagId)}
            onCreate={(name, color) => createMutation.mutate({ name, color })}
          />
        </div>
      )}
      <div className="flex-1">
        <EditorPage />
      </div>
    </div>
  )
}
