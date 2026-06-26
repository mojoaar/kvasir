import { create } from "zustand"
import { persist } from "zustand/middleware"

export type ThemeName = "kvasir" | "dracula" | "nord" | "github" | "cyberpunk"
export type ThemeMode = "dark" | "light"

interface ThemeState {
  theme: ThemeName
  mode: ThemeMode
  setTheme: (theme: ThemeName) => void
  setMode: (mode: ThemeMode) => void
  toggleMode: () => void
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      theme: "kvasir",
      mode: "dark",
      setTheme: (theme) => {
        set({ theme })
        if (typeof document !== "undefined") {
          document.documentElement.setAttribute("data-theme", theme)
          localStorage.setItem("kvasir-theme", theme)
        }
      },
      setMode: (mode) => {
        set({ mode })
        if (typeof document !== "undefined") {
          document.documentElement.setAttribute("data-mode", mode)
          localStorage.setItem("kvasir-mode", mode)
        }
      },
      toggleMode: () => {
        set((state) => {
          const next = state.mode === "dark" ? "light" : "dark"
          if (typeof document !== "undefined") {
            document.documentElement.setAttribute("data-mode", next)
            localStorage.setItem("kvasir-mode", next)
          }
          return { mode: next }
        })
      },
    }),
    {
      name: "kvasir-theme-store",
      partialize: (state) => ({ theme: state.theme, mode: state.mode }),
    }
  )
)
