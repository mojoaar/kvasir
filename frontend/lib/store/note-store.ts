import { create } from "zustand"

export interface Note {
  id: number
  title: string
  content: string
  vaultId?: number | null
  parentId?: number | null
  sortOrder: number
  createdAt: string
  updatedAt: string
  deletedAt?: string | null
}

interface NotesState {
  notes: Note[]
  activeNoteId: number | null
  isLoading: boolean
  setNotes: (notes: Note[]) => void
  setActiveNoteId: (id: number | null) => void
  setIsLoading: (loading: boolean) => void
  addNote: (note: Note) => void
  updateNote: (id: number, updates: Partial<Note>) => void
  removeNote: (id: number) => void
}

export const useNotesStore = create<NotesState>()((set) => ({
  notes: [],
  activeNoteId: null,
  isLoading: false,
  setNotes: (notes) => set({ notes }),
  setActiveNoteId: (id) => set({ activeNoteId: id }),
  setIsLoading: (loading) => set({ isLoading: loading }),
  addNote: (note) => set((state) => ({ notes: [note, ...state.notes] })),
  updateNote: (id, updates) =>
    set((state) => ({
      notes: state.notes.map((n) => (n.id === id ? { ...n, ...updates } : n)),
    })),
  removeNote: (id) =>
    set((state) => ({
      notes: state.notes.filter((n) => n.id !== id),
      activeNoteId: state.activeNoteId === id ? null : state.activeNoteId,
    })),
}))
