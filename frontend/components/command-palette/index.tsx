"use client"

import { useEffect } from "react"
import { Command } from "cmdk"
import { useRouter } from "next/navigation"
import {
  Search,
  FilePlus,
  FolderPlus,
  Settings,
  Sun,
  Moon,
  Keyboard,
  StickyNote,
} from "lucide-react"
import { useUIStore } from "@/lib/store/ui-store"
import { useThemeStore, type ThemeName } from "@/lib/store/theme-store"

const themes: { name: ThemeName; label: string }[] = [
  { name: "kvasir", label: "Kvasir (Nordic)" },
  { name: "dracula", label: "Dracula" },
  { name: "nord", label: "Nord" },
  { name: "github", label: "GitHub" },
  { name: "cyberpunk", label: "Cyberpunk" },
]

export function CommandPalette() {
  const router = useRouter()
  const open = useUIStore((s) => s.commandPaletteOpen)
  const setOpen = useUIStore((s) => s.setCommandPaletteOpen)
  const { theme, mode, setTheme, toggleMode } = useThemeStore()

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        setOpen(!open)
      }
    }
    document.addEventListener("keydown", down)
    return () => document.removeEventListener("keydown", down)
  }, [open, setOpen])

  return (
    <Command.Dialog
      open={open}
      onOpenChange={setOpen}
      label="Command Palette"
      className="[&_[cmdk-group-heading]]:text-muted-foreground [&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:py-1.5 [&_[cmdk-group-heading]]:text-xs [&_[cmdk-group-heading]]:font-semibold [&_[cmdk-group-heading]]:tracking-wider [&_[cmdk-item]]:px-2 [&_[cmdk-item]]:py-3 [&_[cmdk-item]_svg]:size-4 [&_[cmdk-item]_svg]:mr-2 [&_[cmdk-input]]:h-12"
    >
      <div className="flex items-center border-b px-3">
        <Search className="mr-2 size-4 shrink-0 opacity-50" />
        <Command.Input
          placeholder="Search notes, create, switch theme..."
          className="flex h-12 w-full bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground"
        />
      </div>
      <Command.List className="max-h-[300px] overflow-y-auto overflow-x-hidden p-2">
        <Command.Empty className="py-6 text-center text-sm text-muted-foreground">
          No results found.
        </Command.Empty>

        <Command.Group heading="Navigation">
          <Command.Item
            value="notes"
            onSelect={() => {
              router.push("/notes")
              setOpen(false)
            }}
          >
            <StickyNote className="size-4" />
            Notes
          </Command.Item>
          <Command.Item
            value="search"
            onSelect={() => {
              router.push("/search")
              setOpen(false)
            }}
          >
            <Search className="size-4" />
            Search
          </Command.Item>
          <Command.Item
            value="settings"
            onSelect={() => {
              router.push("/settings")
              setOpen(false)
            }}
          >
            <Settings className="size-4" />
            Settings
          </Command.Item>
        </Command.Group>

        <Command.Group heading="Actions">
          <Command.Item
            value="create note"
            onSelect={() => {
              router.push("/notes?new=true")
              setOpen(false)
            }}
          >
            <FilePlus className="size-4" />
            Create Note
          </Command.Item>
          <Command.Item
            value="create folder"
            onSelect={() => {
              router.push("/notes?new=folder")
              setOpen(false)
            }}
          >
            <FolderPlus className="size-4" />
            Create Folder
          </Command.Item>
        </Command.Group>

        <Command.Group heading="Theme">
          {themes.map((t) => (
            <Command.Item
              key={t.name}
              value={`theme ${t.label}`}
              onSelect={() => {
                setTheme(t.name)
                setOpen(false)
              }}
            >
              <span
                className="mr-2 inline-block size-3 rounded-full ring-1 ring-border"
                style={{
                  background:
                    t.name === theme
                      ? mode === "dark"
                        ? "var(--color-primary)"
                        : "var(--color-primary)"
                      : "transparent",
                }}
              />
              {t.label}
              {t.name === theme && (
                <span className="ml-auto text-xs text-muted-foreground">
                  Active
                </span>
              )}
            </Command.Item>
          ))}
          <Command.Item
            value="toggle mode"
            onSelect={() => {
              toggleMode()
              setOpen(false)
            }}
          >
            {mode === "dark" ? (
              <Sun className="size-4" />
            ) : (
              <Moon className="size-4" />
            )}
            Switch to {mode === "dark" ? "Light" : "Dark"} Mode
          </Command.Item>
        </Command.Group>

        <Command.Group heading="Help">
          <Command.Item value="keyboard shortcuts">
            <Keyboard className="size-4" />
            Keyboard Shortcuts
          </Command.Item>
        </Command.Group>
      </Command.List>
    </Command.Dialog>
  )
}
