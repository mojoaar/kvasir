import "@testing-library/jest-dom/vitest"

const store: Record<string, string> = {}

const localStorageMock: Storage = {
  getItem: (key: string) => store[key] ?? null,
  setItem: (key: string, value: string) => { store[key] = value },
  removeItem: (key: string) => { delete store[key] },
  clear: () => { Object.keys(store).forEach((k) => delete store[k]) },
  key: (index: number) => Object.keys(store)[index] ?? null,
  length: 0,
}

Object.defineProperty(localStorageMock, "length", {
  get: () => Object.keys(store).length,
})

Object.defineProperty(globalThis, "localStorage", {
  value: localStorageMock,
  writable: true,
})

globalThis.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

Element.prototype.scrollIntoView = () => {}

