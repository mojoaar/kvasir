import { Node } from "@tiptap/react"

declare global {
  interface Window {
    mermaid?: {
      render: (id: string, code: string) => Promise<{ svg: string }>
    }
  }
}

export const MermaidExtension = Node.create({
  name: "mermaid",

  group: "block",
  atom: true,

  addAttributes() {
    return {
      code: { default: "" },
      svg: { default: "" },
    }
  },

  parseHTML() {
    return [{ tag: "div[data-mermaid]" }]
  },

  renderHTML({ node }) {
    const svg = node.attrs.svg as string
    const code = node.attrs.code as string
    if (svg) {
      return ["div", { "data-mermaid": "" }, svg]
    }
    return [
      "div",
      { "data-mermaid": "" },
      ["pre", {}, ["code", {}, code]],
    ]
  },

  addNodeView() {
    return ({ node, editor }) => {
      const dom = document.createElement("div")
      dom.setAttribute("data-mermaid", "")

      const pre = document.createElement("pre")
      const code = document.createElement("code")
      code.textContent = node.attrs.code as string
      code.contentEditable = "true"

      code.addEventListener("input", () => {
        const pos = editor.state.selection.from
        const resolved = editor.state.doc.resolve(pos)
        const parentPos = resolved.start()

        editor
          .chain()
          .setNodeSelection(parentPos)
          .updateAttributes("mermaid", { code: code.textContent || "" })
          .run()
      })

      pre.appendChild(code)
      dom.appendChild(pre)

      const renderBtn = document.createElement("button")
      renderBtn.textContent = "Render"
      renderBtn.className =
        "mt-1 px-2 py-0.5 text-xs rounded bg-primary text-primary-foreground hover:opacity-90"
      renderBtn.addEventListener("click", async () => {
        if (window.mermaid) {
          try {
            const id = `mermaid-${Date.now()}`
            const { svg } = await window.mermaid.render(id, code.textContent || "")

            const resolved = editor.state.doc.resolve(
              editor.state.selection.from
            )
            const parentPos = resolved.start()

            editor
              .chain()
              .setNodeSelection(parentPos)
              .updateAttributes("mermaid", { code: code.textContent || "", svg })
              .run()
          } catch {
            // render error — keep code
          }
        }
      })
      dom.appendChild(renderBtn)

      return { dom }
    }
  },
})
