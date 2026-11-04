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
              <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-dark-200" for="tools-video-model">
                {{ t('tools.video.model') }}
              </label>
              <input
                id="tools-video-model"
                v-model="videoModel"
                list="tools-video-models"
                type="text"
                :placeholder="t('tools.video.modelPlaceholder')"
                class="input h-[38px]"
                autocomplete="off"
                spellcheck="false"
              />
              <datalist id="tools-video-models">
                <option v-for="m in availableModels" :key="m" :value="m" />
              </datalist>
            </div>

            <!-- 提示词 -->
            <div>
              <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-dark-200" for="tools-video-prompt">
                {{ t('tools.video.prompt') }}
              </label>
              <textarea
                id="tools-video-prompt"
                v-model="prompt"
                rows="5"
                :placeholder="t('tools.video.promptPlaceholder')"
                class="input resize-none"
              ></textarea>
            </div>

            <button
              type="button"
              class="btn btn-primary w-full"
              :disabled="!canGenerate || submitting"
              @click="submitTask"
            >
              <Icon v-if="submitting" name="sync" size="sm" class="animate-spin" />
              {{ submitting ? t('tools.video.submitting') : t('tools.video.generate') }}
            </button>

            <p class="text-xs leading-5 text-gray-400 dark:text-dark-500">
              {{ t('tools.video.asyncHint') }}
            </p>
          </div>

          <!-- 错误提示 -->
          <div
            v-if="error"
            class="rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-900/20 dark:text-red-400"
          >
            {{ error }}
          </div>
        </div>

        <!-- 任务列表 -->
        <div class="lg:col-span-2">
          <div class="rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800">
            <div class="mb-3 flex items-center justify-between">
              <span class="text-sm font-medium text-gray-900 dark:text-white">
                {{ t('tools.video.tasks') }}
              </span>
              <button
                v-if="tasks.length"
                type="button"
                class="btn btn-secondary h-[32px] text-xs"
                @click="clearTasks"
              >
                <Icon name="trash" size="xs" />
                {{ t('tools.video.clear') }}
              </button>
            </div>

            <div v-if="tasks.length === 0" class="flex h-72 flex-col items-center justify-center gap-2 text-gray-400 dark:text-dark-500">
              <Icon name="play" size="xl" />
              <p class="text-sm">{{ t('tools.video.empty') }}</p>
            </div>

            <div v-else class="space-y-4">
              <div
                v-for="task in tasks"
                :key="task.id"
                class="rounded-xl border border-gray-200 p-3 dark:border-dark-700"
              >
                <div class="mb-2 flex items-center justify-between gap-2 text-xs text-gray-500 dark:text-dark-400">
                  <span class="truncate">{{ task.prompt }}</span>
                  <div class="flex shrink-0 items-center gap-2">
                    <span class="font-medium">{{ task.model }}</span>
                    <span
                      class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium"
                      :class="statusBadgeClass(task.status)"
                    >
                      <span v-if="task.status === 'pending'" class="h-1.5 w-1.5 animate-pulse rounded-full bg-current"></span>
                      {{ statusText(task.status) }}
                    </span>
                  </div>
                </div>

                <!-- 视频：成功后播放 -->
                <video
                  v-if="task.status === 'done' && task.videoUrl"
                  :src="task.videoUrl"
                  controls
                  playsinline
                  class="w-full rounded-lg bg-black"
                ></video>

                <!-- 进行中提示 -->
                <div v-else-if="task.status === 'pending'" class="flex h-40 flex-col items-center justify-center gap-2 rounded-lg bg-gray-50 text-sm text-gray-400 dark:bg-dark-700/50 dark:text-dark-400">
                  <Icon name="sync" size="lg" class="animate-spin" />
                  <span>{{ t('tools.video.polling', { elapsed: elapsedSeconds(task) }) }}</span>
                </div>

                <!-- 失败提示 -->
                <div v-else-if="task.error" class="flex h-40 flex-col items-center justify-center gap-1 rounded-lg bg-red-50 px-4 text-center text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">
                  <Icon name="xCircle" size="lg" />
                  <span class="break-all">{{ task.error }}</span>
                </div>

                <!-- 下载按钮 -->
                <div v-if="task.status === 'done' && task.videoUrl" class="mt-2 flex justify-end">
                  <button
                    type="button"
                    class="btn btn-secondary h-[32px] text-xs"
                    @click="downloadVideo(task)"
                  >
                    <Icon name="download" size="xs" />
                    {{ t('tools.video.download') }}
                  </button>
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
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ToolsApiKeyPanel from '@/components/tools/ToolsApiKeyPanel.vue'
import { GatewayError, useToolsSettings } from '@/composables/useToolsSettings'

