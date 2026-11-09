import { beforeEach, describe, expect, it } from 'vitest'
import { DOCS_OS_STORAGE_KEY, enhanceDocsOsTabs } from '../docsOsTabs'

function mountDocs(): HTMLElement {
  document.body.innerHTML = `
    <div id="docs">
      <h3>Claude Code</h3>
      <h4>macOS（zsh）</h4>
      <p>mac-install</p>
      <h4>Linux（bash）</h4>
      <p>linux-install</p>
      <h4>Windows PowerShell</h4>
      <p>windows-install</p>

      <h3>Codex</h3>
      <h4>macOS</h4>
      <p>mac-codex</p>
      <h4>Linux</h4>
      <p>linux-codex</p>
      <h4>Windows</h4>
      <p>windows-codex</p>
    </div>
  `
  return document.querySelector('#docs') as HTMLElement
}

describe('enhanceDocsOsTabs', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('groups adjacent OS sections into tabs and keeps one panel visible', () => {
    localStorage.setItem(DOCS_OS_STORAGE_KEY, 'windows')
    const container = mountDocs()
    enhanceDocsOsTabs(container)

    expect(container.querySelectorAll('.os-tabs')).toHaveLength(2)
    expect(container.querySelectorAll('[role="tab"]')).toHaveLength(6)
    expect(container.querySelectorAll('.os-panel')).toHaveLength(6)
    expect(container.querySelector('h4')).toBeNull()
    expect(
      Array.from(container.querySelectorAll('[role="tab"]'))
        .slice(0, 3)
        .map((tab) => tab.textContent?.trim()),
    ).toEqual(['Windows', 'Linux', 'macOS'])

    const visiblePanels = Array.from(container.querySelectorAll<HTMLElement>('.os-panel')).filter(
      (panel) => !panel.hidden,
    )
    expect(visiblePanels).toHaveLength(2)
    expect(visiblePanels.map((panel) => panel.textContent?.trim())).toEqual(['windows-install', 'windows-codex'])
  })

  it('switches every group to the selected OS', () => {
    const container = mountDocs()
    enhanceDocsOsTabs(container)

    const linuxTab = container.querySelector<HTMLButtonElement>('[role="tab"][data-os="linux"]')
    linuxTab?.click()

    const selectedTabs = Array.from(container.querySelectorAll<HTMLButtonElement>('[role="tab"][aria-selected="true"]'))
    expect(selectedTabs).toHaveLength(2)
    expect(selectedTabs.every((tab) => tab.dataset.os === 'linux')).toBe(true)
    expect(localStorage.getItem(DOCS_OS_STORAGE_KEY)).toBe('linux')

    const visiblePanels = Array.from(container.querySelectorAll<HTMLElement>('.os-panel')).filter(
      (panel) => !panel.hidden,
    )
    expect(visiblePanels.map((panel) => panel.textContent?.trim())).toEqual(['linux-install', 'linux-codex'])
  })

  it('supports arrow-key navigation and remembers the selection', () => {
    localStorage.setItem(DOCS_OS_STORAGE_KEY, 'windows')
    const container = mountDocs()
    enhanceDocsOsTabs(container)

    const tablist = container.querySelector<HTMLElement>('.os-tablist')
    tablist?.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight', bubbles: true }))

    const selectedTabs = Array.from(container.querySelectorAll<HTMLButtonElement>('[role="tab"][aria-selected="true"]'))
    expect(selectedTabs.every((tab) => tab.dataset.os === 'linux')).toBe(true)
    expect(localStorage.getItem(DOCS_OS_STORAGE_KEY)).toBe('linux')
  })

  it('does not create tabs for a single OS section', () => {
    document.body.innerHTML = `
      <div id="docs">
        <h3>安装</h3>
        <h4>Windows</h4>
        <p>only windows</p>
      </div>
    `
    const container = document.querySelector('#docs') as HTMLElement
    enhanceDocsOsTabs(container)

    expect(container.querySelector('.os-tabs')).toBeNull()
    expect(container.querySelector('h4')?.textContent).toBe('Windows')
  })
})
