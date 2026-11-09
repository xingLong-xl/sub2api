<template>
  <AppLayout>
    <div class="-m-[11px] flex h-[calc(100dvh-77px)] flex-col overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800 md:-m-[19px] lg:-m-[27px]">
      <div
        class="grid min-h-0 flex-1 grid-cols-1 divide-y divide-gray-200 overflow-y-auto dark:divide-dark-700 xl:grid-cols-[380px_minmax(0,1fr)_300px] xl:divide-x xl:divide-y-0 xl:overflow-hidden"
      >
        <!-- 第一栏：生成参数 -->
        <ToolsPanel :title="t('tools.image.params')" icon="cog">
          <div class="space-y-3">
            <ToolsConnectionBar
              v-model="imageModel"
              tool="image"
              inline
              class="border-b border-gray-100 pb-3 dark:border-dark-700"
            />
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-image-prompt">
                {{ t('tools.image.prompt') }}
              </label>
              <textarea
                id="tools-image-prompt"
                v-model="prompt"
                rows="5"
                :placeholder="t('tools.image.promptPlaceholder')"
                class="input resize-none text-sm"
              ></textarea>
            </div>

            <div class="grid grid-cols-2 gap-2">
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-image-size">
                  {{ t('tools.image.size') }}
                </label>
                <select id="tools-image-size" v-model="size" class="input !py-1.5 text-sm">
                  <option value="auto">{{ t('tools.image.sizeAuto') }}</option>
                  <option value="1024x1024">1024 × 1024</option>
                  <option value="1024x1536">1024 × 1536</option>
                  <option value="1536x1024">1536 × 1024</option>
                  <option value="1024x1792">1024 × 1792</option>
                  <option value="1792x1024">1792 × 1024</option>
                </select>
              </div>
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-image-count">
                  {{ t('tools.image.count') }}
                </label>
                <input
                  id="tools-image-count"
                  v-model.number="count"
                  type="number"
                  min="1"
                  max="4"
                  class="input !py-1.5 text-sm"
                />
              </div>
            </div>

            <!-- 高级参数：后端支持但通常无需改动的字段 -->
            <div class="rounded-lg border border-gray-200 dark:border-dark-700">
              <button
                type="button"
                class="flex w-full items-center justify-between px-2.5 py-2 text-xs font-medium text-gray-600 dark:text-dark-300"
                @click="showAdvanced = !showAdvanced"
              >
                <span>{{ t('tools.image.advanced') }}</span>
                <Icon name="chevronDown" size="sm" :class="['transition-transform', showAdvanced && 'rotate-180']" />
              </button>
              <div v-if="showAdvanced" class="space-y-2.5 border-t border-gray-200 px-2.5 py-2.5 dark:border-dark-700">
                <div>
                  <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-quality">
                    {{ t('tools.image.quality') }}
                  </label>
                  <select id="tools-image-quality" v-model="quality" class="input !py-1.5 text-sm">
                    <option value="">{{ t('tools.image.qualityAuto') }}</option>
                    <option value="low">{{ t('tools.image.qualityLow') }}</option>
                    <option value="medium">{{ t('tools.image.qualityMedium') }}</option>
                    <option value="high">{{ t('tools.image.qualityHigh') }}</option>
                  </select>
                </div>
                <div>
                  <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-style">
                    {{ t('tools.image.style') }}
                  </label>
                  <select id="tools-image-style" v-model="style" class="input !py-1.5 text-sm">
                    <option value="">{{ t('tools.image.qualityAuto') }}</option>
                    <option value="natural">{{ t('tools.image.styleNatural') }}</option>
                    <option value="vivid">{{ t('tools.image.styleVivid') }}</option>
                  </select>
                </div>
                <div>
                  <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-format">
                    {{ t('tools.image.responseFormat') }}
                  </label>
                  <select id="tools-image-format" v-model="responseFormat" class="input !py-1.5 text-sm">
                    <option value="">{{ t('tools.image.qualityAuto') }}</option>
                    <option value="url">{{ t('tools.image.formatUrl') }}</option>
                    <option value="b64_json">{{ t('tools.image.formatB64') }}</option>
                  </select>
                </div>
                <div>
                  <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-background">
                    {{ t('tools.image.background') }}
                  </label>
                  <select id="tools-image-background" v-model="background" class="input !py-1.5 text-sm">
                    <option value="">{{ t('tools.image.backgroundAuto') }}</option>
                    <option value="transparent">{{ t('tools.image.backgroundTransparent') }}</option>
                    <option value="opaque">{{ t('tools.image.backgroundOpaque') }}</option>
                  </select>
                </div>
                <div>
                  <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-output">
                    {{ t('tools.image.outputFormat') }}
                  </label>
                  <select id="tools-image-output" v-model="outputFormat" class="input !py-1.5 text-sm">
                    <option value="">png</option>
                    <option value="png">png</option>
                    <option value="jpeg">jpeg</option>
                    <option value="webp">webp</option>
                  </select>
                </div>
                <div>
                  <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-moderation">
                    {{ t('tools.image.moderation') }}
                  </label>
                  <select id="tools-image-moderation" v-model="moderation" class="input !py-1.5 text-sm">
                    <option value="">{{ t('tools.image.moderationAuto') }}</option>
                    <option value="auto">{{ t('tools.image.moderationAuto') }}</option>
                    <option value="low">{{ t('tools.image.moderationLow') }}</option>
                  </select>
                </div>
              </div>
            </div>

            <button
              type="button"
              class="btn btn-primary w-full"
              :disabled="!canGenerate"
              @click="generate"
            >
              <Icon name="sparkles" size="sm" />
              {{ generating ? t('tools.image.generating') : t('tools.image.generate') }}
            </button>

            <p v-if="error" class="rounded-lg bg-red-50 px-2.5 py-2 text-xs text-red-600 dark:bg-red-900/20 dark:text-red-400">
              {{ error }}
            </p>
          </div>
        </ToolsPanel>

        <!-- 第二栏：生成结果 -->
        <ToolsPanel
          :title="t('tools.image.results')"
          icon="grid"
          :badge="activeRecord ? activeRecord.images.length : null"
        >
          <template #actions>
            <button
              v-if="activeRecord"
              type="button"
              class="btn btn-secondary !px-2 !py-1 text-xs"
              @click="activeRecord = null"
            >
              {{ t('tools.common.clear') }}
            </button>
          </template>

          <div v-if="generating" class="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <div v-for="n in count" :key="n" class="h-56 animate-pulse rounded-lg bg-gray-100 dark:bg-dark-700"></div>
          </div>

          <div v-else-if="activeRecord && activeRecord.images.length" class="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <figure
              v-for="(image, index) in activeRecord.images"
              :key="index"
              class="group relative overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-dark-700 dark:bg-dark-900"
            >
              <img :src="image.src" :alt="activeRecord.prompt" class="h-auto w-full object-contain" />
              <button
                type="button"
                class="absolute right-2 top-2 rounded-md bg-black/60 p-1.5 text-white opacity-0 transition-opacity group-hover:opacity-100"
                :title="t('tools.common.download')"
                @click="downloadImage(image, activeRecord.model, index)"
              >
                <Icon name="download" size="sm" />
              </button>
            </figure>
          </div>

          <div v-else class="flex h-full min-h-[240px] flex-col items-center justify-center gap-2 text-gray-400 dark:text-dark-500">
            <Icon name="sparkles" size="xl" />
            <p class="text-sm">{{ t('tools.image.empty') }}</p>
          </div>
        </ToolsPanel>

        <!-- 第三栏：历史记录（暂为会话内） -->
        <ToolsPanel :title="t('tools.common.history')" icon="clock" :badge="history.length">
          <p class="mb-2 text-xs text-gray-400 dark:text-dark-500">{{ t('tools.image.historyHint') }}</p>
          <p v-if="history.length === 0" class="py-6 text-center text-xs text-gray-400 dark:text-dark-500">
            {{ t('tools.common.emptyHistory') }}
          </p>
          <div v-else class="space-y-2">
            <button
              v-for="(record, index) in history"
              :key="index"
              type="button"
              class="w-full rounded-lg border border-gray-200 p-2 text-left transition-colors hover:border-primary-300 hover:bg-primary-50/40 dark:border-dark-700 dark:hover:border-primary-700 dark:hover:bg-primary-900/10"
              :class="record === activeRecord ? 'border-primary-400 bg-primary-50/60 dark:border-primary-600 dark:bg-primary-900/20' : ''"
              @click="activeRecord = record"
            >
              <div class="flex items-center justify-between gap-2 text-[11px] text-gray-400 dark:text-dark-500">
                <span>{{ formatTime(record.createdAt) }}</span>
                <span class="truncate">{{ record.model }}</span>
              </div>
              <p class="mt-1 truncate text-xs text-gray-700 dark:text-dark-200">{{ record.prompt }}</p>
              <div class="mt-1.5 flex items-center gap-1.5">
                <img
                  v-for="(image, imageIndex) in record.images.slice(0, 4)"
                  :key="imageIndex"
                  :src="image.src"
                  class="h-9 w-9 rounded object-cover"
                  alt=""
                />
                <span
                  class="ml-auto inline-flex items-center gap-1 text-[11px] text-primary-600 dark:text-primary-400"
                  @click.stop="reuse(record)"
                >
                  <Icon name="refresh" size="xs" />
                  {{ t('tools.common.reuse') }}
                </span>
              </div>
            </button>
          </div>
        </ToolsPanel>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ToolsConnectionBar from '@/components/tools/ToolsConnectionBar.vue'
