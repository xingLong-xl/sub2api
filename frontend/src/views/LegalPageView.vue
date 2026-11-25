<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white">
    <header class="border-b border-gray-200/70 bg-white/80 backdrop-blur dark:border-dark-800 dark:bg-dark-900/80">
      <div class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <router-link to="/home" class="flex items-center gap-3">
          <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-9 w-9 rounded-lg object-contain" />
          <span class="text-base font-semibold tracking-tight">{{ siteName }}</span>
        </router-link>
        <div class="flex items-center gap-4 text-sm">
          <router-link to="/about" class="text-gray-500 hover:text-primary-600 dark:text-dark-400 dark:hover:text-primary-400">关于我们</router-link>
          <router-link to="/home" class="text-gray-500 hover:text-primary-600 dark:text-dark-400 dark:hover:text-primary-400">返回首页</router-link>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-6xl px-6 py-12 sm:py-16">
      <div class="mx-auto max-w-3xl text-center">
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-primary-500">{{ current.title }}</span>
        <h1 class="mt-4 text-3xl font-bold tracking-tight sm:text-4xl">{{ current.title }}</h1>
        <p class="mt-4 text-sm leading-7 text-gray-500 dark:text-dark-400">
          {{ current.description }}
        </p>
        <div class="mt-6 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-xs text-gray-400 dark:text-dark-500">
          <span>最后更新：{{ current.updated }}</span>
          <span>生效日期：{{ current.effective }}</span>
          <span>版本：{{ current.version }}</span>
        </div>
      </div>

      <div class="mt-10 grid gap-8 lg:grid-cols-[220px_minmax(0,1fr)]">
        <aside class="lg:sticky lg:top-8 lg:self-start">
          <p class="mb-3 text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-dark-500">目录</p>
          <nav class="flex gap-2 overflow-x-auto pb-2 lg:block lg:space-y-1 lg:overflow-visible lg:pb-0">
            <button
              v-for="(item, index) in toc"
              :key="item"
              class="whitespace-nowrap rounded-lg px-3 py-2 text-left text-sm text-gray-500 transition-colors hover:bg-white hover:text-primary-600 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-primary-400 lg:block lg:w-full"
              @click="scrollToSection(index)"
            >
              {{ item }}
            </button>
          </nav>
        </aside>

        <article
          class="legal-content markdown-body prose prose-sm max-w-none rounded-2xl border border-gray-200/70 bg-white p-6 shadow-sm dark:border-dark-800 dark:bg-dark-900 sm:p-10 dark:prose-invert"
          v-html="rendered"
        ></article>
      </div>
    </main>

    <footer class="border-t border-gray-200/70 px-6 py-6 dark:border-dark-800">
      <div class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-3 text-sm text-gray-500 dark:text-dark-400 sm:flex-row">
        <p>&copy; {{ currentYear }} {{ siteName }}. 保留所有权利。</p>
        <div class="flex items-center gap-4">
          <router-link to="/about" class="hover:text-primary-600 dark:hover:text-primary-400">关于我们</router-link>
          <router-link to="/privacy-policy" class="hover:text-primary-600 dark:hover:text-primary-400">隐私政策</router-link>
          <router-link to="/user-agreement" class="hover:text-primary-600 dark:hover:text-primary-400">用户协议</router-link>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useAppStore } from '@/stores/app'
import { privacyPolicy, userAgreement } from '@/content/legal'
import '@/styles/announcement-markdown.css'

const route = useRoute()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'TokenPilot')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const currentYear = new Date().getFullYear()

const current = computed(() => {
  if (route.path.startsWith('/user-agreement')) {
    return {
      title: '用户协议',
      description: '使用 TokenPilot 服务前，请阅读以下条款。',
      updated: '2026 年 6 月 8 日',
      effective: '2026 年 6 月 8 日',
      version: 'v1.2',
      content: userAgreement,
    }
  }
  return {
    title: '隐私政策',
    description: '了解 TokenPilot 如何收集、使用、存储和保护您的信息。',
    updated: '2026 年 6 月 8 日',
    effective: '2026 年 6 月 8 日',
    version: 'v1.1',
    content: privacyPolicy,
  }
})

const rendered = computed(() => {
  const body = current.value.content.replace(/^# .+\n+/, '')
  return DOMPurify.sanitize(marked.parse(body) as string)
})

const toc = computed(() =>
  [...current.value.content.matchAll(/^## (.+)$/gm)].map((match) => match[1].trim()),
)

function scrollToSection(index: number): void {
  const headings = document.querySelectorAll<HTMLElement>('.legal-content h2')
  headings[index]?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

watch(
  () => route.path,
  () => window.scrollTo({ top: 0, behavior: 'instant' }),
)
</script>
