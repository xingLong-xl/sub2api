<template>
  <AppLayout>
    <div class="space-y-4">
      <!-- 密钥面板 -->
      <ToolsApiKeyPanel />

      <div class="grid gap-4 lg:grid-cols-3">
        <!-- 参数表单 -->
        <div class="space-y-4 lg:col-span-1">
          <div class="space-y-4 rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800">
            <!-- 模型 -->
            <div>
              <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-dark-200" for="tools-image-model">
                {{ t('tools.image.model') }}
              </label>
              <input
                id="tools-image-model"
                v-model="imageModel"
                list="tools-image-models"
                type="text"
                :placeholder="t('tools.image.modelPlaceholder')"
                class="input h-[38px]"
                autocomplete="off"
                spellcheck="false"
              />
              <datalist id="tools-image-models">
                <option v-for="m in availableModels" :key="m" :value="m" />
              </datalist>
            </div>

            <!-- 提示词 -->
            <div>
              <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-dark-200" for="tools-image-prompt">
                {{ t('tools.image.prompt') }}
              </label>
              <textarea
                id="tools-image-prompt"
                v-model="prompt"
                rows="5"
                :placeholder="t('tools.image.promptPlaceholder')"
                class="input resize-none"
              ></textarea>
            </div>

            <!-- 尺寸 -->
            <div>
              <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-dark-200" for="tools-image-size">
                {{ t('tools.image.size') }}
              </label>
              <select id="tools-image-size" v-model="size" class="input h-[38px]">
                <option value="auto">{{ t('tools.image.sizeAuto') }}</option>
                <option value="1024x1024">1024 × 1024</option>
                <option value="1024x1536">1024 × 1536</option>
                <option value="1536x1024">1536 × 1024</option>
                <option value="1024x1792">1024 × 1792</option>
                <option value="1792x1024">1792 × 1024</option>
              </select>
            </div>

            <!-- 数量 -->
            <div>
              <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-dark-200" for="tools-image-count">
                {{ t('tools.image.count') }}
              </label>
              <select id="tools-image-count" v-model.number="count" class="input h-[38px]">
                <option :value="1">1</option>
                <option :value="2">2</option>
                <option :value="3">3</option>
                <option :value="4">4</option>
              </select>
            </div>

            <button
              type="button"
              class="btn btn-primary w-full"
              :disabled="!canGenerate || generating"
              @click="generate"
            >
              <Icon v-if="generating" name="sync" size="sm" class="animate-spin" />
              {{ generating ? t('tools.image.generating') : t('tools.image.generate') }}
            </button>
          </div>

          <!-- 错误提示 -->
          <div
            v-if="error"
            class="rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-900/20 dark:text-red-400"
          >
            {{ error }}
          </div>
        </div>

        <!-- 结果区 -->
        <div class="lg:col-span-2">
          <div class="rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800">
            <div class="mb-3 flex items-center justify-between">
              <span class="text-sm font-medium text-gray-900 dark:text-white">
                {{ t('tools.image.results') }}
              </span>
              <button
                v-if="history.length"
                type="button"
                class="btn btn-secondary h-[32px] text-xs"
                @click="history = []"
              >
                <Icon name="trash" size="xs" />
                {{ t('tools.image.clear') }}
              </button>
            </div>

            <div v-if="history.length === 0" class="flex h-72 flex-col items-center justify-center gap-2 text-gray-400 dark:text-dark-500">
              <Icon name="sparkles" size="xl" />
              <p class="text-sm">{{ t('tools.image.empty') }}</p>
            </div>

            <div v-else class="space-y-6">
              <div v-for="(item, index) in history" :key="index" class="space-y-2">
                <div class="flex items-center justify-between gap-2 text-xs text-gray-500 dark:text-dark-400">
                  <span class="truncate">{{ item.prompt }}</span>
                  <span class="shrink-0">{{ item.model }} · {{ formatTime(item.createdAt) }}</span>
                </div>
                <div class="grid grid-cols-2 gap-3 sm:grid-cols-3">
                  <div
                    v-for="(image, imgIndex) in item.images"
                    :key="imgIndex"
                    class="group relative overflow-hidden rounded-lg border border-gray-200 dark:border-dark-700"
                  >
                    <img :src="image.src" :alt="`generated-${index}-${imgIndex}`" class="aspect-square w-full object-cover" />
                    <button
                      type="button"
                      class="absolute bottom-2 right-2 hidden rounded-md bg-black/60 px-2 py-1 text-white backdrop-blur group-hover:block"
                      @click="downloadImage(image, item.model, imgIndex)"
                    >
                      <Icon name="download" size="xs" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ToolsApiKeyPanel from '@/components/tools/ToolsApiKeyPanel.vue'
import { GatewayError, useToolsSettings } from '@/composables/useToolsSettings'

interface GeneratedImage {
  /** data: URL（b64_json）或上游 URL */
  src: string
  /** 是否为 data URL，决定下载方式 */
  isDataUrl: boolean
}

interface GenerationRecord {
  prompt: string
  model: string
  images: GeneratedImage[]
  createdAt: Date
}

const { t } = useI18n()
const { hasApiKey, imageModel, gatewayFetch, fetchGatewayModels } = useToolsSettings()

const prompt = ref('')
const size = ref('auto')
const count = ref(1)
const generating = ref(false)
const error = ref('')

const availableModels = ref<string[]>([])

const history = ref<GenerationRecord[]>([])

const canGenerate = computed(() => !generating.value && prompt.value.trim() !== '' && imageModel.value.trim() !== '' && hasApiKey.value)

function localizedError(err: unknown): string {
  if (err instanceof GatewayError) {
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

async function generate(): Promise<void> {
  if (!canGenerate.value) return
  generating.value = true
  error.value = ''
  try {
    const res = await gatewayFetch('/v1/images/generations', {
      method: 'POST',
      body: JSON.stringify({
        model: imageModel.value.trim(),
        prompt: prompt.value.trim(),
        size: size.value,
        n: count.value,
      }),
    })
    const body: unknown = await res.json()
    const images = parseImages(body)
    if (images.length === 0) {
      error.value = t('tools.image.noResults')
      return
    }
    history.value.unshift({
      prompt: prompt.value.trim(),
      model: imageModel.value.trim(),
      images,
      createdAt: new Date(),
    })
  } catch (err) {
    error.value = localizedError(err)
  } finally {
    generating.value = false
  }
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

onMounted(async () => {
  if (hasApiKey.value) {
    availableModels.value = await fetchGatewayModels()
  }
})
</script>
