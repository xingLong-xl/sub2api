export type DocsOsKey = 'macos' | 'linux' | 'windows'

export const DOCS_OS_STORAGE_KEY = 'tokenpilot.docs.os'

const OS_LABELS: Record<DocsOsKey, string> = {
  macos: 'macOS',
  linux: 'Linux',
  windows: 'Windows',
}

const OS_ORDER: DocsOsKey[] = ['windows', 'linux', 'macos']

interface OsSection {
  key: DocsOsKey
  heading: HTMLElement
  nodes: HTMLElement[]
}

interface OsTabGroup {
  keys: Set<DocsOsKey>
  buttons: Map<DocsOsKey, HTMLButtonElement>
  panels: Map<DocsOsKey, HTMLElement>
  wrapper: HTMLElement
  tablist: HTMLElement
}

function detectOsKey(heading: HTMLElement): DocsOsKey | null {
  if (heading.tagName !== 'H4') return null
  const text = (heading.textContent || '').replace(/\s+/g, ' ').trim()
  if (/^macOS\b/i.test(text)) return 'macos'
  if (/^Linux\b/i.test(text)) return 'linux'
  if (/^Windows\b/i.test(text)) return 'windows'
  return null
}

function isHeading(element: HTMLElement): boolean {
  return /^H[1-6]$/.test(element.tagName)
}

function getStoredOs(container: HTMLElement): DocsOsKey | null {
  try {
    const value = container.ownerDocument.defaultView?.localStorage.getItem(DOCS_OS_STORAGE_KEY)
    return OS_ORDER.includes(value as DocsOsKey) ? (value as DocsOsKey) : null
  } catch {
    return null
  }
}

function storeOs(container: HTMLElement, key: DocsOsKey): void {
  try {
    container.ownerDocument.defaultView?.localStorage.setItem(DOCS_OS_STORAGE_KEY, key)
  } catch {
    // Ignore storage failures (private mode, quota, etc.).
  }
}

function detectPreferredOs(container: HTMLElement): DocsOsKey {
  const userAgent = container.ownerDocument.defaultView?.navigator.userAgent || ''
  if (/Macintosh|Mac OS X/i.test(userAgent)) return 'macos'
  if (/Linux/i.test(userAgent) && !/Android/i.test(userAgent)) return 'linux'
  return 'windows'
}

function buildGroup(container: HTMLElement, sections: OsSection[], groupIndex: number): OsTabGroup {
  const doc = container.ownerDocument
  const wrapper = doc.createElement('div')
  wrapper.className = 'os-tabs'

  const tablist = doc.createElement('div')
  tablist.className = 'os-tablist'
  tablist.setAttribute('role', 'tablist')
  tablist.setAttribute('aria-label', '选择操作系统')
  wrapper.appendChild(tablist)

  const buttons = new Map<DocsOsKey, HTMLButtonElement>()
  const panels = new Map<DocsOsKey, HTMLElement>()

  // Keep the visual order Windows → Linux → macOS regardless of the order the
  // sections appear in Markdown, and merge duplicate OS headings if present.
  const byKey = new Map<DocsOsKey, OsSection>()
  for (const section of sections) {
    const existing = byKey.get(section.key)
    if (existing) {
      existing.nodes.push(...section.nodes)
      section.heading.remove()
    } else {
      byKey.set(section.key, { ...section, nodes: [...section.nodes] })
    }
  }

  for (const section of OS_ORDER.map((key) => byKey.get(key)).filter((item): item is OsSection => Boolean(item))) {
    const tabId = `docs-os-tab-${groupIndex}-${section.key}`
    const panelId = `docs-os-panel-${groupIndex}-${section.key}`

    const button = doc.createElement('button')
    button.type = 'button'
    button.className = 'os-tab'
    button.dataset.os = section.key
    button.id = tabId
    button.textContent = OS_LABELS[section.key]
    button.setAttribute('role', 'tab')
    button.setAttribute('aria-controls', panelId)
    button.setAttribute('aria-selected', 'false')
    button.tabIndex = -1
    tablist.appendChild(button)

    const panel = doc.createElement('div')
    panel.className = 'os-panel'
    panel.id = panelId
    panel.setAttribute('role', 'tabpanel')
    panel.setAttribute('aria-labelledby', tabId)
    panel.hidden = true
    for (const node of section.nodes) panel.appendChild(node)
    wrapper.appendChild(panel)

    buttons.set(section.key, button)
    panels.set(section.key, panel)
  }

  const firstHeading = sections[0]?.heading
  firstHeading?.parentNode?.insertBefore(wrapper, firstHeading)
  for (const section of sections) section.heading.remove()

  return { keys: new Set(buttons.keys()), buttons, panels, wrapper, tablist }
}

function findScopeHeading(group: OsTabGroup): HTMLElement | null {
  let node = group.wrapper.previousElementSibling
  while (node) {
    if (node.tagName === 'H2') return node as HTMLElement
    if (node.tagName === 'H1') return null
    node = node.previousElementSibling
  }
  return null
}

/**
 * One client section can contain several OS-dependent blocks (install, env,
 * settings...). Keep a single tab bar at the top of that section and let it
 * control every block, instead of repeating the switcher above each block.
 */
