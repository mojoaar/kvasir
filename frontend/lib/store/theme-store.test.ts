import { describe, it, expect, beforeEach } from "vitest"
import { useThemeStore } from "@/lib/store/theme-store"

describe("useThemeStore", () => {
  beforeEach(() => {
    useThemeStore.setState({ theme: "kvasir", mode: "dark" })
  })

  it("has kvasir dark as default", () => {
    const state = useThemeStore.getState()
    expect(state.theme).toBe("kvasir")
    expect(state.mode).toBe("dark")
  })

  it("setTheme changes theme and persists", () => {
    useThemeStore.getState().setTheme("dracula")
    expect(useThemeStore.getState().theme).toBe("dracula")
    expect(localStorage.getItem("kvasir-theme")).toBe("dracula")
  })

  it("setTheme accepts all 5 theme names", () => {
    const themes = ["kvasir", "dracula", "nord", "github", "cyberpunk"] as const
    for (const theme of themes) {
      useThemeStore.getState().setTheme(theme)
      expect(useThemeStore.getState().theme).toBe(theme)
    }
  })

  it("setMode changes mode", () => {
    useThemeStore.getState().setMode("light")
    expect(useThemeStore.getState().mode).toBe("light")
    expect(localStorage.getItem("kvasir-mode")).toBe("light")
  })

  it("toggleMode switches dark to light", () => {
    useThemeStore.getState().toggleMode()
    expect(useThemeStore.getState().mode).toBe("light")
    expect(localStorage.getItem("kvasir-mode")).toBe("light")
  })

  it("toggleMode switches light to dark", () => {
    useThemeStore.getState().setMode("light")
    useThemeStore.getState().toggleMode()
    expect(useThemeStore.getState().mode).toBe("dark")
  })

  it("persists theme and mode to localStorage", () => {
    useThemeStore.getState().setTheme("nord")
    useThemeStore.getState().setMode("light")
    const raw = localStorage.getItem("kvasir-theme-store")
    const parsed = JSON.parse(raw!)
    expect(parsed.state.theme).toBe("nord")
    expect(parsed.state.mode).toBe("light")
  })
})
