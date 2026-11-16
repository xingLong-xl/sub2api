<template>
  <AppLayout>
    <div class="custom-page-layout">
      <div class="card flex-1 min-h-0 overflow-hidden">
        <div v-if="loading" class="flex h-full items-center justify-center py-12">
          <div
            class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"
          ></div>
        </div>

        <div
          v-else-if="!menuItem"
          class="flex h-full items-center justify-center p-10 text-center"
        >
          <div class="max-w-md">
            <div
              class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 dark:bg-dark-700"
            >
              <Icon name="link" size="lg" class="text-gray-400" />
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('customPage.notFoundTitle') }}
            </h3>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
              {{ t('customPage.notFoundDesc') }}
            </p>
          </div>
        </div>

        <!-- Markdown mode with TOC -->
        <div v-else-if="isMarkdownMode" class="flex h-full overflow-hidden">
          <!-- TOC Sidebar -->
          <aside class="toc-sidebar">
            <div class="toc-header">
              <span class="toc-title">{{ t('customPage.tableOfContents') }}</span>
            </div>
            <nav class="toc-nav">
              <a
                v-for="item in tocItems"
                :key="item.id"
                :href="'#' + item.id"
                class="toc-item"
                :class="[
                  `toc-level-${item.level}`,
                  { 'toc-active': activeHeadingId === item.id }
                ]"
                @click.prevent="scrollToHeading(item.id)"
              >
                {{ item.text }}
              </a>
            </nav>
          </aside>

          <!-- Content -->
          <div
            ref="markdownContainer"
            class="markdown-page-content flex-1 h-full overflow-auto p-6 md:p-10"
            v-html="renderedHtml"
            @scroll="onContentScroll"
          ></div>
        </div>

        <!-- URL not configured -->
        <div v-else-if="!isValidUrl" class="flex h-full items-center justify-center p-10 text-center">
          <div class="max-w-md">
            <div
              class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 dark:bg-dark-700"
            >
              <Icon name="link" size="lg" class="text-gray-400" />
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('customPage.notConfiguredTitle') }}
            </h3>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
              {{ t('customPage.notConfiguredDesc') }}
            </p>
          </div>
        </div>

        <!-- Iframe embed mode -->
        <div v-else ref="embedShell" class="custom-embed-shell">
          <a
            v-if="!menuItem?.hide_open_button"
            ref="openButton"
            :href="embeddedUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="btn btn-secondary btn-sm custom-open-fab"
            :style="openButtonPosition ? { left: `${openButtonPosition.x}px`, top: `${openButtonPosition.y}px`, right: 'auto' } : undefined"
            @pointerdown="startButtonDrag"
            @pointermove="moveButtonDrag"
            @pointerup="endButtonDrag"
            @pointercancel="endButtonDrag"
            @lostpointercapture="endButtonDrag"
            @dragstart.prevent
            @click="handleOpenButtonClick"
          >
            <Icon name="externalLink" size="sm" class="mr-1.5" :stroke-width="2" />
            {{ t('customPage.openInNewTab') }}
          </a>
          <iframe
            :src="embeddedUrl"
            class="custom-embed-frame"
            allowfullscreen
          ></iframe>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useResizeObserver } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import { useAdminSettingsStore } from '@/stores/adminSettings'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { buildApiUrl } from '@/api/client'
import { buildEmbeddedUrl, detectTheme } from '@/utils/embedded-url'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

interface TocItem {
  id: string
  text: string
  level: number
}

const { t, locale } = useI18n()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()
const adminSettingsStore = useAdminSettingsStore()

const loading = ref(false)
const pageTheme = ref<'light' | 'dark'>('light')
const renderedHtml = ref('')
const markdownContainer = ref<HTMLElement | null>(null)
const tocItems = ref<TocItem[]>([])
const activeHeadingId = ref('')
let themeObserver: MutationObserver | null = null
let markdownRequestVersion = 0

const embedShell = ref<HTMLElement | null>(null)
const openButton = ref<HTMLAnchorElement | null>(null)
const openButtonPosition = ref<{ x: number; y: number } | null>(null)
let buttonDrag: { pointerId: number; startX: number; startY: number; x: number; y: number } | null = null
let suppressButtonClick = false

