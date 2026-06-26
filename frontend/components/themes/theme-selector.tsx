"use client"

import { Paintbrush } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useThemeStore, type ThemeName } from "@/lib/store/theme-store"
import { themes } from "@/lib/themes"
import { cn } from "@/lib/utils"

const themeLabels: Record<ThemeName, string> = {
  kvasir: "Kvasir",
  dracula: "Dracula",
  nord: "Nord",
  github: "GitHub",
  cyberpunk: "Cyberpunk",
}

export function ThemeSelector() {
  const { theme, setTheme } = useThemeStore()

  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Button variant="ghost" size="icon" className="h-8 w-8">
          <Paintbrush className="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {(Object.keys(themes) as ThemeName[]).map((key) => (
          <DropdownMenuItem
            key={key}
            onClick={() => setTheme(key)}
            className={cn(theme === key && "bg-accent")}
          >
            <div className="flex items-center gap-2">
              <div
                className="h-3 w-3 rounded-full"
                style={{ backgroundColor: themes[key].dark.primary }}
              />
              {themeLabels[key]}
            </div>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
