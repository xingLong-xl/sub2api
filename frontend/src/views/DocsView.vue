<template>
  <div
    class="docs-shell relative min-h-screen bg-slate-50 text-gray-900 dark:bg-dark-950 dark:text-gray-100"
  >
    <div class="docs-backdrop pointer-events-none absolute inset-x-0 top-0 h-[460px]" aria-hidden="true"></div>

    <header
      class="sticky top-0 z-30 border-b border-gray-200/70 bg-white/80 backdrop-blur-xl dark:border-dark-800/80 dark:bg-dark-950/80"
    >
      <div class="mx-auto flex h-16 w-full max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <router-link to="/home" class="flex min-w-0 items-center gap-2.5">
          <span
            class="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900"
          >
            <img :src="siteLogo || '/logo.svg'" alt="" class="h-full w-full object-contain" />
          </span>
          <span class="truncate text-[15px] font-semibold tracking-tight text-gray-900 dark:text-white">
            {{ siteName }}
          </span>
          <span
            class="hidden rounded-full border border-gray-200 bg-gray-50 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-gray-500 sm:inline dark:border-dark-700 dark:bg-dark-900 dark:text-dark-400"
          >
            Docs
          </span>
        </router-link>

        <div class="flex items-center gap-1">
          <router-link
            to="/model-plaza"
            class="hidden items-center gap-1.5 rounded-lg px-3 py-2 text-[13px] font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 sm:flex dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white"
          >
            <Icon name="grid" size="sm" />
            {{ t('nav.modelPlaza') }}
          </router-link>
          <button
            type="button"
            class="inline-flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? '切换到浅色模式' : '切换到深色模式'"
            :aria-label="isDark ? '切换到浅色模式' : '切换到深色模式'"
            @click="toggleTheme"
          >
            <Icon :name="isDark ? 'sun' : 'moon'" size="sm" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="ml-1 inline-flex items-center rounded-lg bg-gray-900 px-3.5 py-2 text-[13px] font-medium text-white shadow-sm transition-colors hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </div>
    </header>

    <main class="relative mx-auto w-full max-w-7xl px-4 pb-16 pt-10 sm:px-6 sm:pt-14 lg:px-8">
      <section class="mx-auto max-w-3xl text-center">
        <span
          class="inline-flex items-center gap-2 rounded-full border border-gray-200/80 bg-white/80 px-3 py-1 text-[11px] font-medium text-gray-500 shadow-sm backdrop-blur dark:border-dark-700/80 dark:bg-dark-900/60 dark:text-dark-300"
        >
          <span class="relative flex h-1.5 w-1.5">
            <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-60"></span>
            <span class="relative inline-flex h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
          </span>
          {{ siteName }} Docs
        </span>

        <h1
          class="mt-5 text-[32px] font-bold leading-tight tracking-tight text-gray-900 sm:text-[40px] dark:text-white"
        >
          配置文档
        </h1>
        <p class="mx-auto mt-4 max-w-2xl text-[15px] leading-7 text-gray-600 dark:text-dark-300">
          从一个 API Key 开始：一个 Base URL，兼容 OpenAI、Anthropic、Gemini 三种协议。按下方步骤完成注册、创建密钥、选择模型，即可在几分钟内接入 Claude Code、Codex CLI 或你自己的应用。
        </p>

        <div class="mt-7 flex flex-wrap items-center justify-center gap-3">
          <router-link
            to="/keys"
            class="inline-flex items-center gap-2 rounded-xl bg-gray-900 px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:-translate-y-0.5 hover:bg-gray-800 hover:shadow-md dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
          >
            创建 API Key
            <Icon name="arrowRight" size="sm" />
          </router-link>
          <router-link
            to="/model-plaza"
            class="inline-flex items-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 shadow-sm transition-colors hover:border-gray-300 hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:bg-dark-800"
          >
            查看模型与价格
          </router-link>
        </div>

        <div class="mt-6 flex flex-wrap items-center justify-center gap-2">
          <span class="docs-chip">OpenAI 兼容</span>
          <span class="docs-chip">Anthropic 兼容</span>
          <span class="docs-chip">Gemini 兼容</span>
          <span class="docs-chip">流式输出</span>
          <span class="docs-chip">按量计费</span>
        </div>
      </section>

      <div
        v-if="loading"
        class="mx-auto mt-14 max-w-3xl rounded-3xl border border-gray-200/70 bg-white/70 p-8 shadow-sm sm:p-10 dark:border-dark-800 dark:bg-dark-900/60"
      >
        <div class="docs-skeleton h-5 w-1/3"></div>
        <div class="docs-skeleton mt-6 h-3 w-full"></div>
        <div class="docs-skeleton mt-3 h-3 w-11/12"></div>
        <div class="docs-skeleton mt-3 h-3 w-4/5"></div>
        <div class="docs-skeleton mt-8 h-4 w-1/4"></div>
        <div class="docs-skeleton mt-4 h-3 w-full"></div>
        <div class="docs-skeleton mt-3 h-3 w-10/12"></div>
      </div>

      <div
        v-else-if="loadError"
        class="mx-auto mt-14 max-w-xl rounded-3xl border border-gray-200 bg-white p-8 text-center shadow-sm dark:border-dark-800 dark:bg-dark-900"
      >
        <div
          class="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-gray-100 text-gray-400 dark:bg-dark-800 dark:text-dark-400"
        >
          <Icon name="book" size="lg" />
        </div>
        <h2 class="mt-4 text-lg font-semibold">{{ t('customPage.notFoundTitle') }}</h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ t('docs.unavailableDesc') }}</p>
        <router-link
          to="/"
          class="mt-6 inline-flex items-center rounded-lg border border-gray-200 bg-white px-3.5 py-2 text-[13px] font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700"
        >
          返回首页
        </router-link>
      </div>

      <div v-else class="mx-auto mt-12 max-w-[880px] lg:mt-16">
        <article
          class="docs-card rounded-3xl border border-gray-200/70 bg-white/90 p-6 shadow-[0_1px_2px_rgba(16,24,40,0.04),0_16px_48px_-24px_rgba(16,24,40,0.18)] backdrop-blur sm:p-8 md:p-10 dark:border-dark-800 dark:bg-dark-900/80"
        >
          <div ref="contentRef" class="docs-markdown" v-html="renderedHtml"></div>
        </article>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAppStore, useAuthStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import { buildApiUrl } from '@/api/client'