import ToolsPanel from '@/components/tools/ToolsPanel.vue'
import { GatewayError, useToolsSettings } from '@/composables/useToolsSettings'

interface GeneratedImage {
  /** data: URL（b64_json）或上游 URL */
  src: string
  /** 是否为 data URL，决定下载方式 */
  isDataUrl: boolean
}

interface HistoryParams {
  size: string
  count: number
  quality: string
  style: string
  responseFormat: string
  background: string
  outputFormat: string
  moderation: string
}

interface GenerationRecord {
  prompt: string
  model: string
  images: GeneratedImage[]
  params: HistoryParams
  createdAt: Date
}

const { t } = useI18n()
const { hasApiKey, imageModel, gatewayFetch } = useToolsSettings()

const prompt = ref('')
const size = ref('auto')
const count = ref(1)
const quality = ref('')
const style = ref('')
const responseFormat = ref('')
const background = ref('')
const outputFormat = ref('')
const moderation = ref('')
const showAdvanced = ref(false)

const generating = ref(false)
const error = ref('')

const history = ref<GenerationRecord[]>([])
const activeRecord = ref<GenerationRecord | null>(null)

const canGenerate = computed(
  () => !generating.value && prompt.value.trim() !== '' && imageModel.value.trim() !== '' && hasApiKey.value,
)

