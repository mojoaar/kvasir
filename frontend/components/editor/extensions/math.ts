import { Node } from "@tiptap/react"

declare global {
  interface Window {
    katex?: {
      renderToString: (latex: string, options?: Record<string, unknown>) => string
    }
  }
}

export const MathExtension = Node.create({
  name: "math",

  group: "inline",
  inline: true,
  atom: true,

  addAttributes() {
    return {
      latex: { default: "" },
    }
  },

  parseHTML() {
    return [{ tag: "span[data-math]" }]
  },

  renderHTML({ node }) {
    const latex = node.attrs.latex as string
    if (typeof window !== "undefined" && window.katex) {
      try {
        const html = window.katex.renderToString(latex, {
          throwOnError: false,
          displayMode: false,
        })
        return ["span", { "data-math": "" }, html]
      } catch {
        return ["span", { "data-math": "" }, latex]
      }
    }
    return ["span", { "data-math": "" }, latex]
  },
})

export const DisplayMathExtension = Node.create({
  name: "displayMath",

  group: "block",
  atom: true,

  addAttributes() {
    return {
      latex: { default: "" },
    }
  },

  parseHTML() {
    return [{ tag: "div[data-display-math]" }]
  },

  renderHTML({ node }) {
    const latex = node.attrs.latex as string
    if (typeof window !== "undefined" && window.katex) {
      try {
        const html = window.katex.renderToString(latex, {
          throwOnError: false,
          displayMode: true,
        })
        return ["div", { "data-display-math": "" }, html]
      } catch {
        return ["div", { "data-display-math": "" }, latex]
      }
    }
    return ["div", { "data-display-math": "" }, latex]
  },
})
