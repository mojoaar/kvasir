import { api } from "./client"
import type { Tag, CreateTag } from "./schemas"

export const tagsApi = {
  list: () => api.get<Tag[]>("/api/v1/tags"),
  create: (body: CreateTag) => api.post<Tag>("/api/v1/tags", body),
  get: (id: number) => api.get<Tag>(`/api/v1/tags/${id}`),
  update: (id: number, body: CreateTag) =>
    api.put<Tag>(`/api/v1/tags/${id}`, body),
  delete: (id: number) => api.delete<{ status: string }>(`/api/v1/tags/${id}`),
  getNoteTags: (noteId: number) => api.get<Tag[]>(`/api/v1/notes/${noteId}/tags`),
  addToNote: (noteId: number, tagId: number) =>
    api.post<{ status: string }>(`/api/v1/notes/${noteId}/tags`, { tagId }),
  removeFromNote: (noteId: number, tagId: number) =>
    api.delete<{ status: string }>(`/api/v1/notes/${noteId}/tags`, { tagId }),
}
