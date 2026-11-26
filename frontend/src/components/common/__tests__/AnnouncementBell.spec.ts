import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { enableAutoUnmount, flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import AnnouncementBell from '../AnnouncementBell.vue'
import { useAnnouncementStore } from '@/stores/announcements'

// The bell refetches the list every time it opens, so `list` must resolve:
// an unmocked call rejects and wipes the seeded announcements.
const announcements = vi.hoisted(() => [
  { id: 1, title: 'Notice', content: 'Details', notify_mode: 'silent', created_at: '2026-09-30', updated_at: '2026-09-30' }
])
const listMock = vi.hoisted(() => vi.fn())
vi.mock('@/api', () => ({ announcementsAPI: { list: listMock, markRead: vi.fn() } }))
vi.mock('@/stores/app', () => ({ useAppStore: () => ({ showError: vi.fn(), showSuccess: vi.fn() }) }))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: (key: string) => key }) }))
vi.mock('@/utils/format', () => ({ formatRelativeTime: () => 'now', formatRelativeWithDateTime: () => 'now' }))
enableAutoUnmount(afterEach)

// Attach to the document so click events bubble up to document-level listeners
// (that is what the assertions on defaultPrevented observe).
function mountBell() {
  return mount(AnnouncementBell, {
    attachTo: document.body,
    global: { stubs: { Teleport: true, Transition: true, Icon: true } }
  })
}

function announcementWithLink(content: string) {
  return [{ id: 2, title: 'Promo', content, notify_mode: 'silent', created_at: '2026-09-30', updated_at: '2026-09-30' }]
}

beforeEach(() => {
  setActivePinia(createPinia())
  vi.clearAllMocks()
  listMock.mockResolvedValue(announcements)
  vi.spyOn(console, 'error').mockImplementation(() => {})
  useAnnouncementStore().announcements = [...announcements] as never
})
afterEach(() => { vi.restoreAllMocks(); document.body.style.overflow = '' })

describe('notification centre', () => {
  it('opens from the bell and renders the fetched announcements', async () => {
    const wrapper = mountBell()
    await wrapper.get('button').trigger('click')
    await flushPromises()
    expect(listMock).toHaveBeenCalled()
    expect(wrapper.text()).toContain('Notice')
    expect(wrapper.find('.markdown-body').exists()).toBe(true)
  })
})

describe('links inside an announcement body', () => {
  it('switches to the group tab instead of navigating away', async () => {
    listMock.mockResolvedValue(announcementWithLink('<p><a href="/home?announcements=group">点此入群</a></p>'))
    const wrapper = mountBell()
    await wrapper.get('button').trigger('click')
    await flushPromises()

    const preventDefault = vi.spyOn(Event.prototype, 'preventDefault')

    await wrapper.get('.markdown-body a').trigger('click')
    await flushPromises()

    expect(preventDefault).toHaveBeenCalled()
    expect(wrapper.find('img[src="/announcement/qq-group.png"]').exists()).toBe(true)
  })

  it('leaves ordinary links untouched', async () => {
    // A hash link keeps jsdom from attempting a real cross-document navigation.
    listMock.mockResolvedValue(announcementWithLink('<p><a href="#docs">文档</a></p>'))
    const wrapper = mountBell()
    await wrapper.get('button').trigger('click')
    await flushPromises()

    const preventDefault = vi.spyOn(Event.prototype, 'preventDefault')

    await wrapper.get('.markdown-body a').trigger('click')
    await flushPromises()

    expect(preventDefault).not.toHaveBeenCalled()
    expect(wrapper.find('img[src="/announcement/qq-group.png"]').exists()).toBe(false)
  })
})
