<template>
  <div :class="rootClass">
    <div>
      <select
        id="tools-connection-group"
        v-model="keyChoice"
        :class="groupSelectClass"
        :aria-label="t('tools.connection.selectGroup')"
      >
        <option value="">{{ t('tools.connection.selectGroup') }}</option>
        <option v-for="option in keyOptions" :key="option.id" :value="option.id">
          {{ option.label }}
        </option>
      </select>
    </div>

    <div :class="inline ? 'flex min-w-0 flex-1 flex-wrap items-center gap-[8px]' : 'ml-auto flex flex-wrap items-center gap-2'">
      <span v-if="!inline" class="text-xs font-medium text-gray-500 dark:text-dark-400">
        {{ t('tools.connection.model') }}
      </span>
      <div class="flex items-center gap-1.5" :class="inline ? 'min-w-0 flex-1' : ''">
        <select
          v-if="!hasGroup"
          disabled
          :class="[modelSelectClass, 'cursor-not-allowed text-gray-400 dark:text-dark-500']"
          :aria-label="t('tools.connection.model')"
        >
          <option>{{ t('tools.connection.selectGroupFirst') }}</option>
        </select>
        <select
          v-else-if="modelsLoading"
          disabled
          :class="[modelSelectClass, 'cursor-not-allowed text-gray-400 dark:text-dark-500']"
          :aria-label="t('tools.connection.model')"
        >
          <option>{{ t('tools.common.loading') }}</option>
        </select>
        <select
          v-else-if="modelOptions.length > 0"
          v-model="model"
          :class="modelSelectClass"
          :aria-label="t('tools.connection.model')"
        >
          <option value="">{{ t('tools.connection.selectModel') }}</option>
          <option v-for="item in modelOptions" :key="item" :value="item">{{ item }}</option>
        </select>
        <select
          v-else
          disabled
          :class="[modelSelectClass, 'cursor-not-allowed text-gray-400 dark:text-dark-500']"
          :aria-label="t('tools.connection.model')"
        >
          <option>{{ t('tools.connection.noModels') }}</option>
        </select>
      </div>
    </div>

    <p
      v-if="hasGroup && !modelsLoading && modelOptions.length === 0 && !capabilityUnsupported"
      class="w-full text-[11px] leading-4 text-amber-600 dark:text-amber-400"
    >
      {{ modelsError ? t('tools.connection.modelListUnavailable') : t('tools.connection.noModelsHint') }}
    </p>

    <p
      v-if="capabilityUnsupported"
      class="flex w-full flex-wrap items-center gap-1 text-xs leading-5 text-amber-600 dark:text-amber-400"
    >
      <Icon name="exclamationTriangle" size="xs" class="shrink-0" />
      <span>{{ t('tools.connection.groupUnsupported', { group: currentGroupLabel }) }}</span>
      <template v-if="suggestedGroups.length">
        <span class="text-amber-500/70">·</span>
        <span>{{ t('tools.connection.supportedGroups') }}:</span>
        <button
          v-for="group in suggestedGroups"
          :key="group.id"
          type="button"
          class="rounded bg-amber-100 px-1.5 py-0.5 text-amber-700 transition-colors hover:bg-amber-200 dark:bg-amber-900/30 dark:text-amber-300 dark:hover:bg-amber-900/50"
          @click="switchGroup(group.id)"
        >
          {{ group.label }}<span v-if="group.models.length" class="opacity-75">（{{ group.models.slice(0, 3).join(' / ') }}）</span>
        </button>
      </template>
      <template v-else>
        <span class="text-amber-500/70">·</span>
        <span>{{ t('tools.connection.noSupportedGroup') }}</span>
      </template>
    </p>

    <span v-if="keysError" class="text-xs text-red-500 dark:text-red-400">
      {{ keysError }}
    </span>
    <span v-else-if="!keysLoading && keyOptions.length === 0" class="text-xs text-amber-600 dark:text-amber-400">
      {{ t('tools.connection.noKeys') }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useToolsSettings } from '@/composables/useToolsSettings'
import type { ToolCapability } from '@/composables/useToolsSettings'
import type { ApiKey } from '@/types'