function localizedError(err: unknown): string {
  if (err instanceof GatewayError) {
    if (err.code === 'INSUFFICIENT_BALANCE') return t('tools.errors.insufficientBalance')
    if (err.status === 401) return t('tools.errors.unauthorized')
    if (err.status === 403) return t('tools.errors.forbidden')
    if (err.status === 429) return t('tools.errors.rateLimited')
    return err.message
  }
  return err instanceof Error ? err.message : t('tools.errors.unknown')
}

function formatTime(date: Date): string {
  return date.toLocaleTimeString()
}

/** 解析 /v1/images/generations 响应，兼容 b64_json 与 url 两种返回 */
function parseImages(body: unknown): GeneratedImage[] {
  const data = (body as { data?: unknown })?.data
  if (!Array.isArray(data)) return []
  const images: GeneratedImage[] = []
  for (const item of data) {
    const b64 = (item as { b64_json?: unknown })?.b64_json
    const url = (item as { url?: unknown })?.url
    if (typeof b64 === 'string' && b64) {
      images.push({ src: `data:image/png;base64,${b64}`, isDataUrl: true })
    } else if (typeof url === 'string' && url) {
      images.push({ src: url, isDataUrl: false })
    }
  }
  return images
}

function currentParams(): HistoryParams {
  return {
    size: size.value,
    count: count.value,
    quality: quality.value,
    style: style.value,
    responseFormat: responseFormat.value,
    background: background.value,
    outputFormat: outputFormat.value,
    moderation: moderation.value,
  }
}

async function generate(): Promise<void> {
  if (!canGenerate.value) return
  generating.value = true
  error.value = ''
  try {
    const payload: Record<string, unknown> = {
      model: imageModel.value.trim(),
      prompt: prompt.value.trim(),
      size: size.value,
      n: count.value,
    }
    // 空字符串表示「用上游默认值」，不发送
    if (quality.value) payload.quality = quality.value
    if (style.value) payload.style = style.value
    if (responseFormat.value) payload.response_format = responseFormat.value
    if (background.value) payload.background = background.value
    if (outputFormat.value) payload.output_format = outputFormat.value
    if (moderation.value) payload.moderation = moderation.value

    const res = await gatewayFetch('/v1/images/generations', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    const body: unknown = await res.json()
    const images = parseImages(body)
    if (images.length === 0) {
      error.value = t('tools.image.noResults')
      return
    }
    const record: GenerationRecord = {
      prompt: prompt.value.trim(),
      model: imageModel.value.trim(),
      images,
      params: currentParams(),
      createdAt: new Date(),
    }
    history.value.unshift(record)
    activeRecord.value = record
  } catch (err) {
    error.value = localizedError(err)
  } finally {
    generating.value = false
  }
}

/** 把历史记录里的参数回填到表单，方便微调后重试 */
function reuse(record: GenerationRecord): void {
  prompt.value = record.prompt
  imageModel.value = record.model
  size.value = record.params.size
  count.value = record.params.count
  quality.value = record.params.quality
  style.value = record.params.style
  responseFormat.value = record.params.responseFormat
  background.value = record.params.background
  outputFormat.value = record.params.outputFormat
  moderation.value = record.params.moderation
  activeRecord.value = record
}

function downloadImage(image: GeneratedImage, model: string, index: number): void {
  const link = document.createElement('a')
  const safeModel = model.replace(/[^a-zA-Z0-9._-]/g, '-') || 'image'
  if (image.isDataUrl) {
    link.href = image.src
    link.download = `${safeModel}-${Date.now()}-${index + 1}.png`
  } else {
    // 上游 URL 无法保证同源下载，退化为新标签页打开
    link.href = image.src
    link.target = '_blank'
    link.rel = 'noopener'
  }
  link.click()
}
</script>