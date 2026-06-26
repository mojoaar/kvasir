import { z } from "zod"

export const noteSchema = z.object({
  id: z.number(),
  title: z.string(),
  content: z.string(),
  vaultId: z.number().nullable().optional(),
  parentId: z.number().nullable().optional(),
  isFolder: z.boolean(),
  sortOrder: z.number(),
  createdAt: z.string(),
  updatedAt: z.string(),
  deletedAt: z.string().nullable().optional(),
})

export type Note = z.infer<typeof noteSchema>

export const createNoteSchema = z.object({
  title: z.string().min(1, "Title is required"),
  content: z.string().default(""),
  vaultId: z.number().nullable().optional(),
  parentId: z.number().nullable().optional(),
  isFolder: z.boolean().default(false),
  sortOrder: z.number().default(0),
})

export type CreateNote = z.infer<typeof createNoteSchema>

export const updateNoteSchema = z.object({
  title: z.string().min(1, "Title is required"),
  content: z.string().default(""),
  vaultId: z.number().nullable().optional(),
  parentId: z.number().nullable().optional(),
  isFolder: z.boolean().default(false),
  sortOrder: z.number().default(0),
})

export type UpdateNote = z.infer<typeof updateNoteSchema>

export const notesQuerySchema = z.object({
  offset: z.number().int().min(0).default(0),
  limit: z.number().int().min(1).max(200).default(50),
  vaultId: z.number().int().optional(),
  parentId: z.number().int().optional(),
})

export type NotesQuery = z.infer<typeof notesQuerySchema>

export const tagSchema = z.object({
  id: z.number(),
  name: z.string(),
  color: z.string(),
  createdAt: z.string(),
})

export type Tag = z.infer<typeof tagSchema>

export const createTagSchema = z.object({
  name: z.string().min(1, "Name is required"),
  color: z.string().min(1, "Color is required"),
})

export type CreateTag = z.infer<typeof createTagSchema>

export const addTagRequestSchema = z.object({
  tagId: z.number(),
})
