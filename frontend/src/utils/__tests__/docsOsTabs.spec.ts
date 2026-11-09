import { beforeEach, describe, expect, it } from 'vitest'
import { DOCS_OS_STORAGE_KEY, enhanceDocsOsTabs } from '../docsOsTabs'

function mountDocs(): HTMLElement {
  document.body.innerHTML = `
    <div id="docs">
      <h2>Claude Code</h2>
      <p>claude-intro</p>
      <h3>安装</h3>
      <h4>macOS（zsh）</h4>
      <p>mac-install</p>
      <h4>Linux（bash）</h4>
      <p>linux-install</p>
      <h4>Windows PowerShell</h4>
      <p>windows-install</p>

      <h3>配置</h3>
      <h4>macOS</h4>
      <p>mac-config</p>
      <h4>Linux</h4>
      <p>linux-config</p>
      <h4>Windows</h4>
      <p>windows-config</p>

      <h2>Codex</h2>
      <p>codex-intro</p>
      <h3>安装</h3>
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

  it('groups adjacent OS sections and shows one master switcher per client', () => {
    localStorage.setItem(DOCS_OS_STORAGE_KEY, 'windows')
    const container = mountDocs()
    enhanceDocsOsTabs(container)

    expect(container.querySelectorAll('.os-tabs')).toHaveLength(3)
    expect(container.querySelectorAll('[role="tab"]')).toHaveLength(9)
    expect(container.querySelectorAll('.os-panel')).toHaveLength(9)
    expect(container.querySelector('h4')).toBeNull()

    expect(container.querySelectorAll('.os-tablist--master')).toHaveLength(2)
    expect(container.querySelectorAll('.os-tablist:not([hidden])')).toHaveLength(2)
    expect(container.querySelectorAll('.os-tablist[hidden]')).toHaveLength(1)

    const firstMaster = container.querySelector('.os-tablist--master')
    const firstStepHeading = container.querySelector('h3')
    expect(firstMaster?.compareDocumentPosition(firstStepHeading!) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()

    expect(
      Array.from(container.querySelectorAll('[role="tab"]'))
        .slice(0, 3)
        .map((tab) => tab.textContent?.trim()),
    ).toEqual(['Windows', 'Linux', 'macOS'])

    const visiblePanels = Array.from(container.querySelectorAll<HTMLElement>('.os-panel')).filter(
      (panel) => !panel.hidden,
    )
    expect(visiblePanels).toHaveLength(3)
    expect(visiblePanels.map((panel) => panel.textContent?.trim())).toEqual([
      'windows-install',
      'windows-config',
      'windows-codex',
    ])
  })

  it('keeps each client section on its own OS selection', () => {
    localStorage.setItem(DOCS_OS_STORAGE_KEY, 'windows')
    const container = mountDocs()
    enhanceDocsOsTabs(container)

    const claudeTablist = container.querySelectorAll<HTMLElement>('.os-tablist--master')[0]
    claudeTablist?.querySelector<HTMLButtonElement>('[data-os="linux"]')?.click()

    const selectedTabs = Array.from(container.querySelectorAll<HTMLButtonElement>('[role="tab"][aria-selected="true"]'))
    expect(selectedTabs).toHaveLength(3)
    expect(selectedTabs.filter((tab) => tab.dataset.os === 'linux')).toHaveLength(2)
    expect(selectedTabs.filter((tab) => tab.dataset.os === 'windows')).toHaveLength(1)
    expect(localStorage.getItem(`${DOCS_OS_STORAGE_KEY}:Claude Code`)).toBe('linux')

    const visiblePanels = Array.from(container.querySelectorAll<HTMLElement>('.os-panel')).filter(
      (panel) => !panel.hidden,
    )
    expect(visiblePanels.map((panel) => panel.textContent?.trim())).toEqual([
      'linux-install',
      'linux-config',
      'windows-codex',
    ])

    const codexTablist = container.querySelectorAll<HTMLElement>('.os-tablist--master')[1]
    codexTablist?.querySelector<HTMLButtonElement>('[data-os="macos"]')?.click()

    const afterCodex = Array.from(container.querySelectorAll<HTMLButtonElement>('[role="tab"][aria-selected="true"]'))
    expect(afterCodex.filter((tab) => tab.dataset.os === 'linux')).toHaveLength(2)
    expect(afterCodex.filter((tab) => tab.dataset.os === 'macos')).toHaveLength(1)
  })

  it('supports arrow-key navigation inside one client section', () => {
    localStorage.setItem(DOCS_OS_STORAGE_KEY, 'windows')
    const container = mountDocs()
    enhanceDocsOsTabs(container)

    const tablist = container.querySelector<HTMLElement>('.os-tablist--master')
    tablist?.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight', bubbles: true }))

    const selectedTabs = Array.from(container.querySelectorAll<HTMLButtonElement>('[role="tab"][aria-selected="true"]'))
    expect(selectedTabs.filter((tab) => tab.dataset.os === 'linux')).toHaveLength(2)
    expect(selectedTabs.filter((tab) => tab.dataset.os === 'windows')).toHaveLength(1)
    expect(localStorage.getItem(`${DOCS_OS_STORAGE_KEY}:Claude Code`)).toBe('linux')
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
