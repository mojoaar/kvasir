"use client"

import { Suspense } from "react"
import { useSearchParams, useRouter } from "next/navigation"
import { useQuery } from "@tanstack/react-query"
import { Input } from "@/components/ui/input"
import { Search as SearchIcon } from "lucide-react"
import { api } from "@/lib/api/client"

interface SearchResult {
  id: number
  title: string
  content: string
  snippet: string
  rank: number
  isFolder: boolean
  createdAt: string
}

function SearchContent() {
  const searchParams = useSearchParams()
  const router = useRouter()

  const q = searchParams.get("q") || ""

  const { data: results = [], isLoading } = useQuery({
    queryKey: ["search", q],
    queryFn: () =>
      api.get<SearchResult[]>(
        `/api/v1/search?q=${encodeURIComponent(q)}&limit=20`
      ),
    enabled: q.trim().length > 0,
  })

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)
    const value = formData.get("q") as string
    router.push(`/search?q=${encodeURIComponent(value)}`)
  }

  return (
    <div className="flex flex-col h-full max-w-3xl mx-auto px-6 py-8">
      <form onSubmit={handleSubmit} className="mb-6">
        <div className="relative">
          <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            type="search"
            name="q"
            placeholder="Search notes..."
            defaultValue={q}
            className="pl-10 h-10 text-sm"
            autoFocus
          />
        </div>
      </form>

      {isLoading && (
        <p className="text-sm text-muted-foreground">Searching...</p>
      )}

      {!isLoading && q.trim() && results.length === 0 && (
        <p className="text-sm text-muted-foreground">
          No results found for &ldquo;{q}&rdquo;
        </p>
      )}

      {results.length > 0 && (
        <div className="space-y-4">
          <p className="text-xs text-muted-foreground">
            {results.length} result{results.length !== 1 ? "s" : ""}
          </p>
          {results.map((result) => (
            <div key={result.id} className="space-y-1">
              <h3 className="text-sm font-medium">{result.title}</h3>
              {result.snippet && (
                <p
                  className="text-xs text-muted-foreground leading-relaxed"
                  dangerouslySetInnerHTML={{ __html: result.snippet }}
                />
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default function SearchPage() {
  return (
    <Suspense fallback={<div className="p-8 text-sm text-muted-foreground">Loading...</div>}>
      <SearchContent />
    </Suspense>
  )
}
