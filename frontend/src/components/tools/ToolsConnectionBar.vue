<template>
  <div
    :class="inline
      ? 'flex flex-wrap items-center gap-[8px]'
      : 'flex shrink-0 flex-wrap items-center gap-2 border-b border-gray-200 px-3 py-1.5 dark:border-dark-700'"
  >
    <select
      v-model="keyChoice"
      :class="['input !w-auto !rounded-lg !px-2.5 !py-1.5 text-sm', inline ? 'min-w-[150px] max-w-[200px]' : 'min-w-[180px] max-w-[280px]']"
      :aria-label="t('tools.connection.selectGroup')"
    >
      <option value="">{{ t('tools.connection.selectGroup') }}</option>
      <option v-for="option in keyOptions" :key="option.id" :value="option.id">
        {{ option.label }}
      </option>
    </select>

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

    <span v-if="!keysLoading && keyOptions.length === 0" class="text-xs text-amber-600 dark:text-amber-400">
      {{ t('tools.connection.noKeys') }}
    </span>

    <!-- 模型选择：右对齐，按工具能力过滤 -->
    <div :class="inline ? 'flex flex-wrap items-center gap-[8px]' : 'ml-auto flex flex-wrap items-center gap-2'">
      <span v-if="!inline" class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('tools.connection.model') }}</span>
      <select
        v-if="!hasGroup"
        disabled
        :class="['input !w-auto cursor-not-allowed !rounded-lg !px-2.5 !py-1.5 text-sm text-gray-400 dark:text-dark-500', inline ? 'min-w-[160px] max-w-[240px]' : 'min-w-[220px] max-w-[360px]']"
        :aria-label="t('tools.connection.model')"
      >
        <option>{{ t('tools.connection.selectGroupFirst') }}</option>
      </select>
      <select
        v-else-if="modelOptions.length > 0"
        v-model="model"
        :class="['input !w-auto !rounded-lg !px-2.5 !py-1.5 text-sm', inline ? 'min-w-[160px] max-w-[240px]' : 'min-w-[220px] max-w-[360px]']"
        :aria-label="t('tools.connection.model')"
      >
        <option value="">{{ t('tools.connection.selectModel') }}</option>
        <option v-for="item in modelOptions" :key="item" :value="item">{{ item }}</option>
      </select>
      <input
        v-else
        v-model="model"
        type="text"
        :class="['input !w-auto flex-1 !rounded-lg !px-2.5 !py-1.5 text-sm placeholder:text-[11px]', inline ? 'min-w-[200px] max-w-[320px]' : 'min-w-[300px] max-w-[440px]']"
        :placeholder="modelPlaceholder"
        :title="modelsError || undefined"
        autocomplete="off"
        spellcheck="false"
      />
      <button
        type="button"
        class="btn btn-secondary !px-2 !py-1.5"
        :disabled="modelsLoading || !hasGroup"
        :title="t('tools.connection.refreshModels')"
        @click="refreshModels"
      >
        <Icon name="refresh" size="sm" :class="{ 'animate-spin': modelsLoading }" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
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
    if (option) selectKey(option.key.id)
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

/** 没有模型清单时，把原因直接放进输入框占位符，而不是在框外再挂一行提示 */
/** 是否已经选中分组：未选中时模型框只提示、不接受输入 */
const hasGroup = computed(() => selectedGroupID.value !== null)

const modelPlaceholder = computed(() =>
  modelsError.value
    ? t('tools.connection.modelListUnavailable')
    : t('tools.connection.manualModelPlaceholder'),
)

const modelOptions = computed<string[]>(() => {
  const entry = selectedGroupID.value !== null ? groupCatalog.value[String(selectedGroupID.value)] : undefined
  let scoped: string[]
  if (entry && entry.models.length > 0) {
    // 渠道目录带 billing_mode，能力判定比模型名更准
    scoped = catalogModelsForGroup(selectedGroupID.value, props.tool)
  } else {
    // 目录未覆盖该分组时退回网关 /v1/models（余额/权限不足时会为空）
    scoped = modelsForCapability(props.tool, gatewayModels.value)
  }
  const current = props.modelValue.trim()
  if (current && !scoped.includes(current)) return [current, ...scoped]
  return scoped
})

// 换分组后模型清单会变，重新拉取
watch(selectedKeyId, async () => {
  if (hasApiKey.value) await fetchGatewayModels(true)
})

// 模型清单到达且当前为空时，默认选中第一个
watch(modelOptions, (options) => {
  if (!props.modelValue.trim() && options.length > 0) {
    emit('update:modelValue', options[0])
  }
})

async function refreshModels(): Promise<void> {
  await fetchGatewayModels(true)
}

onMounted(async () => {
  await Promise.all([loadUserKeys(), loadModelCatalog()])
  if (hasApiKey.value) await fetchGatewayModels()
})
</script>