import { sanitizeUrl } from '@/utils/url'
import { enhanceDocsOsTabs } from '@/utils/docsOsTabs'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const { t } = useI18n()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'TokenPilot')
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
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  isDark.value = savedTheme ? savedTheme === 'dark' : prefersDark
  document.documentElement.classList.toggle('dark', isDark.value)
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
    const raw = (await resp.text()).trim()
    if (!raw) {
      loadError.value = true
      return
    }

    const html = marked.parse(raw, { gfm: true, breaks: false }) as string
    const sanitized = DOMPurify.sanitize(html, {
      ADD_TAGS: ['iframe'],
      ADD_ATTR: ['allowfullscreen', 'frameborder', 'src'],
    })

    let headingIndex = 0
    const withIds = sanitized.replace(
      /<(h[1-4])[^>]*>(.*?)<\/h[1-4]>/gi,
      (_, tag: string, content: string) => {
        const text = content.replace(/<[^>]+>/g, '').trim()
        const id = generateHeadingId(text, headingIndex++)
        return `<${tag} id="${id}">${content}</${tag}>`
      },
    )

    renderedHtml.value = withIds
  } catch {
    loadError.value = true
  } finally {
    loading.value = false
    await nextTick()
    enhanceContent()
    scrollToHashOnLoad()
  }
}

