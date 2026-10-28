<template>
  <div class="relative flex min-h-screen flex-col bg-gray-50 dark:bg-dark-950">
    <!-- Header (same pattern as HomeView) -->
    <header class="sticky top-0 z-20 border-b border-gray-200/60 bg-white/80 px-6 py-3 backdrop-blur dark:border-dark-800/60 dark:bg-dark-950/80">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <router-link to="/home" class="flex items-center gap-3">
          <div class="h-9 w-9 overflow-hidden rounded-xl shadow-md">
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-lg font-semibold tracking-tight text-gray-900 dark:text-white">{{ siteName }}</span>
        </router-link>
        <div class="flex items-center gap-2">
          <button
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Content -->
    <main class="mx-auto flex w-full max-w-6xl flex-1 gap-8 px-4 py-8 sm:px-6">
      <!-- TOC -->
      <aside v-if="tocItems.length > 0" class="hidden w-56 shrink-0 lg:block">
        <div class="sticky top-24 max-h-[calc(100vh-8rem)] overflow-y-auto">
          <p class="mb-3 text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-dark-500">
            {{ t('customPage.tableOfContents') }}
          </p>
          <nav class="flex flex-col gap-0.5 border-l border-gray-200 dark:border-dark-700">
            <a
              v-for="item in tocItems"
              :key="item.id"
              :href="'#' + item.id"
              class="docs-toc-item block truncate border-l-2 border-transparent py-1 pr-2 text-sm text-gray-600 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white"
              :class="[
                `docs-toc-level-${item.level}`,
                { 'docs-toc-active': activeHeadingId === item.id },
              ]"
              @click.prevent="scrollToHeading(item.id)"
            >
              {{ item.text }}
            </a>
          </nav>
        </div>
      </aside>

      <!-- Markdown -->
      <article class="min-w-0 flex-1 rounded-2xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-800 dark:bg-dark-900 md:p-10">
        <div v-if="loading" class="flex items-center justify-center py-24">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
        </div>

        <div v-else-if="loadError" class="py-20 text-center">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('customPage.notFoundTitle') }}</h2>
          <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ t('docs.unavailableDesc') }}</p>
        </div>

        <div
          v-else
          ref="contentRef"
          class="docs-markdown"
          v-html="renderedHtml"
        ></div>
      </article>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import { buildApiUrl } from '@/api/client'
import { sanitizeUrl } from '@/utils/url'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

interface TocItem {
  id: string
  text: string
  level: number
}

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() =>
  sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', {
    allowRelative: true,
    allowDataUrl: true,
  }),
)

const loading = ref(true)
const loadError = ref(false)
const renderedHtml = ref('')
const contentRef = ref<HTMLElement | null>(null)
const tocItems = ref<TocItem[]>([])
const activeHeadingId = ref('')
let scrollRafId = 0

const isDark = ref(document.documentElement.classList.contains('dark'))
const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

function generateHeadingId(text: string, index: number): string {
  const base = text
    .toLowerCase()
    .replace(/[^\w一-鿿]+/g, '-')
    .replace(/^-+|-+$/g, '')
  return base ? `${base}-${index}` : `heading-${index}`
}

async function fetchDocs() {
  loading.value = true
  loadError.value = false
  try {
    const resp = await fetch(buildApiUrl('/docs'))
    if (!resp.ok) {
      loadError.value = true
      return
    }
    const raw = await resp.text()
    const html = marked.parse(raw) as string
    const sanitized = DOMPurify.sanitize(html, {
      ADD_TAGS: ['iframe'],
      ADD_ATTR: ['allowfullscreen', 'frameborder', 'src'],
    })

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
      },
    )

    renderedHtml.value = withIds
    tocItems.value = toc
  } catch {
    loadError.value = true
  } finally {
    loading.value = false
    await nextTick()
    injectCopyButtons()
  }
}

