import { nextTick } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { enableAutoUnmount, flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createMemoryHistory, createRouter } from 'vue-router'
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

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }]
  })
}

async function mountBell(path = '/') {
  const router = createTestRouter()
  await router.push(path)
  await router.isReady()
  const wrapper = mount(AnnouncementBell, {
    global: { plugins: [router], stubs: { Teleport: true, Transition: true, Icon: true } }
  })
  // Only flush the render queue: awaiting the list request would drop the
  // seeded announcements before a test can assert on them.
  await nextTick()
  return { wrapper, router }
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
    const { wrapper } = await mountBell()
    await wrapper.get('button').trigger('click')
    await flushPromises()
    expect(listMock).toHaveBeenCalled()
    expect(wrapper.text()).toContain('Notice')
    expect(wrapper.find('.markdown-body').exists()).toBe(true)
  })
})

describe('announcement deep link', () => {
  it('opens the notification centre from ?announcements=1 and clears the parameter', async () => {
    const { wrapper, router } = await mountBell('/home?announcements=1')
    await flushPromises()
    expect(wrapper.text()).toContain('announcements.systemNotifications')
    expect(wrapper.text()).toContain('announcements.siteAnnouncement')
    expect(router.currentRoute.value.query.announcements).toBeUndefined()
    expect(router.currentRoute.value.path).toBe('/home')
  })

  it('opens the group QR-code tab from ?announcements=group', async () => {
    const { wrapper } = await mountBell('/home?announcements=group')
    expect(wrapper.find('img[src="/announcement/qq-group.png"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('扫码加入 QQ 群')
  })

  it('ignores a disabled parameter value', async () => {
    const { wrapper } = await mountBell('/home?announcements=0')
    await flushPromises()
    expect(wrapper.find('img[src="/announcement/qq-group.png"]').exists()).toBe(false)
  })
})
