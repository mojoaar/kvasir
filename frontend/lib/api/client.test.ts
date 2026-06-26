import { describe, it, expect, vi, beforeEach, afterEach } from "vitest"
import { api, healthCheck } from "@/lib/api/client"

describe("api client", () => {
  let originalFetch: typeof globalThis.fetch

  beforeEach(() => {
    originalFetch = globalThis.fetch
    vi.restoreAllMocks()
  })

  it("healthCheck returns status ok", async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ status: "ok", version: "0.1.0" }),
    }) as unknown as typeof fetch

    const result = await healthCheck()
    expect(result.status).toBe("ok")
    expect(result.version).toBe("0.1.0")
  })

  it("get makes a GET request", async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve([{ id: 1 }]),
    }) as unknown as typeof fetch

    const result = await api.get("/api/v1/notes")
    expect(result).toEqual([{ id: 1 }])
    expect(globalThis.fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/v1/notes"),
      expect.objectContaining({ headers: expect.any(Object) })
    )
  })

  it("post sends JSON body", async () => {
    const mockNote = { id: 1, title: "Test" }
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockNote),
    }) as unknown as typeof fetch

    const result = await api.post("/api/v1/notes", { title: "Test" })
    expect(result).toEqual(mockNote)
    expect(globalThis.fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({ method: "POST" })
    )
  })

  it("put sends JSON body", async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ id: 1, title: "Updated" }),
    }) as unknown as typeof fetch

    const result = await api.put("/api/v1/notes/1", { title: "Updated" })
    expect(result.title).toBe("Updated")
    expect(globalThis.fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({ method: "PUT" })
    )
  })

  it("delete sends DELETE request", async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ success: true }),
    }) as unknown as typeof fetch

    await api.delete("/api/v1/notes/1")
    expect(globalThis.fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({ method: "DELETE" })
    )
  })

  it("handles HTTP errors with error message", async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 404,
      statusText: "Not Found",
      json: () => Promise.resolve({ error: "note not found" }),
    }) as unknown as typeof fetch

    await expect(api.get("/api/v1/notes/999")).rejects.toThrow("note not found")
  })

  afterEach(() => {
    globalThis.fetch = originalFetch
  })
})
