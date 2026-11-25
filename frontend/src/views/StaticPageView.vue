<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white">
    <header class="border-b border-gray-200/70 bg-white/80 backdrop-blur dark:border-dark-800 dark:bg-dark-900/80">
      <div class="mx-auto flex max-w-5xl items-center justify-between px-6 py-4">
        <router-link to="/home" class="flex items-center gap-3">
          <img
            :src="siteLogo || '/logo.svg'"
            alt="Logo"
            class="h-9 w-9 rounded-lg object-contain"
          />
          <span class="text-base font-semibold tracking-tight">{{ siteName }}</span>
        </router-link>
        <router-link
          to="/home"
          class="rounded-lg px-3 py-2 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
        >
          返回首页
        </router-link>
      </div>
    </header>

    <main class="mx-auto max-w-4xl px-6 py-12 sm:py-16">
      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      </div>

      <div
        v-else-if="error"
        class="rounded-2xl border border-red-200 bg-red-50 p-6 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-300"
      >
        {{ error }}
      </div>

      <article
        v-else
        class="markdown-body prose prose-sm max-w-none dark:prose-invert"
        v-html="rendered"
      ></article>
    </main>

    <footer class="border-t border-gray-200/70 px-6 py-6 dark:border-dark-800">
      <p class="mx-auto max-w-5xl text-center text-sm text-gray-500 dark:text-dark-400">
        &copy; {{ currentYear }} {{ siteName }}. 保留所有权利。
      </p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useAppStore } from '@/stores/app'
import '@/styles/announcement-markdown.css'

const route = useRoute()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'TokenPilot')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const currentYear = new Date().getFullYear()

const pages: Record<string, { file: string; title: string }> = {
  '/about': { file: 'about.md', title: '关于我们' },
  '/privacy-policy': { file: 'privacy-policy.md', title: '隐私政策' },
  '/user-agreement': { file: 'user-agreement.md', title: '用户协议' },
}

const page = computed(() => {
  const path = route.path.replace(/\/+$/, '') || '/about'
  return pages[path] ?? pages['/about']
})

const loading = ref(true)
const error = ref('')
const rendered = ref('')

async function loadPage(): Promise<void> {
  loading.value = true
  error.value = ''
  rendered.value = ''
  try {
    const response = await fetch(`/pages/${page.value.file}`, { cache: 'no-store' })
    if (!response.ok) throw new Error(`页面加载失败（${response.status}）`)
    const markdown = await response.text()
    rendered.value = DOMPurify.sanitize(marked.parse(markdown) as string)
    document.title = `${page.value.title} - ${siteName.value}`
  } catch (err) {
    console.error('Failed to load static page', err)
    error.value = '页面暂时无法加载，请稍后重试。'
  } finally {
    loading.value = false
  }
}

watch(() => route.path, loadPage, { immediate: true })
</script>
