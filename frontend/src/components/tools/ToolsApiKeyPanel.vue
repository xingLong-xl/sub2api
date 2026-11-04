<template>
  <div class="rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-start">
      <div class="flex min-w-0 flex-1 flex-col gap-1.5">
        <span class="flex items-center gap-1.5 text-sm font-medium text-gray-900 dark:text-white">
          <Icon name="key" size="sm" />
          {{ t('tools.apiKey.title') }}
        </span>
        <div class="flex items-center gap-2">
          <Input
            v-model="localKey"
            :type="showKey ? 'text' : 'password'"
            :placeholder="t('tools.apiKey.placeholder')"
            autocomplete="off"
            spellcheck="false"
            class="flex-1"
            @update:model-value="onKeyInput"
          />
          <button
            type="button"
            class="btn btn-secondary h-[38px] px-3"
            :title="showKey ? t('tools.apiKey.hide') : t('tools.apiKey.show')"
            @click="showKey = !showKey"
          >
            <Icon :name="showKey ? 'eyeOff' : 'eye'" size="sm" />
          </button>
          <button
            v-if="apiKey"
            type="button"
            class="btn btn-secondary h-[38px] px-3"
            :title="t('tools.apiKey.clear')"
            @click="clearKey"
          >
            <Icon name="x" size="sm" />
          </button>
        </div>
        <p class="text-xs leading-5 text-gray-500 dark:text-dark-400">
          {{ t('tools.apiKey.hint') }}
          <router-link to="/keys" class="font-medium text-primary-600 hover:underline dark:text-primary-400">
            {{ t('tools.apiKey.manage') }}
          </router-link>
        </p>
      </div>
      <div class="flex shrink-0 items-center gap-1.5 sm:pt-6">
        <span
          class="inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium"
          :class="hasApiKey
            ? 'bg-green-50 text-green-700 dark:bg-green-900/30 dark:text-green-400'
            : 'bg-amber-50 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'"
        >
          <span class="h-1.5 w-1.5 rounded-full" :class="hasApiKey ? 'bg-green-500' : 'bg-amber-500'"></span>
          {{ hasApiKey ? t('tools.apiKey.configured') : t('tools.apiKey.notConfigured') }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import Input from '@/components/common/Input.vue'
import { useToolsSettings } from '@/composables/useToolsSettings'

const { t } = useI18n()
const { apiKey, hasApiKey, setApiKey } = useToolsSettings()

const showKey = ref(false)
// 本地缓冲输入值，失焦式保存交给 composable 的 watch（trim 后落库）
const localKey = ref(apiKey.value)

function onKeyInput(value: string): void {
  setApiKey(value)
}

function clearKey(): void {
  setApiKey('')
  localKey.value = ''
}
</script>
