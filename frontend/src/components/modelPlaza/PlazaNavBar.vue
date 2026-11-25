<template>
  <header
    class="glass sticky top-0 z-30 border-b border-gray-200/50 dark:border-dark-700/50"
  >
    <div class="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-3.5 sm:px-6">
      <!-- 左:站点 logo + 名称 -->
      <RouterLink to="/home" class="group flex min-w-0 items-center gap-3" aria-label="返回首页">
        <template v-if="settings">
          <span
            class="flex h-9 w-9 flex-shrink-0 items-center justify-center overflow-hidden rounded-xl bg-white shadow-sm ring-1 ring-gray-200 transition-opacity group-hover:opacity-80 dark:bg-dark-800 dark:ring-dark-700"
          >
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </span>
          <span class="truncate text-base font-semibold text-gray-950 transition-colors group-hover:text-primary-600 dark:text-white dark:group-hover:text-primary-400">
            {{ siteName }}
          </span>
        </template>
        <template v-else>
          <span class="h-9 w-9 flex-shrink-0 animate-pulse rounded-xl bg-gray-200 dark:bg-dark-700" aria-hidden="true"></span>
          <span class="h-5 w-28 animate-pulse rounded bg-gray-200 dark:bg-dark-700" aria-hidden="true"></span>
        </template>
      </RouterLink>

      <!-- 右: 固定导航按钮 -->
      <div class="flex flex-shrink-0 items-center gap-1">
        <RouterLink to="/docs" class="inline-flex items-center justify-center whitespace-nowrap rounded-lg px-3 py-2 text-[13px] font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white">
          {{ t('nav.docs') }}
        </RouterLink>
        <RouterLink to="/model-plaza" class="inline-flex items-center justify-center whitespace-nowrap rounded-lg px-3 py-2 text-[13px] font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white">
          {{ t('nav.modelPlaza') }}
        </RouterLink>
        <RouterLink
          :to="isAuthenticated ? backTarget : { path: '/login', query: { redirect: '/model-plaza' } }"
          class="inline-flex items-center justify-center whitespace-nowrap rounded-lg px-3 py-2 text-[13px] font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white"
        >
          {{ t('nav.console') }}
        </RouterLink>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { sanitizeUrl } from '@/utils/url'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const settings = computed(() => appStore.cachedPublicSettings)
const siteName = computed(() => settings.value?.site_name || 'Sub2API')
const siteLogo = computed(() =>
  sanitizeUrl(settings.value?.site_logo || '', { allowRelative: true, allowDataUrl: true })
)
const isAuthenticated = computed(() => authStore.isAuthenticated)
const backTarget = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))
</script>
