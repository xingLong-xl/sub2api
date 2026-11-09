<template>
  <AppLayout>
    <div class="-m-[11px] flex h-[calc(100dvh-77px)] flex-col overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800 md:-m-[19px] lg:-m-[27px]">
      <div
        class="grid min-h-0 flex-1 grid-cols-1 divide-y divide-gray-200 overflow-y-auto dark:divide-dark-700 xl:grid-cols-[380px_minmax(0,1fr)_300px] xl:divide-x xl:divide-y-0 xl:overflow-hidden"
      >
        <!-- 第一栏：生成参数 -->
        <ToolsPanel :title="t('tools.video.params')" icon="cog">
          <div class="space-y-3">
            <ToolsConnectionBar
              v-model="videoModel"
              tool="video"
              inline
              class="border-b border-gray-100 pb-3 dark:border-dark-700"
            />
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-video-prompt">
                {{ t('tools.video.prompt') }}
              </label>
              <textarea
                id="tools-video-prompt"
                v-model="prompt"
                rows="5"
                :placeholder="t('tools.video.promptPlaceholder')"
                class="input resize-none text-sm"
              ></textarea>
            </div>

            <div class="grid grid-cols-2 gap-2">
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-video-duration">
                  {{ t('tools.video.duration') }}
                </label>
                <input
                  id="tools-video-duration"
                  v-model.number="duration"
                  type="number"
                  min="1"
                  max="30"
                  :placeholder="t('tools.video.durationHint')"
                  class="input !py-1.5 text-sm"
                />
              </div>
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-video-count">
                  {{ t('tools.video.count') }}
                </label>
                <input
                  id="tools-video-count"
                  v-model.number="count"
                  type="number"
                  min="1"
                  max="4"
                  class="input !py-1.5 text-sm"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-2">
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-video-ratio">
                  {{ t('tools.video.aspectRatio') }}
                </label>
                <select id="tools-video-ratio" v-model="aspectRatio" class="input !py-1.5 text-sm">
                  <option value="">{{ t('tools.video.resolutionDefault') }}</option>
                  <option value="16:9">16:9</option>
                  <option value="9:16">9:16</option>
                  <option value="1:1">1:1</option>
                  <option value="4:3">4:3</option>
                  <option value="3:4">3:4</option>
                </select>
              </div>
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-video-resolution">
                  {{ t('tools.video.resolution') }}
                </label>
                <select id="tools-video-resolution" v-model="resolution" class="input !py-1.5 text-sm">
                  <option value="">{{ t('tools.video.resolutionDefault') }}</option>
                  <option value="480p">480p</option>
                  <option value="720p">720p</option>
                  <option value="1080p">1080p</option>
                </select>
              </div>
            </div>

            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-video-image">
                {{ t('tools.video.referenceImage') }}
              </label>
              <input
                id="tools-video-image"
                v-model="referenceImage"
                type="text"
                :placeholder="t('tools.video.referenceImagePlaceholder')"
                class="input !py-1.5 text-sm"
                autocomplete="off"
                spellcheck="false"
              />
            </div>

            <button
              type="button"
              class="btn btn-primary w-full"
              :disabled="!canGenerate"
              @click="submitTask()"
            >
              <Icon name="play" size="sm" />
              {{ submitting ? t('tools.video.submitting') : t('tools.video.generate') }}
            </button>

            <p class="text-xs text-gray-400 dark:text-dark-500">{{ t('tools.video.asyncHint') }}</p>

            <p v-if="error" class="rounded-lg bg-red-50 px-2.5 py-2 text-xs text-red-600 dark:bg-red-900/20 dark:text-red-400">
              {{ error }}
            </p>
          </div>
        </ToolsPanel>

        <!-- 第二栏：任务与播放 -->
        <ToolsPanel :title="t('tools.video.tasks')" icon="play" :badge="tasks.length">
          <template #actions>
            <button
              v-if="tasks.length"
              type="button"
              class="btn btn-secondary !px-2 !py-1 text-xs"
              @click="clearTasks"
            >
              {{ t('tools.common.clear') }}
            </button>
          </template>

          <p v-if="tasks.length === 0" class="flex h-full min-h-[240px] flex-col items-center justify-center gap-2 text-center text-gray-400 dark:text-dark-500">
            <Icon name="play" size="xl" />
            <span class="text-sm">{{ t('tools.video.empty') }}</span>
          </p>

          <div v-else class="space-y-3">
            <article
              v-for="task in tasks"
              :key="task.id"
              class="rounded-lg border border-gray-200 p-2.5 dark:border-dark-700"
            >
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0">
                  <p class="truncate text-xs text-gray-700 dark:text-dark-200">{{ task.prompt }}</p>
                  <p class="mt-0.5 truncate text-[11px] text-gray-400 dark:text-dark-500">
                    {{ task.model }} · {{ task.id }}
                  </p>
                </div>
                <span
                  class="shrink-0 rounded-full px-2 py-0.5 text-[11px] font-medium"
                  :class="statusBadgeClass(task.status)"
                >
                  {{ statusText(task.status) }}
                </span>
              </div>

              <video
                v-if="task.videoUrl"
                :src="task.videoUrl"
                controls
                playsinline
                class="mt-2 w-full rounded-md bg-black"
              ></video>

              <p v-if="task.status === 'pending'" class="mt-2 text-[11px] text-blue-600 dark:text-blue-400">
                {{ t('tools.video.polling', { elapsed: elapsedSeconds(task) }) }}
              </p>
              <p v-if="task.error" class="mt-2 text-[11px] text-red-600 dark:text-red-400">{{ task.error }}</p>

              <div class="mt-2 flex items-center gap-1.5">
                <button
                  v-if="task.videoUrl"
                  type="button"
                  class="btn btn-secondary !px-2 !py-1 text-xs"
                  @click="downloadVideo(task)"
                >
                  <Icon name="download" size="xs" />
                  {{ t('tools.common.download') }}
                </button>
                <button
                  v-if="task.status === 'failed' || task.status === 'expired' || task.status === 'error'"
                  type="button"
                  class="btn btn-secondary !px-2 !py-1 text-xs"
                  :disabled="submitting"
                  @click="retryTask(task)"
                >
                  <Icon name="refresh" size="xs" />
                  {{ t('tools.video.retry') }}
                </button>
                <button
                  type="button"
                  class="btn btn-ghost !px-2 !py-1 text-xs text-gray-400"
                  @click="removeTask(task)"
                >
                  <Icon name="trash" size="xs" />
                </button>
              </div>
            </article>
          </div>
        </ToolsPanel>

        <!-- 第三栏：历史记录（暂为会话内） -->
        <ToolsPanel :title="t('tools.common.history')" icon="clock" :badge="doneTasks.length">
          <p class="mb-2 text-xs text-gray-400 dark:text-dark-500">{{ t('tools.video.historyHint') }}</p>
          <p v-if="doneTasks.length === 0" class="py-6 text-center text-xs text-gray-400 dark:text-dark-500">
            {{ t('tools.common.emptyHistory') }}
          </p>
          <div v-else class="space-y-2">
            <button
              v-for="task in doneTasks"
              :key="task.id"
              type="button"
              class="w-full rounded-lg border border-gray-200 p-2 text-left transition-colors hover:border-primary-300 hover:bg-primary-50/40 dark:border-dark-700 dark:hover:border-primary-700 dark:hover:bg-primary-900/10"
              @click="reuseTask(task)"
            >
              <div class="flex items-center justify-between gap-2 text-[11px] text-gray-400 dark:text-dark-500">
                <span>{{ formatTime(task.createdAt) }}</span>
                <span class="truncate">{{ task.model }}</span>
              </div>
              <p class="mt-1 truncate text-xs text-gray-700 dark:text-dark-200">{{ task.prompt }}</p>
              <div class="mt-1.5 flex items-center gap-2 text-[11px] text-primary-600 dark:text-primary-400">
                <span class="inline-flex items-center gap-1">
                  <Icon name="refresh" size="xs" />
                  {{ t('tools.common.reuse') }}
                </span>
                <span
                  class="inline-flex items-center gap-1"
                  @click.stop="downloadVideo(task)"
                >
                  <Icon name="download" size="xs" />
                  {{ t('tools.common.download') }}
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
import { computed, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ToolsConnectionBar from '@/components/tools/ToolsConnectionBar.vue'
import ToolsPanel from '@/components/tools/ToolsPanel.vue'
import { GatewayError, useToolsSettings } from '@/composables/useToolsSettings'

type VideoTaskStatus = 'pending' | 'done' | 'failed' | 'expired' | 'error'

interface VideoParams {
  duration: number | ''
  aspectRatio: string
  resolution: string
  count: number
  referenceImage: string
}

interface VideoTask {
  id: string
  prompt: string
  model: string
  status: VideoTaskStatus
  createdAt: number
  params: VideoParams
  videoUrl?: string
  error?: string
}

const POLL_INTERVAL_MS = 4000
const MAX_POLL_COUNT = 150 // 最长约 10 分钟

const { t } = useI18n()
const { hasApiKey, videoModel, gatewayFetch } = useToolsSettings()

const prompt = ref('')
const duration = ref<number | ''>('')
const count = ref(1)
const aspectRatio = ref('')
const resolution = ref('')
const referenceImage = ref('')

const submitting = ref(false)
const error = ref('')

const tasks = ref<VideoTask[]>([])

const doneTasks = computed(() => tasks.value.filter((task) => task.status === 'done' && task.videoUrl))

/** 任务 id → 轮询定时器句柄 */
const pollTimers = new Map<string, ReturnType<typeof setTimeout>>()
/** 任务 id → 轮询次数 */
const pollCounts = new Map<string, number>()

const canGenerate = computed(
  () => !submitting.value && prompt.value.trim() !== '' && videoModel.value.trim() !== '' && hasApiKey.value,
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

function formatTime(timestamp: number): string {
  return new Date(timestamp).toLocaleTimeString()
}

function elapsedSeconds(task: VideoTask): number {
  return Math.max(0, Math.round((Date.now() - task.createdAt) / 1000))
}

function currentParams(): VideoParams {
  return {
    duration: duration.value,
    aspectRatio: aspectRatio.value,
    resolution: resolution.value,
    count: count.value,
    referenceImage: referenceImage.value,
  }
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
  const countdown = (pollCounts.get(task.id) ?? 0) + 1
  pollCounts.set(task.id, countdown)

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
    if (countdown >= MAX_POLL_COUNT) {
      stopPolling(task.id)
      task.status = 'error'
      task.error = t('tools.video.timeout')
      return
    }
    pollTimers.set(task.id, setTimeout(() => void pollTask(task), POLL_INTERVAL_MS))
  } catch (err) {
    // 单次轮询失败不终止任务（可能是瞬时网络问题），达到上限后放弃
    if (countdown >= MAX_POLL_COUNT) {
      stopPolling(task.id)
      task.status = 'error'
      task.error = localizedError(err)
      return
    }
    pollTimers.set(task.id, setTimeout(() => void pollTask(task), POLL_INTERVAL_MS))
  }
}

function buildPayload(taskPrompt: string, model: string, params: VideoParams): Record<string, unknown> {
  const payload: Record<string, unknown> = {
    model,
    prompt: taskPrompt,
  }
  if (typeof params.duration === 'number' && params.duration > 0) payload.duration = params.duration
  if (params.aspectRatio) payload.aspect_ratio = params.aspectRatio
  if (params.resolution) payload.resolution = params.resolution
  if (params.count > 1) payload.n = params.count
  if (params.referenceImage.trim()) payload.image = params.referenceImage.trim()
  return payload
}

async function submitTask(promptOverride?: string, modelOverride?: string, paramsOverride?: VideoParams): Promise<void> {
  const taskPrompt = (promptOverride ?? prompt.value).trim()
  const taskModel = (modelOverride ?? videoModel.value).trim()
  const params = paramsOverride ?? currentParams()
  if (!taskPrompt || !taskModel || !hasApiKey.value) return

  submitting.value = true
  error.value = ''
  try {
    const res = await gatewayFetch('/v1/videos/generations', {
      method: 'POST',
      body: JSON.stringify(buildPayload(taskPrompt, taskModel, params)),
    })
    const body: unknown = await res.json()
    const id = extractTaskId(body)
    if (!id) {
      error.value = t('tools.video.missingTaskId')
      return
    }

    const task: VideoTask = {
      id,
      prompt: taskPrompt,
      model: taskModel,
      status: 'pending',
      createdAt: Date.now(),
      params,
    }
    tasks.value.unshift(task)
    // 取回数组中的响应式代理再轮询：直接改 unshift 进去的原始对象不会触发视图更新，
    // 表现为任务卡片状态永远停在"生成中"。
    const trackedTask = tasks.value[0]
    pollCounts.set(id, 0)
    pollTimers.set(id, setTimeout(() => void pollTask(trackedTask), POLL_INTERVAL_MS))
  } catch (err) {
    error.value = localizedError(err)
  } finally {
    submitting.value = false
  }
}

function retryTask(task: VideoTask): void {
  void submitTask(task.prompt, task.model, task.params)
}

/** 历史记录点击：回填参数，便于微调后重新生成 */
function reuseTask(task: VideoTask): void {
  prompt.value = task.prompt
  videoModel.value = task.model
  duration.value = task.params.duration
  aspectRatio.value = task.params.aspectRatio
  resolution.value = task.params.resolution
  count.value = task.params.count
  referenceImage.value = task.params.referenceImage
}

function downloadVideo(task: VideoTask): void {
  if (!task.videoUrl) return
  const link = document.createElement('a')
  const safeModel = task.model.replace(/[^a-zA-Z0-9._-]/g, '-') || 'video'
  link.download = `${safeModel}-${task.id}.mp4`
  link.href = task.videoUrl
  link.click()
}

function removeTask(task: VideoTask): void {
  stopPolling(task.id)
  if (task.videoUrl?.startsWith('blob:')) {
    URL.revokeObjectURL(task.videoUrl)
  }
  tasks.value = tasks.value.filter((item) => item.id !== task.id)
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

onBeforeUnmount(() => {
  for (const task of tasks.value) {
    stopPolling(task.id)
    if (task.videoUrl?.startsWith('blob:')) {
      URL.revokeObjectURL(task.videoUrl)
    }
  }
})
</script>