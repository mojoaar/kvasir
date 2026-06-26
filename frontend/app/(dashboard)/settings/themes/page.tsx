"use client"

import { useThemeStore, type ThemeName } from "@/lib/store/theme-store"
import { themes } from "@/lib/themes"
import { cn } from "@/lib/utils"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

const themeLabels: Record<ThemeName, string> = {
  kvasir: "Kvasir",
  dracula: "Dracula",
  nord: "Nord",
  github: "GitHub",
  cyberpunk: "Cyberpunk",
}

const themeDescriptions: Record<ThemeName, string> = {
  kvasir: "Nordic dark elegance — the default Kvasir look.",
  dracula: "Purple-toned dark theme inspired by the Dracula editor theme.",
  nord: "Arctic, blue-tinted palette based on the Nord design system.",
  github: "Clean, familiar look matching the GitHub UI.",
  cyberpunk: "High-contrast neon on black for the bold.",
}

export default function ThemesPage() {
  const { theme, setTheme } = useThemeStore()

  return (
    <div className="p-6 space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Themes</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Choose your color scheme. Each theme supports dark and light mode.
        </p>
      </div>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {(Object.keys(themes) as ThemeName[]).map((key) => (
          <Card
            key={key}
            className={cn(
              "cursor-pointer transition-colors hover:border-primary/50",
              theme === key && "border-primary ring-2 ring-primary/30"
            )}
            onClick={() => setTheme(key)}
          >
            <CardHeader className="pb-3">
              <div className="flex items-center gap-2">
                <div
                  className="h-4 w-4 rounded-full"
                  style={{ backgroundColor: themes[key].dark.primary }}
                />
                <CardTitle className="text-base">{themeLabels[key]}</CardTitle>
              </div>
              <CardDescription>{themeDescriptions[key]}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex gap-2">
                <div className="flex-1 h-8 rounded border" style={{
                  background: themes[key].dark.background,
                  borderColor: themes[key].dark.border,
                }}>
                  <div className="flex items-center gap-1 p-1.5">
                    <div className="h-3 w-3 rounded-full" style={{ background: themes[key].dark.primary }} />
                    <div className="h-2 flex-1 rounded-sm" style={{ background: themes[key].dark.muted }} />
                  </div>
                </div>
                <div className="flex-1 h-8 rounded border" style={{
                  background: themes[key].light.background,
                  borderColor: themes[key].light.border,
                }}>
                  <div className="flex items-center gap-1 p-1.5">
                    <div className="h-3 w-3 rounded-full" style={{ background: themes[key].light.primary }} />
                    <div className="h-2 flex-1 rounded-sm" style={{ background: themes[key].light.muted }} />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
