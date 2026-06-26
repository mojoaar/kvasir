"use client"

import { useState, type FormEvent } from "react"
import Link from "next/link"
import { usePathname, useRouter } from "next/navigation"
import { cn } from "@/lib/utils"
import { useUIStore } from "@/lib/store/ui-store"
import {
  Search,
  Settings,
  PanelLeftClose,
  PanelLeft,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import { ThemeToggle } from "@/components/themes/theme-toggle"
import { ThemeSelector } from "@/components/themes/theme-selector"
import { NoteTree } from "@/components/sidebar/note-tree"

const navItems = [
  {
    href: "/search",
    label: "Search",
    icon: Search,
  },
  {
    href: "/settings",
    label: "Settings",
    icon: Settings,
  },
]

export function Sidebar() {
  const pathname = usePathname()
  const router = useRouter()
  const { sidebarOpen, toggleSidebar } = useUIStore()
  const [searchQuery, setSearchQuery] = useState("")

  const handleSearch = (e: FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      router.push(`/search?q=${encodeURIComponent(searchQuery.trim())}`)
    }
  }

  if (!sidebarOpen) {
    return (
      <aside className="w-12 border-r flex flex-col items-center py-3 shrink-0">
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8"
          onClick={toggleSidebar}
        >
          <PanelLeft className="h-4 w-4" />
        </Button>
      </aside>
    )
  }

  return (
    <aside className="w-60 border-r flex flex-col h-full shrink-0">
      <div className="flex items-center justify-between px-3 py-3">
        <Link href="/" className="font-semibold text-sm tracking-tight">
          Kvasir
        </Link>
        <div className="flex items-center gap-1">
          <ThemeSelector />
          <ThemeToggle />
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7"
            onClick={toggleSidebar}
          >
            <PanelLeftClose className="h-4 w-4" />
          </Button>
        </div>
      </div>
      <Separator />
      <nav className="px-2 py-1.5 space-y-0.5">
        {navItems.map((item) => {
          const isActive =
            pathname === item.href || pathname.startsWith(`${item.href}/`)
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-2 rounded-md px-3 py-1 text-xs font-medium transition-colors",
                isActive
                  ? "bg-accent text-accent-foreground"
                  : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
              )}
            >
              <item.icon className="h-3.5 w-3.5" />
              {item.label}
            </Link>
          )
        })}
      </nav>
      <div className="px-2 py-2">
        <form onSubmit={handleSearch}>
          <div className="relative">
            <Search className="absolute left-2 top-1/2 -translate-y-1/2 h-3 w-3 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-7 h-7 text-xs"
            />
          </div>
        </form>
      </div>
      <Separator />
      <div className="flex-1 min-h-0">
        <NoteTree />
      </div>
    </aside>
  )
}
