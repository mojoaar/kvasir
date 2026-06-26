import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor, render, cleanup } from "@testing-library/react"
import { useThemeStore } from "@/lib/store/theme-store"

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
  usePathname: () => "/",
  useSearchParams: () => new URLSearchParams(),
}))

import { CommandPalette } from "@/components/command-palette"

function openPalette() {
  document.dispatchEvent(new KeyboardEvent("keydown", { key: "k", metaKey: true }))
}

describe("CommandPalette", () => {
  beforeEach(() => {
    cleanup()
    vi.spyOn(console, "error").mockImplementation(() => {})
    useThemeStore.setState({ theme: "kvasir", mode: "dark" })
  })

  it("renders when Cmd+K is pressed", async () => {
    render(<CommandPalette />)
    openPalette()

    await waitFor(() => {
      expect(screen.getByPlaceholderText("Search notes, create, switch theme...")).toBeInTheDocument()
    })
  })

  it("renders when Ctrl+K is pressed", async () => {
    render(<CommandPalette />)
    document.dispatchEvent(new KeyboardEvent("keydown", { key: "k", ctrlKey: true }))

    await waitFor(() => {
      expect(screen.getByPlaceholderText("Search notes, create, switch theme...")).toBeInTheDocument()
    })
  })

  it("shows navigation commands", async () => {
    render(<CommandPalette />)
    openPalette()

    await waitFor(() => {
      const notes = screen.getAllByText("Notes")
      expect(notes.length).toBeGreaterThan(0)
      expect(screen.getByText("Search")).toBeInTheDocument()
    })
  })

  it("shows theme commands", async () => {
    render(<CommandPalette />)
    openPalette()

    await waitFor(() => {
      expect(screen.getByText("Kvasir (Nordic)")).toBeInTheDocument()
    })
  })

  it("does not open on other key combinations", () => {
    render(<CommandPalette />)

    document.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" }))
    document.dispatchEvent(new KeyboardEvent("keydown", { key: "k", altKey: true }))

    expect(screen.queryByPlaceholderText("Search notes, create, switch theme...")).not.toBeInTheDocument()
  })
})