const props = defineProps<{
  tool: ToolCapability
  modelValue: string
  /** 内联模式：去掉分割线与外边距，用于嵌在输入区内部 */
  inline?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const { t } = useI18n()
const {
  apiKey,
  userKeys,
  keysLoading,
  keysError,
  selectedKey,
  selectedKeyId,
  selectedGroupID,
  groupCatalog,
  loadModelCatalog,
  catalogModelsForGroup,
  loadUserKeys,
  selectKey,
  clearSelection,
  hasApiKey,
  gatewayModels,
  modelsLoading,
  modelsError,
  fetchGatewayModels,
  modelsForCapability,
  platformSupports,
  platformSupportsTool,
} = useToolsSettings()

const rootClass = computed(() => {
  if (props.inline) return 'flex flex-wrap items-center gap-[8px]'
  return 'flex shrink-0 flex-wrap items-center gap-2 border-b border-gray-200 px-3 py-1.5 dark:border-dark-700'
})

const groupSelectClass = computed(() => [
  'input !rounded-lg !px-2.5 !py-1.5 text-sm',
  props.inline
    ? '!w-auto min-w-[110px] max-w-[145px]'
    : '!w-auto min-w-[180px] max-w-[280px]',
])

const modelSelectClass = computed(() => [
  'input !rounded-lg !px-2.5 !py-1.5 text-sm',
  props.inline
    ? '!w-full min-w-0'
    : '!w-auto min-w-[220px] max-w-[360px]',
])


interface KeyGroupOption {
  /** 选项值：分组优先（g:<id>），无分组时退化为单把密钥（k:<id>） */
  id: string
  label: string
  key: ApiKey
}

/** 同一分组只出现一次：工具页按「分组」选择，内部自动使用该分组下的一把可用密钥 */
const keyOptions = computed<KeyGroupOption[]>(() => {
  const options: KeyGroupOption[] = []
  const seen = new Set<string>()
  for (const key of userKeys.value) {
    const group = key.group
    const id = group ? `g:${group.id}` : `k:${key.id}`
    if (seen.has(id)) continue
    seen.add(id)
    const groupName = group?.name?.trim()
    const keyName = key.name?.trim()
    const label = groupName || keyName || key.key.slice(0, 10)
    options.push({ id, label, key })
  }
  return options
})

const keyChoice = computed({
  get: () => {
    const key = selectedKey.value ?? userKeys.value.find((item) => item.key === apiKey.value)
    if (!key) return ''
    const value = key.group ? `g:${key.group.id}` : `k:${key.id}`
    // 选项尚未加载完成时回退到占位项，避免下拉框显示为空白
    return keyOptions.value.some((option) => option.id === value) ? value : ''
  },
  set: (value: string) => {
    // 选回「选择分组」时清空：不再携带密钥，模型也一并清掉
    if (!value) {
      clearSelection()
      emit('update:modelValue', '')
      return
    }
    const option = keyOptions.value.find((item) => item.id === value)
    if (option) {
      selectKey(option.key.id)
      // 换分组后模型列表会整体变化，清空旧模型，等用户重新选择
      emit('update:modelValue', '')
    }
  },
})

const model = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value),
})

/** 当前分组平台不支持该工具时给出明确提示，并列出可用的分组 */
const capabilityUnsupported = computed(() => platformSupports(props.tool) === false)

const currentGroupLabel = computed(() => selectedKey.value?.group?.name?.trim() || t('tools.connection.ungrouped'))

/** 用户其它可用分组中，平台支持当前工具的那些（附带该分组下的同类模型） */
const suggestedGroups = computed(() =>
  keyOptions.value
    .filter((option) => platformSupportsTool(option.key.group?.platform ?? '', props.tool) === true)
    .map((option) => {
      const groupID = option.key.group?.id
      const entry = groupID !== undefined ? groupCatalog.value[String(groupID)] : undefined
      const models = (entry?.models ?? [])
        .filter((model) => model.capability === props.tool)
        .map((model) => model.name)
      return { id: option.id, label: option.label, models }
    }),
)

function switchGroup(id: string): void {
  keyChoice.value = id
}

/** 是否已经选中分组：未选中时模型框只提示、不接受输入 */
const hasGroup = computed(() => selectedGroupID.value !== null)

const modelOptions = computed<string[]>(() => {
  const entry = selectedGroupID.value !== null ? groupCatalog.value[String(selectedGroupID.value)] : undefined
  if (entry && entry.models.length > 0) {
    // 渠道目录带 billing_mode，能力判定比模型名更准
    return catalogModelsForGroup(selectedGroupID.value, props.tool)
  }
  // 目录未覆盖该分组时退回网关 /v1/models（余额/权限不足时会为空）
  return modelsForCapability(props.tool, gatewayModels.value)
})

// 换分组后模型清单会变，重新拉取
watch(selectedKeyId, async () => {
  if (hasApiKey.value) await fetchGatewayModels(true)
})

onMounted(async () => {
  // 强制刷新：另一个浏览器/标签页新建密钥后，当前页也能在进入工具页时看到。
  await Promise.all([loadUserKeys(true), loadModelCatalog(true)])
  if (hasApiKey.value) await fetchGatewayModels()
  window.addEventListener('focus', refreshKeysOnVisible)
  document.addEventListener('visibilitychange', refreshKeysOnVisible)
})

onBeforeUnmount(() => {
  window.removeEventListener('focus', refreshKeysOnVisible)
  document.removeEventListener('visibilitychange', refreshKeysOnVisible)
})

function refreshKeysOnVisible(): void {
  if (document.visibilityState === 'visible') void loadUserKeys(true)
}
</script>