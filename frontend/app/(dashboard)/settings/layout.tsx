"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"

const settingsNav = [
  { href: "/settings/themes", label: "Themes" },
  { href: "/settings/plugins", label: "Plugins" },
]

export default function SettingsLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const pathname = usePathname()

  return (
    <div className="flex h-full">
      <aside className="w-48 border-r shrink-0 p-4 space-y-2">
        <h2 className="font-semibold text-sm tracking-tight mb-3">Settings</h2>
        {settingsNav.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "block text-sm rounded-md px-3 py-1.5 transition-colors",
              pathname === item.href
                ? "bg-accent text-accent-foreground font-medium"
                : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
            )}
          >
            {item.label}
          </Link>
        ))}
      </aside>
      <main className="flex-1 overflow-auto">{children}</main>
    </div>
  )
}
