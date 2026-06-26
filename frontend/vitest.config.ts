import { defineConfig } from "vitest/config"
import react from "@vitejs/plugin-react"
import path from "path"

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    setupFiles: ["./tests/setup.ts"],
    include: ["**/*.test.{ts,tsx}"],
    coverage: {
      provider: "v8",
      reporter: ["text", "json-summary"],
      include: ["lib/store/**/*.ts", "lib/api/client.ts", "components/command-palette/index.tsx"],
      exclude: ["components/ui/**", "lib/utils.ts"],
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname),
    },
  },
})