type VideoTaskStatus = 'pending' | 'done' | 'failed' | 'expired' | 'error'

interface VideoTask {
  id: string
  prompt: string
  model: string
  status: VideoTaskStatus
  /** blob: URL（本地代理内容）或上游绝对 URL */
  videoUrl?: string
  /** 下载文件名后缀对应的 blob URL，用于触发下载 */
  error?: string
  createdAt: number
}

const POLL_INTERVAL_MS = 4000
const MAX_POLL_COUNT = 150 // 最长约 10 分钟

const { t } = useI18n()
const { hasApiKey, videoModel, gatewayFetch, fetchGatewayModels } = useToolsSettings()

const prompt = ref('')
const submitting = ref(false)
const error = ref('')

const availableModels = ref<string[]>([])
const tasks = ref<VideoTask[]>([])

/** 任务 id → 轮询定时器句柄 */
const pollTimers = new Map<string, ReturnType<typeof setTimeout>>()
/** 任务 id → 轮询次数 */
const pollCounts = new Map<string, number>()

const canGenerate = computed(() => !submitting.value && prompt.value.trim() !== '' && videoModel.value.trim() !== '' && hasApiKey.value)

function localizedError(err: unknown): string {
  if (err instanceof GatewayError) {
    if (err.status === 401) return t('tools.errors.unauthorized')
    if (err.status === 403) return t('tools.errors.forbidden')
    if (err.status === 429) return t('tools.errors.rateLimited')
    return err.message
  }
  return err instanceof Error ? err.message : t('tools.errors.unknown')
}

function statusBadgeClass(status: VideoTaskStatus): string {
  switch (status) {
    case 'done':
      return 'bg-green-50 text-green-700 dark:bg-green-900/30 dark:text-green-400'
    case 'pending':
      return 'bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
    default:
      return 'bg-red-50 text-red-700 dark:bg-red-900/30 dark:text-red-400'
  }
}

function statusText(status: VideoTaskStatus): string {
  return t(`tools.video.status.${status}`)
}

function elapsedSeconds(task: VideoTask): number {
  return Math.max(0, Math.round((Date.now() - task.createdAt) / 1000))
}

/** 提交响应中按优先级提取任务 id（上游为 xAI request_id，兼容 id / data.*） */
function extractTaskId(body: unknown): string {
  const root = body as Record<string, unknown> | null
  const data = root && typeof root.data === 'object' ? (root.data as Record<string, unknown>) : null
  for (const source of [root, data]) {
    if (!source) continue
    for (const key of ['request_id', 'id']) {
      const value = source[key]
      if (typeof value === 'string' && value) return value
      if (typeof value === 'number' && Number.isFinite(value)) return String(value)
    }
  }
  return ''
}

function stopPolling(taskId: string): void {
  const timer = pollTimers.get(taskId)
  if (timer) {
    clearTimeout(timer)
    pollTimers.delete(taskId)
  }
}

/** done 后的视频内容需要携带密钥拉取（网关本地代理路径），转成 blob URL 供播放 */
async function resolveVideoUrl(rawUrl: string): Promise<string> {
  if (/^https?:\/\//i.test(rawUrl)) {
    // 上游绝对地址：直接播放
    return rawUrl
  }
  const res = await gatewayFetch(rawUrl.startsWith('/') ? rawUrl : `/${rawUrl}`)
  const blob = await res.blob()
  return URL.createObjectURL(blob)
}

