<template>
  <section class="flex min-h-[340px] flex-col overflow-hidden xl:min-h-0">
    <header class="flex h-9 shrink-0 items-center justify-between gap-2 border-b border-gray-200 px-3 dark:border-dark-700">
      <div class="flex min-w-0 items-center gap-2">
        <Icon v-if="icon" :name="icon" size="sm" class="shrink-0 text-gray-400 dark:text-dark-500" />
        <h2 class="truncate text-sm font-semibold text-gray-900 dark:text-white" :title="titleTooltip || undefined">{{ title }}</h2>
        <span
          v-if="badge !== null && badge !== undefined && badge !== ''"
          class="inline-flex shrink-0 items-center rounded-full bg-gray-100 px-1.5 py-0.5 text-[11px] font-medium text-gray-600 dark:bg-dark-700 dark:text-dark-300"
        >
          {{ badge }}
        </span>
      </div>
      <div class="flex shrink-0 items-center gap-1">
        <slot name="actions" />
      </div>
    </header>
    <div class="min-h-0 flex-1 overflow-y-auto p-3" :class="bodyClass">
      <slot />
    </div>
  </section>
</template>

<script setup lang="ts">
import Icon from '@/components/icons/Icon.vue'

/** Icon 组件的 name 联合类型（沿用组件自身的 props 定义，避免重复维护图标清单） */
type IconName = InstanceType<typeof Icon>['$props']['name']

withDefaults(
  defineProps<{
    title: string
    /** 悬停时展示的完整文本（标题被截断时用） */
    titleTooltip?: string
    icon?: IconName
    badge?: string | number | null
    bodyClass?: string
  }>(),
  {
    titleTooltip: undefined,
    icon: undefined,
    badge: null,
    bodyClass: '',
  },
)
</script>