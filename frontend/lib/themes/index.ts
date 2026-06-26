import { kvasir } from "./kvasir"
import { dracula } from "./dracula"
import { nord } from "./nord"
import { github } from "./github"
import { cyberpunk } from "./cyberpunk"
import type { ThemeName } from "@/lib/store/theme-store"

export type { ThemeName } from "@/lib/store/theme-store"

export interface ThemeColors {
  readonly background: string
  readonly foreground: string
  readonly card: string
  readonly cardForeground: string
  readonly popover: string
  readonly popoverForeground: string
  readonly primary: string
  readonly primaryForeground: string
  readonly secondary: string
  readonly secondaryForeground: string
  readonly muted: string
  readonly mutedForeground: string
  readonly accent: string
  readonly accentForeground: string
  readonly destructive: string
  readonly border: string
  readonly input: string
  readonly ring: string
  readonly chart1: string
  readonly chart2: string
  readonly chart3: string
  readonly chart4: string
  readonly chart5: string
  readonly sidebar: string
  readonly sidebarForeground: string
  readonly sidebarPrimary: string
  readonly sidebarPrimaryForeground: string
  readonly sidebarAccent: string
  readonly sidebarAccentForeground: string
  readonly sidebarBorder: string
  readonly sidebarRing: string
  readonly surface: string
  readonly surfaceHover: string
  readonly success: string
  readonly warning: string
  readonly error: string
  readonly link: string
  readonly codeBlock: string
  readonly highlight: string
  readonly accentSecondary: string
  readonly accentTertiary: string
  readonly editorBackground: string
  readonly editorText: string
}

export interface ThemeDefinition {
  readonly name: string
  readonly dark: ThemeColors
  readonly light: ThemeColors
}

export const themes: Record<ThemeName, ThemeDefinition> = {
  kvasir,
  dracula,
  nord,
  github,
  cyberpunk,
}

export const themeNames = Object.keys(themes) as ThemeName[]

export function getTheme(name: ThemeName) {
  return themes[name]
}

export function getThemeColors(name: ThemeName, mode: "dark" | "light") {
  return themes[name][mode]
}