function setOpenButtonPosition(x: number, y: number) {
  const shell = embedShell.value
  const button = openButton.value
  if (!shell || !button) return
  openButtonPosition.value = {
    x: Math.max(0, Math.min(x, shell.clientWidth - button.offsetWidth)),
    y: Math.max(0, Math.min(y, shell.clientHeight - button.offsetHeight)),
  }
}

function startButtonDrag(event: PointerEvent) {
  if (event.button !== 0 || !event.isPrimary || !openButton.value) return
  const button = openButton.value
  suppressButtonClick = false
  buttonDrag = {
    pointerId: event.pointerId, startX: event.clientX, startY: event.clientY,
    x: button.offsetLeft, y: button.offsetTop,
  }
  // Capture keeps receiving moves when the pointer crosses the embedded iframe.
  button.setPointerCapture(event.pointerId)
}

function moveButtonDrag(event: PointerEvent) {
  if (!buttonDrag || event.pointerId !== buttonDrag.pointerId) return
  const dx = event.clientX - buttonDrag.startX
  const dy = event.clientY - buttonDrag.startY
  if (!suppressButtonClick && Math.hypot(dx, dy) < 4) return
  suppressButtonClick = true
  setOpenButtonPosition(buttonDrag.x + dx, buttonDrag.y + dy)
  event.preventDefault()
}

function endButtonDrag(event: PointerEvent) {
  if (!buttonDrag || event.pointerId !== buttonDrag.pointerId) return
  buttonDrag = null
  if (openButton.value?.hasPointerCapture(event.pointerId)) {
    openButton.value.releasePointerCapture(event.pointerId)
  }
}

function handleOpenButtonClick(event: MouseEvent) {
  if (suppressButtonClick && event.detail !== 0) event.preventDefault()
  suppressButtonClick = false
}

useResizeObserver([embedShell, openButton], () => {
  if (openButtonPosition.value) {
    setOpenButtonPosition(openButtonPosition.value.x, openButtonPosition.value.y)
  }
})

const menuItemId = computed(() => route.params.id as string)

const menuItem = computed(() => {
  const id = menuItemId.value
  const publicItems = appStore.cachedPublicSettings?.custom_menu_items ?? []
  const found = publicItems.find((item) => item.id === id) ?? null
  if (found) return found
  if (authStore.isAdmin) {
    return adminSettingsStore.customMenuItems.find((item) => item.id === id) ?? null
  }
  return null
})

const markdownSlug = computed(() => {
  const item = menuItem.value
  if (!item) return ''
  if (item.page_slug) return item.page_slug
  if (item.url?.startsWith('md:')) return item.url.slice(3)
  return ''
})

const isMarkdownMode = computed(() => !!markdownSlug.value)

const embeddedUrl = computed(() => {
  if (!menuItem.value || isMarkdownMode.value) return ''
  return buildEmbeddedUrl(
    menuItem.value.url,
    authStore.user?.id,
    authStore.token,
    pageTheme.value,
    locale.value,
  )
})

const isValidUrl = computed(() => {
  if (isMarkdownMode.value) return false
  const url = embeddedUrl.value
  return url.startsWith('http://') || url.startsWith('https://')
})

function generateHeadingId(text: string, index: number): string {
  const base = text
    .toLowerCase()
    .replace(/[^\w一-鿿]+/g, '-')
    .replace(/^-+|-+$/g, '')
  return base ? `${base}-${index}` : `heading-${index}`
}

