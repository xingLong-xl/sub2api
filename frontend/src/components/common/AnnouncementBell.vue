<template>
  <div>
    <!-- 铃铛按钮 -->
    <button
      @click="openModal"
      class="relative flex h-9 items-center gap-1.5 rounded-lg px-2.5 text-gray-600 transition-all hover:bg-gray-100 hover:scale-105 dark:text-gray-400 dark:hover:bg-dark-800"
      :class="{ 'text-blue-600 dark:text-blue-400': unreadCount > 0 }"
      :aria-label="t('announcements.title')"
    >
      <span class="relative inline-flex h-5 w-5 items-center justify-center">
        <Icon name="bell" size="md" />
        <!-- 未读红点 -->
        <span
          v-if="unreadCount > 0"
          class="absolute -right-1 -top-1 flex h-2 w-2"
        >
          <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-red-500 opacity-75"></span>
          <span class="relative inline-flex h-2 w-2 rounded-full bg-red-500"></span>
        </span>
      </span>
      <span class="hidden whitespace-nowrap text-sm font-medium sm:inline">
        {{ t('announcements.systemLabel') }}
      </span>
    </button>

    <!-- 通知中心 Modal -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          v-if="isModalOpen"
          class="fixed inset-0 z-[100] flex items-center justify-center bg-black/35 p-4 backdrop-blur-[2px] sm:p-6"
          @click="closeModal"
        >
          <div
            class="flex h-[78vh] max-h-[780px] min-h-[520px] w-full max-w-[720px] flex-col overflow-hidden rounded-2xl bg-white shadow-2xl ring-1 ring-black/10 dark:bg-dark-800 dark:ring-white/10"
            @click.stop
          >
            <!-- Header -->
            <div class="flex items-center justify-between border-b border-gray-100 px-6 py-5 dark:border-dark-700">
              <div class="flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-primary-500 to-primary-600 text-white shadow-lg shadow-primary-500/20">
                  <Icon name="bell" size="md" />
                </div>
                <div>
                  <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
                    {{ t('announcements.notificationCenter') }}
                  </h2>
                </div>
              </div>
              <button
                @click="closeModal"
                class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-700 dark:hover:text-white"
                :aria-label="t('common.close')"
              >
                <Icon name="x" size="md" />
              </button>
            </div>

            <!-- Tabs -->
            <div class="grid w-full grid-cols-2 border-b border-gray-200 dark:border-dark-700">
              <button
                class="relative flex items-center justify-center gap-2 border-b-2 px-4 py-2 text-sm font-medium transition-colors first:border-r first:border-gray-100 dark:first:border-dark-700"
                :class="activeTab === 'notifications'
                  ? 'border-b-transparent text-primary-600 dark:text-primary-400'
                  : 'border-b-transparent text-gray-500 hover:bg-gray-50 hover:text-gray-800 dark:text-dark-400 dark:hover:bg-dark-800/60 dark:hover:text-white'"
                @click="activeTab = 'notifications'"
              >
                {{ t('announcements.systemNotifications') }}
                <span v-if="activeTab === 'notifications'" class="absolute inset-x-0 bottom-0 h-0.5 rounded-full bg-primary-500"></span>
                <span v-if="unreadCount > 0" class="rounded-full bg-primary-100 px-1.5 py-0.5 text-[10px] font-semibold text-primary-700 dark:bg-primary-900/40 dark:text-primary-300">
                  {{ unreadCount }}
                </span>
              </button>
              <button
                class="relative flex items-center justify-center gap-2 border-b-2 px-4 py-2 text-sm font-medium transition-colors first:border-r first:border-gray-100 dark:first:border-dark-700"
                :class="activeTab === 'announcement'
                  ? 'border-b-transparent text-primary-600 dark:text-primary-400'
                  : 'border-b-transparent text-gray-500 hover:bg-gray-50 hover:text-gray-800 dark:text-dark-400 dark:hover:bg-dark-800/60 dark:hover:text-white'"
                @click="activeTab = 'announcement'"
              >
                {{ t('announcements.siteAnnouncement') }}
                <span v-if="activeTab === 'announcement'" class="absolute inset-x-0 bottom-0 h-0.5 rounded-full bg-primary-500"></span>
              </button>
            </div>

            <!-- Body -->
            <div class="min-h-0 flex-1 overflow-y-auto px-6 py-5 [scrollbar-gutter:stable]">
              <div v-if="activeTab === 'notifications'">
                <div v-if="loading" class="flex items-center justify-center py-20">
                  <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
                </div>
                <div v-else-if="announcements.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
                  <div class="mb-3 flex h-14 w-14 items-center justify-center rounded-full bg-gray-100 text-gray-400 dark:bg-dark-700 dark:text-dark-500">
                    <Icon name="inbox" size="lg" />
                  </div>
                  <p class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('announcements.empty') }}</p>
                  <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">{{ t('announcements.emptyDescription') }}</p>
                </div>
                <div v-else class="space-y-2">
                  <article
                    v-for="item in announcements"
                    :key="item.id"
                    class="rounded-xl border border-gray-100 bg-white px-4 py-4 dark:border-dark-700 dark:bg-dark-800/40"
                  >
                    <div class="flex items-start gap-3">
                      <span class="mt-2 flex-shrink-0">
                        <span v-if="!item.read_at" class="block h-2 w-2 rounded-full bg-primary-500"></span>
                        <span v-else class="block h-2 w-2 rounded-full bg-gray-200 dark:bg-dark-600"></span>
                      </span>
                      <div class="min-w-0 flex-1">
                        <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ item.title }}</h3>
                        <div class="mt-1 text-[11px] text-gray-400 dark:text-dark-500">
                          {{ formatRelativeWithDateTime(item.created_at) }}
                        </div>
                      </div>
                    </div>
                    <div
                      class="markdown-body prose prose-sm mt-3 max-w-none text-sm leading-6 text-gray-700 dark:prose-invert dark:text-dark-200"
                      v-html="renderMarkdown(item.content)"
                      @click="handleMarkdownLinkClick"
                    ></div>
                  </article>
                </div>
              </div>

              <div v-else class="space-y-5">
                <div class="rounded-2xl border border-primary-100 bg-gradient-to-br from-primary-50/80 via-white to-indigo-50/50 p-5 dark:border-primary-900/40 dark:from-primary-900/15 dark:via-dark-800 dark:to-indigo-900/10">
                  <div class="flex flex-wrap items-start justify-between gap-4">
                    <div>
                      <p class="text-xs font-semibold uppercase tracking-wider text-primary-600 dark:text-primary-400">
                        {{ t('announcements.officialGroup') }}
                      </p>
                      <h3 class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">
                        TokenPilot 官方交流群
                      </h3>
                      <p class="mt-2 max-w-xl text-sm leading-6 text-gray-600 dark:text-dark-300">
                        欢迎加入我们的用户交流群，获取模型动态、接口教程、使用支持和站点公告。
                      </p>
                    </div>
                    <span class="rounded-full bg-white/80 px-3 py-1 text-xs font-medium text-primary-700 ring-1 ring-primary-100 dark:bg-dark-800/80 dark:text-primary-300 dark:ring-primary-900/50">
                      QQ 群：584779326
                    </span>
                  </div>

                  <div class="mt-5 grid gap-3 text-sm text-gray-600 dark:text-dark-300 sm:grid-cols-2">
                    <p>· GPT / Claude / Codex 接入交流</p>
                    <p>· API Key、余额、用量问题反馈</p>
                    <p>· 模型状态与优惠活动通知</p>
                    <p>· 开发者、学生、项目用户互助</p>
                  </div>
                </div>

                <div class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-gray-200 bg-gray-50/60 px-5 py-7 text-center dark:border-dark-600 dark:bg-dark-900/20">
                  <p class="text-sm font-medium text-gray-700 dark:text-dark-200">扫码加入 QQ 群</p>
                  <div class="flex h-52 w-52 items-center justify-center overflow-hidden rounded-2xl border border-gray-100 bg-white p-3 shadow-lg shadow-primary-500/10 ring-1 ring-black/5 dark:border-dark-600">
                    <img
                      v-if="!qrLoadFailed"
                      :src="qrImageUrl"
                      alt="QQ 群二维码"
                      class="h-full w-full object-contain"
                      loading="lazy"
                      @error="qrLoadFailed = true"
                    />
                    <div v-else class="flex flex-col items-center gap-2 text-xs text-gray-400 dark:text-dark-500">
                      <Icon name="grid" size="lg" />
                      <span>二维码待上传</span>
                      <span class="font-mono text-[10px]">frontend/public/announcement/qq-group.png</span>
                    </div>
                  </div>
                  <p class="text-[15px] text-gray-500 dark:text-dark-400">
                    官网：<a href="https://tokenpilot.site" target="_blank" rel="noopener noreferrer" class="text-primary-600 hover:underline dark:text-primary-400">tokenpilot.site</a>
                  </p>
                </div>

                <p class="text-center text-xs text-gray-400 dark:text-dark-500">
                  国内直连 · 多模型支持 · 稳定 API 中转 · 适合开发者、学生和项目接入使用
                </p>
              </div>
            </div>

            <!-- Footer -->
            <div class="flex items-center justify-between border-t border-gray-100 bg-gray-50/60 px-6 py-4 dark:border-dark-700 dark:bg-dark-900/20">
              <span></span>
              <div class="flex items-center gap-2">
                <button
                  class="rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:border-dark-600 dark:bg-dark-700 dark:text-dark-300 dark:hover:bg-dark-600 dark:hover:text-white"
                  @click="closeForToday"
                >
                  {{ t('announcements.closeToday') }}
                </button>
                <button
                  class="rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
                  @click="closeModal"
                >
                  {{ t('announcements.close') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- 公告详情 Modal -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          v-if="detailModalOpen && selectedAnnouncement"
          class="fixed inset-0 z-[110] flex items-start justify-center overflow-y-auto bg-gradient-to-br from-black/70 via-black/60 to-black/70 p-4 pt-[6vh] backdrop-blur-md"
          @click="closeDetail"
        >
          <div
            class="w-full max-w-[780px] overflow-hidden rounded-3xl bg-white shadow-2xl ring-1 ring-black/5 dark:bg-dark-800 dark:ring-white/10"
            @click.stop
          >
            <!-- Header with Decorative Elements -->
            <div class="relative overflow-hidden border-b border-gray-100 bg-gradient-to-br from-blue-50/80 via-indigo-50/50 to-purple-50/30 px-8 py-6 dark:border-dark-700 dark:from-blue-900/20 dark:via-indigo-900/10 dark:to-purple-900/5">
              <!-- Decorative background elements -->
              <div class="absolute right-0 top-0 h-full w-64 bg-gradient-to-l from-indigo-100/30 to-transparent dark:from-indigo-900/20"></div>
              <div class="absolute -right-8 -top-8 h-32 w-32 rounded-full bg-gradient-to-br from-blue-400/20 to-indigo-500/20 blur-3xl"></div>
              <div class="absolute -left-4 -bottom-4 h-24 w-24 rounded-full bg-gradient-to-tr from-purple-400/20 to-pink-500/20 blur-2xl"></div>

              <div class="relative z-10 flex items-start justify-between gap-4">
                <div class="flex-1 min-w-0">
                  <!-- Icon and Category -->
                  <div class="mb-3 flex items-center gap-2">
                    <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 text-white shadow-lg shadow-blue-500/30">
                      <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="rounded-lg bg-blue-100 px-2.5 py-1 text-xs font-medium text-blue-700 dark:bg-blue-900/40 dark:text-blue-300">
                        {{ t('announcements.title') }}
                      </span>
                      <span
                        v-if="!selectedAnnouncement.read_at"
                        class="inline-flex items-center gap-1.5 rounded-lg bg-gradient-to-r from-blue-500 to-indigo-600 px-2.5 py-1 text-xs font-medium text-white shadow-lg shadow-blue-500/30"
                      >
                        <span class="relative flex h-2 w-2">
                          <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-white opacity-75"></span>
                          <span class="relative inline-flex h-2 w-2 rounded-full bg-white"></span>
                        </span>
                        {{ t('announcements.unread') }}
                      </span>
                    </div>
                  </div>

                  <!-- Title -->
                  <h2 class="mb-3 text-2xl font-bold leading-tight text-gray-900 dark:text-white">
                    {{ selectedAnnouncement.title }}
                  </h2>

                  <!-- Meta Info -->
                  <div class="flex items-center gap-4 text-sm text-gray-600 dark:text-gray-400">
                    <div class="flex items-center gap-1.5">
                      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <time>{{ formatRelativeWithDateTime(selectedAnnouncement.created_at) }}</time>
                    </div>
                    <div class="flex items-center gap-1.5">
                      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                      <span>{{ selectedAnnouncement.read_at ? t('announcements.read') : t('announcements.unread') }}</span>
                    </div>
                  </div>
                </div>

                <!-- Close button -->
                <button
                  @click="closeDetail"
                  class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-xl bg-white/50 text-gray-500 backdrop-blur-sm transition-all hover:bg-white hover:text-gray-700 hover:shadow-lg dark:bg-dark-700/50 dark:text-gray-400 dark:hover:bg-dark-700 dark:hover:text-gray-300"
                  :aria-label="t('common.close')"
                >
                  <Icon name="x" size="md" />
                </button>
              </div>
            </div>

            <!-- Body with Enhanced Markdown -->
            <div class="max-h-[60vh] overflow-y-auto bg-white px-8 py-8 dark:bg-dark-800">
              <!-- Content with decorative border -->
              <div class="relative">
                <!-- Decorative left border -->
                <div class="absolute left-0 top-0 bottom-0 w-1 rounded-full bg-gradient-to-b from-blue-500 via-indigo-500 to-purple-500"></div>

                <div class="pl-6">
                  <div
                    class="markdown-body prose prose-sm max-w-none dark:prose-invert"
                    v-html="renderMarkdown(selectedAnnouncement.content)"
                  ></div>
                </div>
              </div>
            </div>

            <!-- Footer with Actions -->
            <div class="border-t border-gray-100 bg-gray-50/50 px-8 py-5 dark:border-dark-700 dark:bg-dark-900/30">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
                  <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span>{{ selectedAnnouncement.read_at ? t('announcements.readStatus') : t('announcements.markReadHint') }}</span>
                </div>
                <div class="flex items-center gap-3">
                  <button
                    @click="closeDetail"
                    class="rounded-xl border border-gray-300 bg-white px-5 py-2.5 text-sm font-medium text-gray-700 shadow-sm transition-all hover:bg-gray-50 hover:shadow dark:border-dark-600 dark:bg-dark-700 dark:text-gray-300 dark:hover:bg-dark-600"
                  >
                    {{ t('common.close') }}
                  </button>
                  <button
                    v-if="!selectedAnnouncement.read_at"
                    @click="markAsReadAndClose(selectedAnnouncement.id)"
                    class="rounded-xl bg-gradient-to-r from-blue-600 to-indigo-600 px-5 py-2.5 text-sm font-medium text-white shadow-lg shadow-blue-500/30 transition-all hover:shadow-xl hover:scale-105"
                  >
                    <span class="flex items-center gap-2">
                      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                      </svg>
                      {{ t('announcements.markRead') }}
                    </span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useAppStore } from '@/stores/app'
import { useAnnouncementStore } from '@/stores/announcements'
import { formatRelativeWithDateTime } from '@/utils/format'
import type { UserAnnouncement } from '@/types'
import Icon from '@/components/icons/Icon.vue'
import '@/styles/announcement-markdown.css'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const announcementStore = useAnnouncementStore()

// Configure marked
marked.setOptions({
  breaks: true,
  gfm: true,
})

// Use store state (storeToRefs for reactivity)
const { announcements, loading } = storeToRefs(announcementStore)
const unreadCount = computed(() => announcementStore.unreadCount)

// Local modal state
const isModalOpen = ref(false)
const activeTab = ref<'notifications' | 'announcement'>('notifications')
const qrLoadFailed = ref(false)
const qrImageUrl = '/announcement/qq-group.png'
const ANNOUNCEMENT_DISMISS_KEY = 'announcement_center_dismissed_date'

function localDateKey(): string {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
}

function shouldAutoOpen(): boolean {
  try {
    return localStorage.getItem(ANNOUNCEMENT_DISMISS_KEY) !== localDateKey()
  } catch {
    return true
  }
}


const detailModalOpen = ref(false)
const selectedAnnouncement = ref<UserAnnouncement | null>(null)

// Methods
function renderMarkdown(content: string): string {
  if (!content) return ''
  const html = marked.parse(content) as string
  return DOMPurify.sanitize(html)
}


function openModal() {
  isModalOpen.value = true
  // Always refresh when the notification center opens; the store otherwise
  // throttles list requests for 20 minutes and keeps deleted notices visible.
  void announcementStore.fetchAnnouncements(true)
}

function closeModal() {
  isModalOpen.value = false
}

function closeForToday() {
  try {
    localStorage.setItem(ANNOUNCEMENT_DISMISS_KEY, localDateKey())
  } catch {
    // ignore storage failures
  }
  closeModal()
}


function closeDetail() {
  detailModalOpen.value = false
  selectedAnnouncement.value = null
}

async function markAsRead(id: number) {
  try {
    const marked = await announcementStore.markAsRead(id)
    if (!marked) appStore.showError(t('common.unknownError'))
    return marked
  } catch (err: any) {
    appStore.showError(err?.message || t('common.unknownError'))
    return false
  }
}

async function markAsReadAndClose(id: number) {
  if (!await markAsRead(id)) return
  appStore.showSuccess(t('announcements.markedAsRead'))
  closeDetail()
}


// Announcement bodies may link to the group QR code with ?announcements=group.
// Inside the dialog that link must only switch the tab: navigating away would
// close the notification centre the reader is already looking at.
const ANNOUNCEMENT_QUERY_PARAM = 'announcements'
const GROUP_TAB_QUERY_VALUES = ['group', 'qr', 'announcement']

function announcementQueryValue(raw: unknown): string {
  const value = Array.isArray(raw) ? raw[0] : raw
  return typeof value === 'string' ? value.trim().toLowerCase() : ''
}

function handleMarkdownLinkClick(event: MouseEvent) {
  // Leave modified clicks and middle clicks to the browser/new tab.
  if (event.defaultPrevented || event.button !== 0) return
  if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return
  const anchor = (event.target as HTMLElement | null)?.closest?.('a[href]')
  if (!anchor) return
  const href = anchor.getAttribute('href')
  if (!href) return
  let url: URL
  try {
    url = new URL(href, window.location.origin)
  } catch {
    return
  }
  if (url.origin !== window.location.origin) return
  const value = announcementQueryValue(url.searchParams.get(ANNOUNCEMENT_QUERY_PARAM))
  if (!value) return
  event.preventDefault()
  activeTab.value = GROUP_TAB_QUERY_VALUES.includes(value) ? 'announcement' : 'notifications'
}

function consumeAnnouncementQuery(): boolean {
  const value = announcementQueryValue(route.query[ANNOUNCEMENT_QUERY_PARAM])
  if (!value || value === '0' || value === 'false') return false
  if (GROUP_TAB_QUERY_VALUES.includes(value)) {
    activeTab.value = 'announcement'
  }
  openModal()
  // Drop the parameter so closing the dialog and refreshing does not reopen it.
  const restQuery = { ...route.query }
  delete restQuery[ANNOUNCEMENT_QUERY_PARAM]
  void router.replace({ path: route.path, query: restQuery, hash: route.hash })
  return true
}

function handleEscape(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    if (detailModalOpen.value) {
      closeDetail()
    } else if (isModalOpen.value) {
      closeModal()
    }
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
  if (consumeAnnouncementQuery()) return
  if (shouldAutoOpen()) {
    isModalOpen.value = true
    void announcementStore.fetchAnnouncements(true)
  }
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleEscape)
  document.body.style.overflow = ''
})

// The header stays mounted across routes, so an in-app navigation to the deep
// link must be picked up here as well.
watch(
  () => route.query[ANNOUNCEMENT_QUERY_PARAM],
  () => {
    consumeAnnouncementQuery()
  }
)

watch(
  [isModalOpen, detailModalOpen, () => announcementStore.currentPopup],
  ([modal, detail, popup]) => {
    document.body.style.overflow = (modal || detail || popup) ? 'hidden' : ''
  }
)
</script>

<style scoped>
/* Modal Animations */
.modal-fade-enter-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.modal-fade-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 1, 1);
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-from > div {
  transform: scale(0.94) translateY(-12px);
  opacity: 0;
}

.modal-fade-leave-to > div {
  transform: scale(0.96) translateY(-8px);
  opacity: 0;
}

/* Scrollbar Styling */
.overflow-y-auto::-webkit-scrollbar {
  width: 8px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, #cbd5e1, #94a3b8);
  border-radius: 4px;
}

.dark .overflow-y-auto::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, #4b5563, #374151);
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, #94a3b8, #64748b);
}

.dark .overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, #6b7280, #4b5563);
}
</style>
