import { describe, it, expect, beforeEach } from "vitest"
import { useUIStore } from "@/lib/store/ui-store"

describe("useUIStore", () => {
  beforeEach(() => {
    useUIStore.setState({ sidebarOpen: true, commandPaletteOpen: false })
  })

  it("has default state", () => {
    const state = useUIStore.getState()
    expect(state.sidebarOpen).toBe(true)
    expect(state.commandPaletteOpen).toBe(false)
  })

  it("toggleSidebar switches open/closed", () => {
    useUIStore.getState().toggleSidebar()
    expect(useUIStore.getState().sidebarOpen).toBe(false)
    useUIStore.getState().toggleSidebar()
    expect(useUIStore.getState().sidebarOpen).toBe(true)
  })

  it("setSidebarOpen sets explicitly", () => {
    useUIStore.getState().setSidebarOpen(false)
    expect(useUIStore.getState().sidebarOpen).toBe(false)
    useUIStore.getState().setSidebarOpen(true)
    expect(useUIStore.getState().sidebarOpen).toBe(true)
  })

  it("toggleCommandPalette switches open/closed", () => {
    useUIStore.getState().toggleCommandPalette()
    expect(useUIStore.getState().commandPaletteOpen).toBe(true)
    useUIStore.getState().toggleCommandPalette()
    expect(useUIStore.getState().commandPaletteOpen).toBe(false)
  })

  it("setCommandPaletteOpen sets explicitly", () => {
    useUIStore.getState().setCommandPaletteOpen(true)
    expect(useUIStore.getState().commandPaletteOpen).toBe(true)
    useUIStore.getState().setCommandPaletteOpen(false)
    expect(useUIStore.getState().commandPaletteOpen).toBe(false)
  })
})
