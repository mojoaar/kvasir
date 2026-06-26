import { describe, it, expect, beforeEach } from "vitest"
import { useNotesStore } from "@/lib/store/note-store"

const testNote = {
  id: 1,
  title: "Test",
  content: "# Hello",
  vaultId: null as number | null,
  parentId: null as number | null,
  isFolder: false,
  sortOrder: 0,
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-01-01T00:00:00Z",
  deletedAt: null as string | null,
}

describe("useNotesStore", () => {
  beforeEach(() => {
    useNotesStore.setState({ notes: [], activeNoteId: null, isLoading: false })
  })

  it("has empty initial state", () => {
    const state = useNotesStore.getState()
    expect(state.notes).toEqual([])
    expect(state.activeNoteId).toBeNull()
    expect(state.isLoading).toBe(false)
  })

  it("setNotes replaces all notes", () => {
    useNotesStore.getState().setNotes([testNote])
    expect(useNotesStore.getState().notes).toHaveLength(1)
  })

  it("addNote prepends to list", () => {
    useNotesStore.getState().addNote(testNote)
    const second = { ...testNote, id: 2, title: "Second" }
    useNotesStore.getState().addNote(second)
    expect(useNotesStore.getState().notes).toHaveLength(2)
    expect(useNotesStore.getState().notes[0].id).toBe(2)
  })

  it("updateNote updates matching note", () => {
    useNotesStore.getState().addNote(testNote)
    useNotesStore.getState().updateNote(1, { title: "Updated" })
    expect(useNotesStore.getState().notes[0].title).toBe("Updated")
  })

  it("updateNote does nothing for non-existent id", () => {
    useNotesStore.getState().addNote(testNote)
    useNotesStore.getState().updateNote(999, { title: "Nope" })
    expect(useNotesStore.getState().notes[0].title).toBe("Test")
  })

  it("removeNote removes matching note", () => {
    useNotesStore.getState().addNote(testNote)
    useNotesStore.getState().removeNote(1)
    expect(useNotesStore.getState().notes).toHaveLength(0)
  })

  it("removeNote clears activeNoteId if it was the removed note", () => {
    useNotesStore.getState().addNote(testNote)
    useNotesStore.getState().setActiveNoteId(1)
    useNotesStore.getState().removeNote(1)
    expect(useNotesStore.getState().activeNoteId).toBeNull()
  })

  it("removeNote keeps activeNoteId if it was a different note", () => {
    useNotesStore.getState().addNote(testNote)
    useNotesStore.getState().addNote({ ...testNote, id: 2 })
    useNotesStore.getState().setActiveNoteId(2)
    useNotesStore.getState().removeNote(1)
    expect(useNotesStore.getState().activeNoteId).toBe(2)
  })

  it("setActiveNoteId updates active note", () => {
    useNotesStore.getState().setActiveNoteId(42)
    expect(useNotesStore.getState().activeNoteId).toBe(42)
    useNotesStore.getState().setActiveNoteId(null)
    expect(useNotesStore.getState().activeNoteId).toBeNull()
  })

  it("setIsLoading toggles loading state", () => {
    useNotesStore.getState().setIsLoading(true)
    expect(useNotesStore.getState().isLoading).toBe(true)
    useNotesStore.getState().setIsLoading(false)
    expect(useNotesStore.getState().isLoading).toBe(false)
  })

  it("moveNote updates parentId and sortOrder", () => {
    useNotesStore.getState().addNote(testNote)
    useNotesStore.getState().moveNote(1, 5, 3)
    const note = useNotesStore.getState().notes[0]
    expect(note.parentId).toBe(5)
    expect(note.sortOrder).toBe(3)
  })

  it("handles multiple operations sequentially", () => {
    useNotesStore.getState().addNote(testNote)
    useNotesStore.getState().setActiveNoteId(1)
    useNotesStore.getState().updateNote(1, { title: "Changed" })
    useNotesStore.getState().moveNote(1, 10, 1)
    useNotesStore.getState().setIsLoading(true)

    const state = useNotesStore.getState()
    expect(state.activeNoteId).toBe(1)
    expect(state.notes[0].title).toBe("Changed")
    expect(state.notes[0].parentId).toBe(10)
    expect(state.isLoading).toBe(true)
  })
})