function arrangeMasterTabLists(groups: OsTabGroup[]): void {
  const scopes = new Map<HTMLElement, OsTabGroup[]>()

  for (const group of groups) {
    const scope = findScopeHeading(group)
    if (!scope) continue
    const scoped = scopes.get(scope) ?? []
    scoped.push(group)
    scopes.set(scope, scoped)
  }

  for (const [scope, scoped] of scopes) {
    const first = scoped[0]
    if (!first) continue

    let insertBefore: Element | null = null
    let node = scope.nextElementSibling
    while (node && node.tagName !== 'H2') {
      if (node.tagName === 'H3') {
        insertBefore = node
        break
      }
      node = node.nextElementSibling
    }

    const target = insertBefore ?? first.wrapper
    target.parentNode?.insertBefore(first.tablist, target)
    first.tablist.classList.add('os-tablist--master')

    // The master tab now controls every OS panel in this client section.
    for (const group of scoped) {
      group.panels.forEach((panel, key) => {
        const masterButton = first.buttons.get(key)
        if (!masterButton) return
        const controlled = new Set((masterButton.getAttribute('aria-controls') || '').split(/\s+/).filter(Boolean))
        controlled.add(panel.id)
        masterButton.setAttribute('aria-controls', Array.from(controlled).join(' '))
        panel.setAttribute('aria-labelledby', masterButton.id)
      })
    }

    for (const group of scoped.slice(1)) group.tablist.hidden = true
  }
}

function applySelection(group: OsTabGroup, key: DocsOsKey): void {
  group.buttons.forEach((button, buttonKey) => {
    const selected = buttonKey === key
    button.setAttribute('aria-selected', String(selected))
    button.tabIndex = selected ? 0 : -1
  })
  group.panels.forEach((panel, panelKey) => {
    panel.hidden = panelKey !== key
  })
}

function bindKeyboard(tablist: HTMLElement, group: OsTabGroup, activate: (key: DocsOsKey) => void): void {
  tablist.addEventListener('keydown', (event) => {
    const keys = OS_ORDER.filter((key) => group.keys.has(key))
    if (keys.length === 0) return

    const current = keys.find((key) => group.buttons.get(key)?.getAttribute('aria-selected') === 'true') || keys[0]
    const currentIndex = keys.indexOf(current)
    let nextIndex: number | null = null

    if (event.key === 'ArrowRight' || event.key === 'ArrowDown') {
      nextIndex = (currentIndex + 1) % keys.length
    } else if (event.key === 'ArrowLeft' || event.key === 'ArrowUp') {
      nextIndex = (currentIndex - 1 + keys.length) % keys.length
    } else if (event.key === 'Home') {
      nextIndex = 0
    } else if (event.key === 'End') {
      nextIndex = keys.length - 1
    }

    if (nextIndex === null) return
    event.preventDefault()
    const nextKey = keys[nextIndex]
    activate(nextKey)
    group.buttons.get(nextKey)?.focus()
  })
}

/**
 * Turns adjacent `#### macOS` / `#### Linux` / `#### Windows` sections into
 * accessible tabs. Selecting an OS in one group switches every group on the
 * page, so users only choose their system once.
 */
export function enhanceDocsOsTabs(container: HTMLElement): void {
  if (container.dataset.osTabsReady === 'true') return

  const children = Array.from(container.children) as HTMLElement[]
  const groups: OsTabGroup[] = []
  let index = 0

  while (index < children.length) {
    const firstKey = detectOsKey(children[index])
    if (!firstKey) {
      index += 1
      continue
    }

    const sections: OsSection[] = []
    let cursor = index

    while (cursor < children.length) {
      const heading = children[cursor]
      const key = detectOsKey(heading)
      if (!key) break

      const nodes: HTMLElement[] = []
      let next = cursor + 1
      while (next < children.length && !isHeading(children[next])) {
        nodes.push(children[next])
        next += 1
      }

      sections.push({ key, heading, nodes })
      cursor = next

      if (cursor >= children.length) break
      if (!detectOsKey(children[cursor])) break
    }

    if (sections.length < 2) {
      index += 1
      continue
    }

    groups.push(buildGroup(container, sections, groups.length))
    index = cursor
  }

  if (groups.length === 0) {
    container.dataset.osTabsReady = 'true'
    return
  }

  arrangeMasterTabLists(groups)

  const available = new Set<DocsOsKey>()
  groups.forEach((group) => group.keys.forEach((key) => available.add(key)))

  const stored = getStoredOs(container)
  const preferred = stored ?? detectPreferredOs(container)
  const initial = available.has(preferred) ? preferred : (groups[0] && OS_ORDER.find((key) => groups[0].keys.has(key))) || 'windows'

  const activate = (key: DocsOsKey): void => {
    storeOs(container, key)
    groups.forEach((group) => {
      if (group.keys.has(key)) applySelection(group, key)
    })
  }

  for (const group of groups) {
    group.buttons.forEach((button, key) => {
      button.addEventListener('click', () => activate(key))
    })
    const tablist = group.buttons.values().next().value?.parentElement
    if (tablist) bindKeyboard(tablist, group, activate)
  }

  activate(initial)
  container.dataset.osTabsReady = 'true'
}