async function pollTask(task: VideoTask): Promise<void> {
  const count = (pollCounts.get(task.id) ?? 0) + 1
  pollCounts.set(task.id, count)

  try {
    const res = await gatewayFetch(`/v1/videos/generations/${encodeURIComponent(task.id)}`)
    const body: unknown = await res.json()
    const status = (body as { status?: unknown })?.status
    const statusStr = typeof status === 'string' ? status.toLowerCase() : ''

    if (statusStr === 'done') {
      const video = (body as { video?: { url?: unknown } })?.video
      const rawUrl = typeof video?.url === 'string' ? video.url : ''
      stopPolling(task.id)
      if (!rawUrl) {
        task.status = 'error'
        task.error = t('tools.video.missingUrl')
        return
      }
      try {
        task.videoUrl = await resolveVideoUrl(rawUrl)
        task.status = 'done'
      } catch (err) {
        task.status = 'error'
        task.error = localizedError(err)
      }
      return
    }

    if (statusStr === 'failed' || statusStr === 'expired') {
      stopPolling(task.id)
      task.status = statusStr
      return
    }

    // pending / 其他中间态：继续轮询
    if (count >= MAX_POLL_COUNT) {
      stopPolling(task.id)
      task.status = 'error'
      task.error = t('tools.video.timeout')
      return
    }
    pollTimers.set(task.id, setTimeout(() => pollTask(task), POLL_INTERVAL_MS))
  } catch (err) {
    // 单次轮询失败不终止任务（可能是瞬时网络问题），达到上限后放弃
    if (count >= MAX_POLL_COUNT) {
      stopPolling(task.id)
      task.status = 'error'
      task.error = localizedError(err)
      return
    }
    pollTimers.set(task.id, setTimeout(() => pollTask(task), POLL_INTERVAL_MS))
  }
}

async function submitTask(): Promise<void> {
  if (!canGenerate.value) return
  submitting.value = true
  error.value = ''
  try {
    const res = await gatewayFetch('/v1/videos/generations', {
      method: 'POST',
      body: JSON.stringify({
        model: videoModel.value.trim(),
        prompt: prompt.value.trim(),
      }),
    })
    const body: unknown = await res.json()
    const id = extractTaskId(body)
    if (!id) {
      error.value = t('tools.video.missingTaskId')
      return
    }

    const task: VideoTask = {
      id,
      prompt: prompt.value.trim(),
      model: videoModel.value.trim(),
      status: 'pending',
      createdAt: Date.now(),
    }
    tasks.value.unshift(task)
    pollCounts.set(id, 0)
    pollTimers.set(id, setTimeout(() => pollTask(task), POLL_INTERVAL_MS))
  } catch (err) {
    error.value = localizedError(err)
  } finally {
    submitting.value = false
  }
}

function downloadVideo(task: VideoTask): void {
  if (!task.videoUrl) return
  const link = document.createElement('a')
  const safeModel = task.model.replace(/[^a-zA-Z0-9._-]/g, '-') || 'video'
  link.download = `${safeModel}-${task.id}.mp4`
  link.href = task.videoUrl
  link.click()
}

function clearTasks(): void {
  for (const task of tasks.value) {
    stopPolling(task.id)
    if (task.videoUrl?.startsWith('blob:')) {
      URL.revokeObjectURL(task.videoUrl)
    }
  }
  tasks.value = []
}

onMounted(async () => {
  if (hasApiKey.value) {
    availableModels.value = await fetchGatewayModels()
  }
})

onBeforeUnmount(() => {
  for (const task of tasks.value) {
    stopPolling(task.id)
    if (task.videoUrl?.startsWith('blob:')) {
      URL.revokeObjectURL(task.videoUrl)
    }
  }
})
</script>