function scrollToHashOnLoad() {
  const id = decodeURIComponent(window.location.hash.replace(/^#/, ''))
  if (!id) return
  const el = contentRef.value?.querySelector(`#${CSS.escape(id)}`) as HTMLElement | null
  if (el) {
    requestAnimationFrame(() => {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    })
  }
}

function enhanceContent() {
  const container = contentRef.value
  if (!container) return

  container.querySelectorAll('table').forEach((table) => {
    if (table.closest('.table-wrap')) return
    const wrapper = document.createElement('div')
    wrapper.className = 'table-wrap'
    table.parentNode?.insertBefore(wrapper, table)
    wrapper.appendChild(table)
  })

  container.querySelectorAll('pre').forEach((pre) => {
    if (pre.closest('.code-block')) return
    const code = pre.querySelector('code')
    const lang = (code?.className.match(/language-([\w-]+)/i)?.[1] || '').toLowerCase()

    const wrapper = document.createElement('div')
    wrapper.className = 'code-block'
    pre.parentNode?.insertBefore(wrapper, pre)
    wrapper.appendChild(pre)

    const bar = document.createElement('div')
    bar.className = 'code-bar'
    const label = document.createElement('span')
    label.className = 'code-lang'
    label.textContent = lang || 'code'
    const btn = document.createElement('button')
    btn.type = 'button'
    btn.className = 'copy-btn'
    btn.textContent = t('customPage.copyCode')
    btn.addEventListener('click', async () => {
      const codeText = code?.textContent ?? pre.textContent ?? ''
      try {
        await navigator.clipboard.writeText(codeText)
        btn.textContent = t('customPage.copiedCode')
        setTimeout(() => {
          btn.textContent = t('customPage.copyCode')
        }, 2000)
      } catch {
        btn.textContent = t('customPage.copyCodeFailed')
        setTimeout(() => {
          btn.textContent = t('customPage.copyCode')
        }, 2000)
      }
    })
    bar.append(label, btn)
    wrapper.insertBefore(bar, pre)
  })

  enhanceDocsOsTabs(container)

  container.querySelectorAll('h2, h3, h4').forEach((heading) => {
    const id = heading.getAttribute('id')
    if (!id || heading.querySelector('.heading-anchor')) return
    const anchor = document.createElement('a')
    anchor.className = 'heading-anchor'
    anchor.href = `#${id}`
    anchor.setAttribute('aria-label', '复制此节链接')
    anchor.textContent = '#'
    heading.appendChild(anchor)
  })

  container.querySelectorAll('a[href^="http"]').forEach((link) => {
    const href = link.getAttribute('href') || ''
    if (href.startsWith(window.location.origin)) return
    link.setAttribute('target', '_blank')
    link.setAttribute('rel', 'noopener noreferrer')
  })

  container.querySelectorAll<HTMLAnchorElement>('.docs-cta a[href^="/"]').forEach((link) => {
    link.addEventListener('click', (event) => {
      if (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return
      event.preventDefault()
      router.push(link.getAttribute('href') || '/')
    })
  })
}

onMounted(() => {
  initTheme()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  fetchDocs()
})
</script>

<style>
.docs-backdrop {
  background:
    radial-gradient(52% 52% at 50% -10%, rgba(104, 43, 208, 0.14), transparent 70%),
    radial-gradient(36% 36% at 88% 4%, rgba(14, 165, 233, 0.1), transparent 70%);
}

.dark .docs-backdrop {
  background:
    radial-gradient(52% 52% at 50% -10%, rgba(104, 43, 208, 0.2), transparent 70%),
    radial-gradient(36% 36% at 88% 4%, rgba(14, 165, 233, 0.12), transparent 70%);
}

.docs-chip {
  @apply inline-flex items-center gap-1.5 rounded-full border border-gray-200/80 bg-white/70 px-2.5 py-1 text-[11px] font-medium text-gray-600 shadow-sm backdrop-blur dark:border-dark-700/80 dark:bg-dark-900/50 dark:text-dark-300;
}

.docs-skeleton {
  @apply animate-pulse rounded-md bg-gray-200/70 dark:bg-dark-800;
}

.docs-markdown {
  @apply text-[15px] leading-[1.9] text-gray-700 dark:text-dark-200;
  word-break: break-word;
  overflow-wrap: anywhere;
}

.docs-markdown > :first-child {
  margin-top: 0;
}

.docs-markdown p {
  @apply my-4;
}

.docs-markdown a {
  @apply font-medium text-primary-600 underline decoration-primary-300/70 underline-offset-2 transition-colors hover:text-primary-700 hover:decoration-primary-500 dark:text-primary-400 dark:decoration-primary-500/50 dark:hover:text-primary-300;
}

.docs-markdown .docs-cta {
  @apply mt-5 flex flex-wrap items-center gap-2.5;
}

.docs-markdown .docs-cta-btn {
  @apply inline-flex h-9 items-center justify-center rounded-lg border border-gray-200 bg-white px-4 text-[13px] font-medium text-gray-700 no-underline shadow-sm transition-colors hover:border-gray-300 hover:bg-gray-50 hover:text-gray-900 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700 dark:hover:text-white;
}

.docs-markdown .docs-cta-btn--primary {
  @apply border-gray-900 bg-gray-900 text-white hover:border-gray-800 hover:bg-gray-800 hover:text-white dark:border-white dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100 dark:hover:text-gray-900;
}

.docs-markdown strong {
  @apply font-semibold text-gray-900 dark:text-white;
}

.docs-markdown ul {
  @apply my-4 list-disc space-y-1.5 pl-6 marker:text-gray-300 dark:marker:text-dark-600;
}

.docs-markdown ol {
  @apply my-4 list-decimal space-y-1.5 pl-6 marker:text-gray-400 dark:marker:text-dark-500;
}

.docs-markdown li > ul,
.docs-markdown li > ol {
  @apply my-2;
}

.docs-markdown li > p {
  @apply my-2;
}

.docs-markdown h2 {
  @apply mb-4 mt-12 scroll-mt-24 border-t border-gray-100 pt-8 text-[22px] font-bold tracking-tight text-gray-900 dark:border-dark-800 dark:text-white;
}

.docs-markdown h2:first-of-type {
  @apply mt-0 border-t-0 pt-0;
}

.docs-markdown h3 {
  @apply mb-3 mt-8 scroll-mt-24 text-[17px] font-semibold tracking-tight text-gray-900 dark:text-white;
}

.docs-markdown h4 {
  @apply mb-2 mt-6 scroll-mt-24 text-[15px] font-semibold text-gray-900 dark:text-white;
}

.docs-markdown h5 {
  @apply mb-2 mt-6 flex items-center gap-2.5 text-[14px] font-semibold text-gray-900 dark:text-white;
}

.docs-markdown .os-panel {
  counter-reset: docs-step;
}

.docs-markdown .os-panel h5 {
  counter-increment: docs-step;
}

.docs-markdown .os-panel h5::before {
  content: counter(docs-step);
  @apply flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary-50 text-[11px] font-bold text-primary-600 dark:bg-primary-900/30 dark:text-primary-300;
}

.docs-markdown .os-tabs {
  @apply my-5;
}

.docs-markdown .os-tablist {
  @apply inline-flex items-center gap-0.5 rounded-lg border border-gray-200 bg-gray-100/80 p-0.5 dark:border-dark-700 dark:bg-dark-800/60;
}

.docs-markdown .os-tablist--master {
  @apply mb-3;
}

.docs-markdown .os-tablist--master + h3 {
  @apply mt-2;
}

.docs-markdown .os-tablist--master + .os-tabs {
  @apply mt-0;
}

.docs-markdown .os-tablist[hidden] {
  display: none;
}

.docs-markdown .os-tab {
  @apply min-w-[82px] justify-center rounded-md px-5 py-1 text-center text-[12.5px] font-medium leading-5 text-gray-600 transition-colors hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500/40 dark:text-dark-300 dark:hover:text-white;
}

.docs-markdown .os-tab[aria-selected='true'] {
  @apply bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white;
}

.docs-markdown .os-panel {
  @apply pt-3;
}

.docs-markdown .os-panel[hidden] {
  display: none;
}

.docs-markdown .heading-anchor {
  @apply ml-2 text-gray-300 no-underline opacity-0 transition-opacity duration-150 dark:text-dark-600;
}

.docs-markdown :is(h2, h3, h4):hover .heading-anchor {
  @apply opacity-100;
}

.docs-markdown blockquote {
  @apply my-5 rounded-r-xl border-l-[3px] border-primary-400 bg-primary-50/70 px-4 py-3 text-gray-600 dark:border-primary-500 dark:bg-primary-900/20 dark:text-dark-300;
}

.docs-markdown blockquote p {
  @apply my-0;
}

.docs-markdown blockquote > :not(:first-child) {
  @apply mt-2;
}

.docs-markdown code {
  @apply rounded-md border border-gray-200 bg-gray-50 px-1.5 py-0.5 font-mono text-[0.86em] text-gray-800 dark:border-dark-700 dark:bg-dark-800 dark:text-gray-200;
}

.docs-markdown pre code {
  @apply border-0 bg-transparent p-0 text-inherit;
}

.docs-markdown .code-block {
  @apply relative my-5 overflow-hidden rounded-xl border border-slate-800/80 bg-[#0b1220] shadow-sm;
}

.docs-markdown .code-bar {
  @apply flex items-center justify-between border-b border-white/[0.06] bg-white/[0.03] px-3 py-2;
}

.docs-markdown .code-lang {
  @apply font-mono text-[11px] font-semibold uppercase tracking-[0.14em] text-slate-400;
}

.docs-markdown .code-block pre {
  @apply m-0 overflow-x-auto bg-transparent p-4 text-[12.5px] leading-6 text-slate-100;
  tab-size: 2;
}

.docs-markdown .copy-btn {
  @apply rounded-md border border-white/10 bg-white/5 px-2 py-1 text-[11px] font-medium text-slate-300 transition-colors hover:bg-white/10 hover:text-white;
}

.docs-markdown .table-wrap {
  @apply my-5 overflow-x-auto rounded-xl border border-gray-200 dark:border-dark-700;
}

.docs-markdown table {
  @apply w-full min-w-[560px] border-collapse text-left text-[13.5px];
}

.docs-markdown thead th {
  @apply border-b border-gray-200 bg-gray-50/80 px-4 py-2.5 text-[12px] font-semibold text-gray-500 dark:border-dark-700 dark:bg-dark-800/60 dark:text-dark-300;
}

.docs-markdown tbody td {
  @apply border-b border-gray-100 px-4 py-2.5 align-top text-gray-600 dark:border-dark-800 dark:text-dark-300;
}

.docs-markdown tbody tr:last-child td {
  @apply border-b-0;
}

.docs-markdown tbody tr:hover {
  @apply bg-gray-50/70 dark:bg-dark-800/40;
}

.docs-markdown tbody td:first-child {
  @apply font-medium text-gray-800 dark:text-gray-100;
}

.docs-markdown hr {
  @apply my-10 border-gray-100 dark:border-dark-800;
}

.docs-markdown img {
  @apply my-5 max-w-full rounded-xl border border-gray-200 dark:border-dark-700;
}
</style>