function isRelativeMarkdownAsset(src: string): boolean {
  const trimmed = src.trim()
  if (!trimmed || /^[a-z][a-z0-9+.-]*:/i.test(trimmed) || trimmed.startsWith('//') || trimmed.startsWith('/')) {
    return false
  }
  const [pathPart] = trimmed.split(/([?#].*)/, 2)
  return pathPart
    .split('/')
    .filter((part) => part && part !== '.')
    .every((part) => part !== '..' && !part.includes('\\'))
}

function buildPageImageUrl(slug: string, src: string): string {
  const trimmed = src.trim()
  const [pathPart, suffix = ''] = trimmed.split(/([?#].*)/, 2)
  const encodedPath = pathPart
    .split('/')
    .filter((part) => part && part !== '.')
    .map((part) => encodeURIComponent(part))
    .join('/')
  return buildApiUrl(`/pages/${encodeURIComponent(slug)}/images/${encodedPath}${suffix}`)
}

async function fetchAndRenderMarkdown(slug: string) {
  const version = ++markdownRequestVersion
  loading.value = true
  tocItems.value = []
  activeHeadingId.value = ''
  try {
    const resp = await fetch(buildApiUrl(`/pages/${encodeURIComponent(slug)}`), {
      headers: authStore.token ? { Authorization: `Bearer ${authStore.token}` } : {},
    })
    if (version !== markdownRequestVersion) return
    if (!resp.ok) {
      renderedHtml.value = `<p class="text-red-500">${t('common.pageNotFound')}</p>`
      return
    }
    let raw = await resp.text()
    if (version !== markdownRequestVersion) return

    raw = raw.replace(
      /!\[([^\]]*)\]\(([^)]+)\)/g,
      (match, alt, src) => isRelativeMarkdownAsset(src) ? `![${alt}](${buildPageImageUrl(slug, src)})` : match
    )

    const html = marked.parse(raw) as string
    const sanitized = DOMPurify.sanitize(html, {
      ADD_TAGS: ['iframe'],
      ADD_ATTR: ['allowfullscreen', 'frameborder', 'src', 'target', 'rel'],
    })

    // Inject IDs into headings and build TOC
    const toc: TocItem[] = []
    let headingIndex = 0
    const withIds = sanitized.replace(
      /<(h[1-4])[^>]*>(.*?)<\/h[1-4]>/gi,
      (_, tag: string, content: string) => {
        const level = parseInt(tag[1])
        const text = content.replace(/<[^>]+>/g, '').trim()
        const id = generateHeadingId(text, headingIndex++)
        toc.push({ id, text, level })
        return `<${tag} id="${id}">${content}</${tag}>`
      }
    )

    renderedHtml.value = withIds
    tocItems.value = toc
  } catch {
    if (version === markdownRequestVersion) {
      renderedHtml.value = '<p class="text-red-500">页面加载失败</p>'
    }
  } finally {
    if (version === markdownRequestVersion) {
      loading.value = false
      await nextTick()
      await nextTick()
      if (version === markdownRequestVersion) injectCopyButtons()
    }
  }
}

function scrollToHeading(id: string) {
  const container = markdownContainer.value
  if (!container) return
  const el = container.querySelector(`#${CSS.escape(id)}`)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    activeHeadingId.value = id
  }
}

let scrollRafId = 0
function onContentScroll() {
  if (scrollRafId) return
  scrollRafId = requestAnimationFrame(() => {
    scrollRafId = 0
    const container = markdownContainer.value
    if (!container || tocItems.value.length === 0) return

    const containerRect = container.getBoundingClientRect()
    let current = ''

    for (const item of tocItems.value) {
      const el = container.querySelector(`#${CSS.escape(item.id)}`) as HTMLElement | null
      if (el) {
        const elRect = el.getBoundingClientRect()
        if (elRect.top - containerRect.top <= 100) {
          current = item.id
        }
      }
    }
    activeHeadingId.value = current
  })
}

function injectCopyButtons() {
  const container = markdownContainer.value
  if (!container) return

  container.querySelectorAll('pre').forEach((pre) => {
    if (pre.querySelector('.copy-btn')) return
    const btn = document.createElement('button')
    btn.className = 'copy-btn'
    btn.textContent = t('customPage.copyCode')
    btn.addEventListener('click', async () => {
      const code = pre.querySelector('code')?.textContent ?? pre.textContent ?? ''
      try {
        await navigator.clipboard.writeText(code)
        btn.textContent = t('customPage.copiedCode')
        setTimeout(() => { btn.textContent = t('customPage.copyCode') }, 2000)
      } catch {
        btn.textContent = t('customPage.copyCodeFailed')
        setTimeout(() => { btn.textContent = t('customPage.copyCode') }, 2000)
      }
    })
    pre.style.position = 'relative'
    pre.appendChild(btn)
  })
}

watch(markdownSlug, (slug) => {
  if (slug) {
    fetchAndRenderMarkdown(slug)
  } else {
    markdownRequestVersion++
    loading.value = false
    renderedHtml.value = ''
    tocItems.value = []
  }
}, { immediate: true })

onMounted(async () => {
  pageTheme.value = detectTheme()

  if (typeof document !== 'undefined') {
    themeObserver = new MutationObserver(() => {
      pageTheme.value = detectTheme()
    })
    themeObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class'],
    })
  }

  if (appStore.publicSettingsLoaded) return
  loading.value = true
  try {
    await appStore.fetchPublicSettings()
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  markdownRequestVersion++
  if (themeObserver) {
    themeObserver.disconnect()
    themeObserver = null
  }
})
</script>

<style scoped>
.custom-page-layout {
  @apply flex flex-col;
  height: calc(100vh - 64px - 4rem);
}

.toc-sidebar {
  @apply flex flex-col h-full border-r border-gray-200/80 dark:border-dark-700/80;
  @apply bg-gray-50/40 dark:bg-dark-900/20;
  width: min(228px, 28%);
  min-width: 168px;
  max-width: 248px;
  overflow: hidden;
}

@media (max-width: 640px) {
  /* 小屏不展示目录，避免遮挡正文；不再提供折叠/展开交互。 */
  .toc-sidebar {
    display: none;
  }
}

.toc-header {
  @apply flex items-center justify-between px-5 pb-2 pt-5;
}

.toc-title {
  @apply text-[11px] font-semibold uppercase tracking-[0.14em] text-gray-400 dark:text-dark-500;
}

.toc-nav {
  @apply flex-1 space-y-0.5 overflow-y-auto px-3 py-2;
  scrollbar-gutter: stable;
}

.toc-item {
  @apply block truncate rounded-r-md border-l-2 border-transparent py-1.5 pl-3 pr-2 text-[13px] leading-5 transition-colors;
  @apply text-gray-500 hover:bg-gray-100/80 hover:text-gray-800 dark:text-dark-400 dark:hover:bg-dark-800/60 dark:hover:text-dark-100;
}

.toc-item.toc-active {
  @apply border-primary-500 bg-primary-50/70 font-medium text-primary-700;
  @apply dark:border-primary-400 dark:bg-primary-900/20 dark:text-primary-300;
}

.toc-level-1 { padding-left: 12px; }
.toc-level-2 { padding-left: 24px; }
.toc-level-3 { padding-left: 36px; }
.toc-level-4 { padding-left: 48px; }

.custom-embed-shell {
  @apply relative;
  @apply h-full w-full overflow-hidden rounded-2xl;
  @apply bg-gradient-to-b from-gray-50 to-white dark:from-dark-900 dark:to-dark-950;
  @apply p-0;
}

.custom-open-fab {
  @apply absolute right-3 top-3 z-10 w-max max-w-full touch-none select-none transition-colors;
  @apply shadow-sm backdrop-blur supports-[backdrop-filter]:bg-white/80 dark:supports-[backdrop-filter]:bg-dark-800/80;
}

.custom-embed-frame {
  display: block;
  margin: 0;
  width: 100%;
  height: 100%;
  border: 0;
  border-radius: 0;
  box-shadow: none;
  background: transparent;
}
</style>

<style>
.markdown-page-content {
  line-height: 1.7;
  color: inherit;
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
}

.markdown-page-content .guide-hero {
  margin: 0 0 2rem;
  padding: 1.75rem 1.75rem 1.5rem;
  border: 1px solid rgb(196 181 253 / 0.7);
  border-radius: 1rem;
  background:
    radial-gradient(circle at 100% 0%, rgb(237 233 254 / 0.9), transparent 46%),
    linear-gradient(135deg, rgb(245 243 255), rgb(255 255 255));
}

.dark .markdown-page-content .guide-hero {
  border-color: rgb(109 40 217 / 0.45);
  background:
    radial-gradient(circle at 100% 0%, rgb(76 29 149 / 0.35), transparent 46%),
    linear-gradient(135deg, rgb(30 27 75 / 0.7), rgb(17 24 39 / 0.35));
}

.markdown-page-content .guide-eyebrow {
  margin-bottom: 0.65rem;
  color: rgb(124 58 237);
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.16em;
}

.markdown-page-content .guide-hero h1 {
  margin: 0 0 0.7rem;
  padding: 0;
  border: 0;
  font-size: 1.85rem;
  line-height: 1.3;
}

.markdown-page-content .guide-hero p {
  margin: 0;
  max-width: 42rem;
  color: rgb(75 85 99);
  font-size: 0.95rem;
}

.dark .markdown-page-content .guide-hero p {
  color: rgb(209 213 219);
}

.markdown-page-content .guide-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.55rem;
  margin-top: 1.1rem;
}

.markdown-page-content .guide-meta span {
  padding: 0.3rem 0.6rem;
  border: 1px solid rgb(221 214 254);
  border-radius: 999px;
  color: rgb(109 40 217);
  font-size: 0.72rem;
}

.dark .markdown-page-content .guide-meta span {
  border-color: rgb(91 33 182 / 0.6);
  color: rgb(196 181 253);
}

.markdown-page-content h2 {
  padding-bottom: 0.45rem;
  border-bottom: 1px solid rgb(229 231 235);
}

.dark .markdown-page-content h2 {
  border-bottom-color: rgb(55 65 81);
}
.markdown-page-content h1 { @apply text-3xl font-bold mt-8 mb-4 pb-2 border-b border-gray-200 dark:border-dark-600; }
.markdown-page-content h2 { @apply text-2xl font-bold mt-6 mb-3; }
.markdown-page-content h3 { @apply text-xl font-semibold mt-5 mb-2; }
.markdown-page-content h4 { @apply text-lg font-semibold mt-4 mb-2; }
.markdown-page-content p { @apply mb-4; }
.markdown-page-content ul { @apply list-disc pl-6 mb-4; }
.markdown-page-content ol { @apply list-decimal pl-6 mb-4; }
.markdown-page-content li { @apply mb-1; }
.markdown-page-content a { @apply text-primary-500 hover:text-primary-600 underline; }
.markdown-page-content blockquote { @apply border-l-4 border-gray-300 dark:border-dark-500 pl-4 italic text-gray-600 dark:text-dark-300 my-4; }
.markdown-page-content img { @apply max-w-full h-auto rounded-lg my-4; }
.markdown-page-content table { @apply w-full border-collapse my-4; }
.markdown-page-content th { @apply border border-gray-300 dark:border-dark-500 px-3 py-2 bg-gray-50 dark:bg-dark-700 font-semibold text-left; }
.markdown-page-content td { @apply border border-gray-300 dark:border-dark-500 px-3 py-2; }
.markdown-page-content code { @apply bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-sm font-mono; }
.markdown-page-content pre { @apply bg-gray-100 dark:bg-dark-900/70 text-gray-800 dark:text-gray-200 border border-gray-200 dark:border-dark-700 p-4 rounded-lg overflow-x-auto my-4 relative; }
.markdown-page-content pre code { @apply bg-transparent p-0 text-inherit; }
.markdown-page-content hr { @apply my-6 border-gray-200 dark:border-dark-600; }

.copy-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 10px;
  font-size: 12px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.9);
  color: #374151;
  border: 1px solid rgba(0, 0, 0, 0.1);
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s, background 0.2s;
  font-family: inherit;
}
.copy-btn:hover { background: #ffffff; }
.dark .copy-btn {
  background: rgba(55, 65, 81, 0.9);
  color: #e5e7eb;
  border-color: rgba(255, 255, 255, 0.12);
}
.dark .copy-btn:hover { background: #374151; }
pre:hover .copy-btn { opacity: 1; }
</style>