function scrollToHeading(id: string) {
  const container = contentRef.value
  if (!container) return
  const el = container.querySelector(`#${CSS.escape(id)}`)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    activeHeadingId.value = id
  }
}

function onWindowScroll() {
  if (scrollRafId) return
  scrollRafId = requestAnimationFrame(() => {
    scrollRafId = 0
    if (tocItems.value.length === 0) return

    let current = ''
    for (const item of tocItems.value) {
      const el = contentRef.value?.querySelector(`#${CSS.escape(item.id)}`) as HTMLElement | null
      if (el && el.getBoundingClientRect().top <= 120) {
        current = item.id
      }
    }
    activeHeadingId.value = current
  })
}

function injectCopyButtons() {
  const container = contentRef.value
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

onMounted(() => {
  initTheme()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  fetchDocs()
  window.addEventListener('scroll', onWindowScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', onWindowScroll)
  if (scrollRafId) cancelAnimationFrame(scrollRafId)
})
</script>

<style scoped>
.docs-toc-item {
  padding-left: 12px;
}
.docs-toc-level-1 { padding-left: 12px; }
.docs-toc-level-2 { padding-left: 24px; }
.docs-toc-level-3 { padding-left: 36px; }
.docs-toc-level-4 { padding-left: 48px; }

.docs-toc-active {
  @apply border-primary-500 font-medium text-primary-600 dark:text-primary-400;
}
</style>

<style>
.docs-markdown {
  line-height: 1.75;
  color: inherit;
  word-break: break-word;
}
.docs-markdown h1 { @apply text-3xl font-bold mt-8 mb-4 pb-2 border-b border-gray-200 dark:border-dark-600; }
.docs-markdown h1:first-child { @apply mt-0; }
.docs-markdown h2 { @apply text-2xl font-bold mt-8 mb-3 pb-2 border-b border-gray-100 dark:border-dark-700; }
.docs-markdown h3 { @apply text-xl font-semibold mt-6 mb-2; }
.docs-markdown h4 { @apply text-lg font-semibold mt-4 mb-2; }
.docs-markdown p { @apply mb-4; }
.docs-markdown ul { @apply list-disc pl-6 mb-4; }
.docs-markdown ol { @apply list-decimal pl-6 mb-4; }
.docs-markdown li { @apply mb-1.5; }
.docs-markdown a { @apply text-primary-500 hover:text-primary-600 underline; }
.docs-markdown blockquote { @apply border-l-4 border-primary-400 bg-primary-50/60 dark:bg-primary-900/20 dark:border-primary-600 px-4 py-2 text-gray-600 dark:text-dark-300 my-4 rounded-r; }
.docs-markdown img { @apply max-w-full h-auto rounded-lg my-4; }
.docs-markdown table { @apply w-full border-collapse my-4; }
.docs-markdown th { @apply border border-gray-300 dark:border-dark-500 px-3 py-2 bg-gray-50 dark:bg-dark-700 font-semibold text-left; }
.docs-markdown td { @apply border border-gray-300 dark:border-dark-500 px-3 py-2; }
.docs-markdown code { @apply bg-gray-100 dark:bg-dark-700 px-1.5 py-0.5 rounded text-sm font-mono; }
.docs-markdown pre { @apply bg-gray-900 dark:bg-dark-950 text-gray-100 p-4 rounded-lg overflow-x-auto my-4 relative; }
.docs-markdown pre code { @apply bg-transparent p-0 text-inherit; }
.docs-markdown hr { @apply my-8 border-gray-200 dark:border-dark-600; }

.docs-markdown .copy-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 10px;
  font-size: 12px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.15);
  color: #e2e8f0;
  border: 1px solid rgba(255, 255, 255, 0.2);
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s, background 0.2s;
  font-family: inherit;
}
.docs-markdown .copy-btn:hover { background: rgba(255, 255, 255, 0.25); }
.docs-markdown pre:hover .copy-btn { opacity: 1; }
</style>
