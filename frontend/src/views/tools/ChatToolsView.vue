<template>
  <AppLayout>
    <div class="-m-[11px] flex h-[calc(100dvh-77px)] flex-col overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800 md:-m-[19px] lg:-m-[27px]">
      <div
        class="grid min-h-0 flex-1 grid-cols-1 divide-y divide-gray-200 overflow-y-auto dark:divide-dark-700 xl:divide-x xl:divide-y-0 xl:overflow-hidden"
        :class="sessionsCollapsed
          ? (paramsCollapsed ? 'xl:grid-cols-[48px_minmax(0,1fr)_48px]' : 'xl:grid-cols-[48px_minmax(0,1fr)_380px]')
          : (paramsCollapsed ? 'xl:grid-cols-[260px_minmax(0,1fr)_48px]' : 'xl:grid-cols-[260px_minmax(0,1fr)_380px]')"
      >
        <!-- 第一栏：对话记录（可折叠，向左收起） -->
        <section
          v-if="sessionsCollapsed"
          class="flex min-h-[56px] items-center justify-center gap-2 py-2 xl:min-h-0 xl:flex-col xl:justify-start xl:pt-3"
        >
          <button
            type="button"
            class="inline-flex h-6 w-6 items-center justify-center rounded-md border border-gray-200 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:border-dark-600 dark:hover:bg-dark-700 dark:hover:text-dark-200"
            :title="t('tools.chat.expandSessions')"
            @click="sessionsCollapsed = false"
          >
            <Icon name="chevronRight" size="sm" />
          </button>
          <span class="text-xs text-gray-500 dark:text-dark-400 xl:[writing-mode:vertical-rl]">
            {{ t('tools.chat.sessions') }}
          </span>
        </section>

        <ToolsPanel v-else :title="t('tools.chat.sessions')" icon="chatBubble" :badge="sessions.length">
          <template #actions>
            <button
              type="button"
              class="btn btn-secondary !px-2 !py-1 text-xs"
              @click="newSession"
            >
              <Icon name="plus" size="xs" />
              {{ t('tools.chat.newSession') }}
            </button>
            <button
              type="button"
              class="inline-flex h-6 w-6 items-center justify-center rounded-md border border-gray-200 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:border-dark-600 dark:hover:bg-dark-700 dark:hover:text-dark-200"
              :title="t('tools.chat.collapseSessions')"
              @click="sessionsCollapsed = true"
            >
              <Icon name="chevronLeft" size="sm" />
            </button>
          </template>

          <p class="mb-2 text-xs text-gray-400 dark:text-dark-500">{{ t('tools.chat.sessionLocalHint', { days: retentionDays }) }}</p>

          <p v-if="sessions.length === 0" class="py-6 text-center text-xs text-gray-400 dark:text-dark-500">
            {{ t('tools.chat.emptySessions') }}
          </p>

          <div v-else class="space-y-1.5">
            <div
              v-for="session in sessions"
              :key="session.id"
              class="group flex items-start gap-1 rounded-lg border px-2 py-1.5 transition-colors"
              :class="session.id === activeId
                ? 'border-primary-400 bg-primary-50/60 dark:border-primary-600 dark:bg-primary-900/20'
                : 'border-gray-200 hover:border-primary-300 hover:bg-primary-50/40 dark:border-dark-700 dark:hover:border-primary-700 dark:hover:bg-primary-900/10'"
            >
              <button type="button" class="min-w-0 flex-1 text-left" @click="switchSession(session.id)">
                <p class="truncate text-xs font-medium text-gray-800 dark:text-dark-100">{{ session.title }}</p>
                <p class="mt-0.5 text-[11px] text-gray-400 dark:text-dark-500">
                  {{ t('tools.chat.messageCount', { count: session.messages.length }) }}
                </p>
              </button>
              <button
                type="button"
                class="shrink-0 rounded p-0.5 text-gray-300 opacity-0 transition-opacity hover:text-red-500 group-hover:opacity-100"
                :title="t('tools.chat.deleteSession')"
                @click.stop="deleteSession(session.id)"
              >
                <Icon name="trash" size="xs" />
              </button>
            </div>
          </div>
        </ToolsPanel>

        <!-- 第二栏：对话 -->
        <ToolsPanel :title="activeTitle" :title-tooltip="firstUserMessage" icon="chat" body-class="!overflow-hidden !p-0">
          <template #actions>
            <button
              v-if="messages.length"
              type="button"
              class="btn btn-secondary !px-2 !py-1 text-xs"
              @click="clearConversation"
            >
              {{ t('tools.chat.clear') }}
            </button>
          </template>

          <div class="flex h-full min-h-0 flex-col">
            <!-- 消息区：唯一可滚动区域 -->
            <div ref="messagesContainer" class="min-h-0 flex-1 space-y-3 overflow-y-auto p-3">
              <div
                v-if="messages.length === 0"
                class="flex h-full flex-col items-center justify-center gap-2 text-gray-400 dark:text-dark-500"
              >
                <Icon name="chat" size="xl" />
                <p class="text-sm">{{ t('tools.chat.empty') }}</p>
              </div>

              <div
                v-for="(message, index) in messages"
                :key="index"
                class="flex"
                :class="message.role === 'user' ? 'justify-end' : 'justify-start'"
              >
                <div
                  class="group relative max-w-[85%] rounded-2xl px-3.5 py-2 text-sm leading-relaxed"
                  :class="message.role === 'user'
                    ? 'rounded-br-md border border-primary-200 bg-primary-100 text-primary-900 selection:bg-primary-600 selection:text-white dark:border-primary-700/60 dark:bg-primary-900/40 dark:text-primary-100'
                    : 'rounded-bl-md border border-gray-200 bg-gray-50 text-gray-900 dark:border-dark-700 dark:bg-dark-700/50 dark:text-dark-100'"
                >
                  <div v-if="message.role === 'user'" class="whitespace-pre-wrap break-words">{{ message.content }}</div>
                  <template v-else>
                    <div v-if="message.content" class="markdown-body prose prose-sm max-w-none break-words dark:prose-invert" v-html="renderMarkdown(message.content)"></div>
                    <div v-else class="flex items-center gap-1.5 py-1 text-gray-400 dark:text-dark-500">
                      <span class="h-1.5 w-1.5 animate-bounce rounded-full bg-current [animation-delay:0ms]"></span>
                      <span class="h-1.5 w-1.5 animate-bounce rounded-full bg-current [animation-delay:150ms]"></span>
                      <span class="h-1.5 w-1.5 animate-bounce rounded-full bg-current [animation-delay:300ms]"></span>
                    </div>
                  </template>
                  <button
                    v-if="message.content && (message.role === 'user' || !streaming)"
                    type="button"
                    class="absolute -bottom-5 hidden cursor-pointer items-center gap-1 text-xs text-gray-400 group-hover:flex hover:text-gray-600 dark:hover:text-dark-200"
                    :class="message.role === 'user' ? 'right-1' : 'left-1'"
                    :aria-label="t('common.copy')"
                    :title="t('common.copy')"
                    @click.stop="copyToClipboard(message.content)"
                  >
                    <Icon name="copy" size="xs" />
                    {{ t('common.copy') }}
                  </button>
                </div>
              </div>

              <!-- 流式中错误提示 -->
              <div v-if="streamError" class="flex justify-start">
                <div class="max-w-[85%] rounded-2xl rounded-bl-md border border-red-200 bg-red-50 px-3.5 py-2 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-900/20 dark:text-red-400">
                  {{ streamError }}
                </div>
              </div>
            </div>

            <!-- 输入区：分组/模型在顶部，输入框居中，发送按钮在右下 -->
            <div class="shrink-0 border-t border-gray-200 p-2.5 dark:border-dark-700">
              <div
                class="rounded-lg border border-gray-200 bg-white transition-colors focus-within:border-primary-400 focus-within:ring-1 focus-within:ring-primary-400/30 dark:border-dark-600 dark:bg-dark-800"
              >
                <div class="flex flex-wrap items-center px-2.5 pt-1.5">
                  <ToolsConnectionBar v-model="chatModel" tool="chat" inline class="min-w-0 flex-1" />
                </div>

                <textarea
                  ref="inputEl"
                  v-model="input"
                  rows="6"
                  :placeholder="t('tools.chat.inputPlaceholder')"
                  class="max-h-72 w-full resize-none border-0 bg-transparent px-2.5 py-1.5 text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-0 disabled:cursor-not-allowed disabled:opacity-60 dark:text-gray-100"
                  :disabled="streaming || !historyReady"
                  @keydown.enter.exact.prevent="sendMessage"
                ></textarea>

                <div class="flex items-center justify-end px-2 pb-1.5">
                  <button
                    v-if="!streaming"
                    type="button"
                    class="btn btn-primary !px-3.5 !py-1.5 text-sm"
                    :disabled="!canSend"
                    @click="sendMessage"
                  >
                    <Icon name="arrowUp" size="sm" />
                    {{ t('tools.chat.send') }}
                  </button>
                  <button
                    v-else
                    type="button"
                    class="btn btn-secondary !px-3.5 !py-1.5 text-sm"
                    @click="stopStreaming"
                  >
                    <Icon name="x" size="sm" />
                    {{ t('tools.chat.stop') }}
                  </button>
                </div>
              </div>

              <p class="mt-1.5 text-xs text-gray-400 dark:text-dark-500">
                {{ t('tools.chat.enterHint') }}
              </p>
            </div>
          </div>
        </ToolsPanel>

        <!-- 第三栏：对话参数（可折叠，两列布局） -->
        <section
          v-if="paramsCollapsed"
          class="flex min-h-[56px] items-center justify-center gap-2 py-2 xl:min-h-0 xl:flex-col xl:justify-start xl:pt-3"
        >
          <button
            type="button"
            class="inline-flex h-6 w-6 items-center justify-center rounded-md border border-gray-200 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:border-dark-600 dark:hover:bg-dark-700 dark:hover:text-dark-200"
            :title="t('tools.chat.expandParams')"
            @click="paramsCollapsed = false"
          >
            <Icon name="chevronLeft" size="sm" />
          </button>
          <span class="text-xs text-gray-500 dark:text-dark-400 xl:[writing-mode:vertical-rl]">
            {{ t('tools.chat.params') }}
          </span>
        </section>

        <ToolsPanel v-else :title="t('tools.chat.params')" icon="arrowsUpDown" body-class="flex flex-col">
          <template #actions>
            <button
              type="button"
              class="btn btn-ghost !px-2 !py-1 text-xs text-gray-500"
              @click="resetParams"
            >
              {{ t('tools.chat.resetParams') }}
            </button>
            <button
              type="button"
              class="inline-flex h-6 w-6 items-center justify-center rounded-md border border-gray-200 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:border-dark-600 dark:hover:bg-dark-700 dark:hover:text-dark-200"
              :title="t('tools.chat.collapseParams')"
              @click="paramsCollapsed = true"
            >
              <Icon name="chevronRight" size="sm" />
            </button>
          </template>

          <div class="grid grid-cols-2 gap-x-3 gap-y-2.5">
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-chat-temperature">
                {{ t('tools.chat.temperature') }}
              </label>
              <input
                id="tools-chat-temperature"
                v-model.number="temperature"
                type="number"
                min="0"
                max="2"
                step="0.1"
                placeholder="1"
                class="input !py-1.5 text-sm"
              />
              <p class="mt-0.5 text-[11px] leading-4 text-gray-400 dark:text-dark-500">{{ t('tools.chat.temperatureHint') }}</p>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-chat-topp">
                {{ t('tools.chat.topP') }}
              </label>
              <input
                id="tools-chat-topp"
                v-model.number="topP"
                type="number"
                min="0"
                max="1"
                step="0.05"
                placeholder="1"
                class="input !py-1.5 text-sm"
              />
              <p class="mt-0.5 text-[11px] leading-4 text-gray-400 dark:text-dark-500">{{ t('tools.chat.topPHint') }}</p>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-chat-maxtokens">
                {{ t('tools.chat.maxTokens') }}
              </label>
              <input
                id="tools-chat-maxtokens"
                v-model.number="maxTokens"
                type="number"
                min="0"
                step="1"
                placeholder="0"
                class="input !py-1.5 text-sm"
              />
              <p class="mt-0.5 text-[11px] leading-4 text-gray-400 dark:text-dark-500">{{ t('tools.chat.maxTokensHint') }}</p>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-chat-seed">
                {{ t('tools.chat.seed') }}
              </label>
              <input
                id="tools-chat-seed"
                v-model.number="seed"
                type="number"
                step="1"
                placeholder="随机"
                class="input !py-1.5 text-sm"
              />
              <p class="mt-0.5 text-[11px] leading-4 text-gray-400 dark:text-dark-500">{{ t('tools.chat.seedHint') }}</p>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-chat-frequency">
                {{ t('tools.chat.frequencyPenalty') }}
              </label>
              <input
                id="tools-chat-frequency"
                v-model.number="frequencyPenalty"
                type="number"
                min="-2"
                max="2"
                step="0.1"
                placeholder="0"
                class="input !py-1.5 text-sm"
              />
              <p class="mt-0.5 text-[11px] leading-4 text-gray-400 dark:text-dark-500">{{ t('tools.chat.frequencyPenaltyHint') }}</p>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-chat-presence">
                {{ t('tools.chat.presencePenalty') }}
              </label>
              <input
                id="tools-chat-presence"
                v-model.number="presencePenalty"
                type="number"
                min="-2"
                max="2"
                step="0.1"
                placeholder="0"
                class="input !py-1.5 text-sm"
              />
              <p class="mt-0.5 text-[11px] leading-4 text-gray-400 dark:text-dark-500">{{ t('tools.chat.presencePenaltyHint') }}</p>
            </div>
            <label class="col-span-2 flex h-[34px] w-full cursor-pointer items-center gap-2 rounded-lg border border-gray-200 px-2.5 dark:border-dark-600">
              <input
                v-model="streamEnabled"
                type="checkbox"
                class="h-3.5 w-3.5 shrink-0 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
              />
              <span class="text-xs font-medium text-gray-700 dark:text-dark-200">{{ t('tools.chat.stream') }}</span>
            </label>
          </div>

          <!-- 系统提示词：单独一行占满整列，并向下拉伸填满剩余空间 -->
          <div class="mt-3 flex min-h-0 flex-1 flex-col">
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-chat-system">
              {{ t('tools.chat.systemPrompt') }}
            </label>
            <textarea
              id="tools-chat-system"
              v-model="systemPrompt"
              rows="4"
              :placeholder="t('tools.chat.systemPromptPlaceholder')"
              class="input min-h-[120px] w-full flex-1 resize-none text-sm"
            ></textarea>
          </div>
        </ToolsPanel>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ToolsConnectionBar from '@/components/tools/ToolsConnectionBar.vue'
import ToolsPanel from '@/components/tools/ToolsPanel.vue'
import { useClipboard } from '@/composables/useClipboard'
import { GatewayError, useToolsSettings } from '@/composables/useToolsSettings'
import { toolChatAPI } from '@/api/toolChat'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'

interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
  createdAt?: number
}

interface ChatSession {
  id: string
  title: string
  createdAt: number
  messages: ChatMessage[]
}

const SESSIONS_STORAGE_KEY = 'tools_chat_sessions_v1'
const SESSIONS_OWNER_KEY = 'tools_chat_sessions_owner'
const SERVER_MIGRATION_KEY = 'tools_chat_history_server_migrated'
const MAX_SESSIONS = 20
const DEFAULT_CONTEXT_LIMIT = 20
const DEFAULT_RETENTION_DAYS = 7

function newSessionId(): string {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  return `s_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
}

function loadSessions(): ChatSession[] {
  try {
    const raw = localStorage.getItem(SESSIONS_STORAGE_KEY)
    if (!raw) return []
    const parsed: unknown = JSON.parse(raw)
    if (!Array.isArray(parsed)) return []
    const sessions: ChatSession[] = []
    for (const item of parsed) {
      const candidate = item as Partial<ChatSession>
      if (typeof candidate?.id !== 'string' || !Array.isArray(candidate.messages)) continue
      const createdAt = typeof candidate.createdAt === 'number' ? candidate.createdAt : Date.now()
      sessions.push({
        id: candidate.id,
        title: typeof candidate.title === 'string' && candidate.title ? candidate.title : '',
        createdAt,
        messages: candidate.messages
          .filter((m): m is ChatMessage => !!m && (m.role === 'user' || m.role === 'assistant') && typeof m.content === 'string')
          .slice(-200)
          .map((m) => ({
            role: m.role,
            content: m.content,
            createdAt: typeof m.createdAt === 'number' && m.createdAt > 0 ? m.createdAt : createdAt,
          })),
      })
    }
    return sessions.filter((session) => session.messages.length > 0).slice(0, MAX_SESSIONS)
  } catch {
    return []
  }
}

const { t } = useI18n()
const { hasApiKey, chatModel, gatewayFetch } = useToolsSettings()
const { copyToClipboard } = useClipboard()
const authStore = useAuthStore()
const appStore = useAppStore()

function currentOwnerId(): string {
  const id = authStore.user?.id
  return id === undefined || id === null ? '' : String(id)
}

/** 本地缓存只作为服务端历史不可用时的回退，且不能跨账号读取。 */
function loadSessionsForCurrentUser(): ChatSession[] {
  const owner = currentOwnerId()
  try {
    const storedOwner = localStorage.getItem(SESSIONS_OWNER_KEY) ?? ''
    if (storedOwner && owner && storedOwner !== owner) return []
    if (!storedOwner && owner) localStorage.setItem(SESSIONS_OWNER_KEY, owner)
  } catch {
    // 忽略存储失败
  }
  return loadSessions()
}

const sessions = ref<ChatSession[]>(loadSessionsForCurrentUser())
const activeId = ref<string>(sessions.value[0]?.id ?? '')
const historyReady = ref(false)
const contextLimit = ref(DEFAULT_CONTEXT_LIMIT)
const retentionDays = ref(DEFAULT_RETENTION_DAYS)

const input = ref('')
const streaming = ref(false)
const streamError = ref('')
const abortController = ref<AbortController | null>(null)

const temperature = ref<number | ''>('')
const topP = ref<number | ''>('')
const maxTokens = ref<number | ''>('')
const seed = ref<number | ''>('')
const frequencyPenalty = ref<number | ''>('')
const presencePenalty = ref<number | ''>('')
const systemPrompt = ref('')
const streamEnabled = ref(true)

// 左右两栏折叠状态（记忆到本地，刷新后保持）
const SESSIONS_COLLAPSED_KEY = 'tools_chat_sessions_collapsed'
const PARAMS_COLLAPSED_KEY = 'tools_chat_params_collapsed'
function readCollapsedFlag(key: string): boolean {
  try {
    return localStorage.getItem(key) === '1'
  } catch {
    return false
  }
}
function saveCollapsedFlag(key: string, value: boolean): void {
  try {
    localStorage.setItem(key, value ? '1' : '0')
  } catch {
    // 忽略存储失败
  }
}
const sessionsCollapsed = ref(readCollapsedFlag(SESSIONS_COLLAPSED_KEY))
watch(sessionsCollapsed, (value) => saveCollapsedFlag(SESSIONS_COLLAPSED_KEY, value))
const paramsCollapsed = ref(readCollapsedFlag(PARAMS_COLLAPSED_KEY))
watch(paramsCollapsed, (value) => saveCollapsedFlag(PARAMS_COLLAPSED_KEY, value))

const messagesContainer = ref<HTMLElement | null>(null)
const inputEl = ref<HTMLTextAreaElement | null>(null)


/** 会话标题最大显示宽度：中文/全角按 2、英文按 1 计 */
const TITLE_MAX_WIDTH = 20

function charDisplayWidth(char: string): number {
  const code = char.codePointAt(0) ?? 0
  const wide =
    (code >= 0x1100 && code <= 0x115f) ||
    (code >= 0x2e80 && code <= 0xa4cf) ||
    (code >= 0xac00 && code <= 0xd7a3) ||
    (code >= 0xf900 && code <= 0xfaff) ||
    (code >= 0xfe30 && code <= 0xfe4f) ||
    (code >= 0xff00 && code <= 0xff60) ||
    (code >= 0xffe0 && code <= 0xffe6) ||
    (code >= 0x1f300 && code <= 0x1faff) ||
    (code >= 0x20000 && code <= 0x3fffd)
  return wide ? 2 : 1
}

/** 把一段文本压成单行摘要，超出宽度时追加省略号 */
function summarizeTitle(text: string, maxWidth = TITLE_MAX_WIDTH): string {
  const flat = text.replace(/\s+/g, ' ').trim()
  if (!flat) return ''
  let width = 0
  let out = ''
  for (const char of flat) {
    const next = width + charDisplayWidth(char)
    if (next > maxWidth) return `${out}…`
    width = next
    out += char
  }
  return out
}

const activeSession = computed(() => sessions.value.find((session) => session.id === activeId.value) ?? null)
const messages = computed<ChatMessage[]>(() => activeSession.value?.messages ?? [])

/** 当前会话首条用户消息（用于标题与悬停全文） */
const firstUserMessage = computed(
  () => activeSession.value?.messages.find((message) => message.role === 'user')?.content?.trim() ?? '',
)

/** 面板标题：首条用户消息摘要，没有消息时显示「新对话」 */
const activeTitle = computed(() => summarizeTitle(firstUserMessage.value) || t('tools.chat.untitled'))


const canSend = computed(
  () =>
    historyReady.value &&
    !streaming.value &&
    input.value.trim() !== '' &&
    chatModel.value.trim() !== '' &&
    hasApiKey.value,
)

function persistLocal(): void {
  try {
    const owner = currentOwnerId()
    if (owner) localStorage.setItem(SESSIONS_OWNER_KEY, owner)
    localStorage.setItem(SESSIONS_STORAGE_KEY, JSON.stringify(sessions.value.slice(0, MAX_SESSIONS)))
  } catch {
    // 超出配额时放弃本地缓存，不影响当前会话
  }
}

let serverSaveTimer: ReturnType<typeof setTimeout> | null = null
let serverSaveQueue: Promise<void> = Promise.resolve()
let lastSaveErrorNotice = ''

/** 把服务端保存限制翻译成用户能理解、可操作的提示。 */
function serverSaveErrorMessage(err: unknown): string {
  const reason = extractApiErrorMessage(err, '').toLowerCase()
  if (reason.includes('content is too long')) return t('tools.chat.historyMessageTooLong')
  if (reason.includes('content is too large')) return t('tools.chat.historySessionTooLarge')
  if (reason.includes('too many sessions')) return t('tools.chat.historyTooManySessions')
  return t('tools.chat.historySaveFailed')
}

/** 合并短时间内的多次变更，并串行写入服务端，避免旧快照覆盖新快照。 */
function scheduleServerSave(): void {
  if (serverSaveTimer !== null) clearTimeout(serverSaveTimer)
  serverSaveTimer = setTimeout(() => {
    serverSaveTimer = null
    const snapshot = sessions.value.slice(0, MAX_SESSIONS).map((session) => ({
      id: session.id,
      title: session.title,
      createdAt: session.createdAt,
      messages: session.messages.map((message) => ({
        role: message.role,
        content: message.content,
        createdAt: message.createdAt,
      })),
    }))
    serverSaveQueue = serverSaveQueue
      .catch(() => undefined)
      .then(() => toolChatAPI.save(snapshot))
      .then(() => {
        lastSaveErrorNotice = ''
      })
      .catch((err) => {
        console.warn('Failed to save chat history to server', err)
        const message = serverSaveErrorMessage(err)
        // 同一类保存失败只提示一次，避免每次防抖重试都弹 toast。
        if (message !== lastSaveErrorNotice) {
          lastSaveErrorNotice = message
          appStore.showError(message)
        }
      })
  }, 300)
}

function persist(): void {
  persistLocal()
  scheduleServerSave()
}

function serverMigrationDone(): boolean {
  try {
    return localStorage.getItem(SERVER_MIGRATION_KEY) === '1'
  } catch {
    return false
  }
}

function markServerMigrationDone(): void {
  try {
    localStorage.setItem(SERVER_MIGRATION_KEY, '1')
  } catch {
    // 忽略存储失败
  }
}

/** 服务端历史是主数据源；首次上线时把浏览器里已有的旧记录迁移过去。 */
async function loadServerSessions(): Promise<void> {
  try {
    const result = await toolChatAPI.list()
    contextLimit.value = result.contextLimit > 0 ? result.contextLimit : DEFAULT_CONTEXT_LIMIT
    retentionDays.value = result.retentionDays > 0 ? result.retentionDays : DEFAULT_RETENTION_DAYS
    const remoteSessions = result.sessions.slice(0, MAX_SESSIONS)
    if (remoteSessions.length > 0) {
      sessions.value = remoteSessions
      markServerMigrationDone()
      persistLocal()
      return
    }
    if (sessions.value.length > 0 && !serverMigrationDone()) {
      await toolChatAPI.save(sessions.value.slice(0, MAX_SESSIONS))
      markServerMigrationDone()
      return
    }
    // 服务端明确为空时以服务端为准，避免已删除的记录被本地缓存复活。
    sessions.value = []
  } catch (err) {
    // 服务端不可用时退回本地缓存，不阻塞聊天。
    console.warn('Failed to load chat history from server', err)
  }
}

function createSession(): ChatSession {
  const session: ChatSession = {
    id: newSessionId(),
    title: t('tools.chat.untitled'),
    createdAt: Date.now(),
    messages: [],
  }
  sessions.value.unshift(session)
  if (sessions.value.length > MAX_SESSIONS) sessions.value = sessions.value.slice(0, MAX_SESSIONS)
  // 返回响应式代理，确保首条消息推入后视图和流式输出能正常更新。
  return sessions.value[0] ?? session
}

function newSession(): void {
  const current = activeSession.value
  if (streaming.value) stopStreaming()

  // 真实会话才需要清空输入；已经处于空白草稿时重复点击只是保持当前界面。
  if (current && current.messages.length > 0) {
    input.value = ''
  }
  // 兼容旧数据中遗留的 0 条消息会话，顺手清掉。
  if (current && current.messages.length === 0) {
    sessions.value = sessions.value.filter((session) => session.id !== current.id)
    persist()
  }

  // 只是切回空白草稿；真正发送第一条消息时才创建会话记录。
  activeId.value = ''
  streamError.value = ''
  void nextTick(() => inputEl.value?.focus())
}

function switchSession(id: string): void {
  if (streaming.value) stopStreaming()
  activeId.value = id
}

function deleteSession(id: string): void {
  if (streaming.value && activeId.value === id) stopStreaming()
  sessions.value = sessions.value.filter((session) => session.id !== id)
  if (activeId.value === id) activeId.value = sessions.value[0]?.id ?? ''
  persist()
}

/** 助手消息 Markdown 渲染（与 DocsView 相同的 marked + DOMPurify 管线） */
function renderMarkdown(content: string): string {
  try {
    const html = marked.parse(content) as string
    return DOMPurify.sanitize(html)
  } catch {
    return DOMPurify.sanitize(content)
  }
}

function localizedError(err: unknown): string {
  if (err instanceof GatewayError) {
    if (err.code === 'INSUFFICIENT_BALANCE') return t('tools.errors.insufficientBalance')
    if (err.status === 401) return t('tools.errors.unauthorized')
    if (err.status === 403) return t('tools.errors.forbidden')
    if (err.status === 429) return t('tools.errors.rateLimited')
    return err.message
  }
  if (err instanceof Error && err.name === 'AbortError') return t('tools.chat.stopped')
  return err instanceof Error ? err.message : t('tools.errors.unknown')
}

async function scrollToBottom(): Promise<void> {
  await nextTick()
  const el = messagesContainer.value
  if (el) el.scrollTop = el.scrollHeight
}

watch(() => messages.value.length, () => void scrollToBottom())
watch(activeId, () => void scrollToBottom())

// 同一浏览器切换账号时，重新按新账号加载历史；不能把上一个账号的会话留在内存里。
watch(
  () => authStore.user?.id,
  async (nextOwner, previousOwner) => {
    if (nextOwner === previousOwner) return
    if (streaming.value) stopStreaming()
    activeId.value = ''
    historyReady.value = false
    if (nextOwner === undefined || nextOwner === null) {
      sessions.value = []
      return
    }
    sessions.value = loadSessionsForCurrentUser()
    await loadServerSessions()
    if (!sessions.value.some((session) => session.id === activeId.value)) {
      activeId.value = sessions.value[0]?.id ?? ''
    }
    historyReady.value = true
    persistLocal()
  },
)

function stopStreaming(): void {
  abortController.value?.abort()
}

function clearConversation(): void {
  if (streaming.value) stopStreaming()
  const current = activeSession.value
  if (!current) return
  current.messages = []
  streamError.value = ''
  // 清空后同样回到空白草稿，不保留 0 条消息的会话记录。
  sessions.value = sessions.value.filter((session) => session.id !== current.id)
  activeId.value = ''
  persist()
}

function resetParams(): void {
  temperature.value = ''
  topP.value = ''
  maxTokens.value = ''
  seed.value = ''
  frequencyPenalty.value = ''
  presencePenalty.value = ''
  systemPrompt.value = ''
  streamEnabled.value = true
}

/** 从 SSE data 行解析增量文本，兼容 OpenAI chat.completion.chunk 结构 */
function extractDeltaText(jsonText: string): string {
  try {
    const chunk: unknown = JSON.parse(jsonText)
    const choices = (chunk as { choices?: unknown })?.choices
    if (!Array.isArray(choices) || choices.length === 0) return ''
    const delta = (choices[0] as { delta?: { content?: unknown } })?.delta
    const content = delta?.content
    return typeof content === 'string' ? content : ''
  } catch {
    return ''
  }
}

/** 非流式响应：从 choices[0].message.content 取正文（兼容字符串与分段数组） */
function extractMessageContent(body: unknown): string {
  const choices = (body as { choices?: unknown })?.choices
  if (!Array.isArray(choices) || choices.length === 0) return ''
  const content = (choices[0] as { message?: { content?: unknown } })?.message?.content
  if (typeof content === 'string') return content
  if (Array.isArray(content)) {
    return content
      .map((part) => (part && typeof (part as { text?: unknown }).text === 'string' ? (part as { text: string }).text : ''))
      .join('')
  }
  return ''
}

function buildPayload(history: ChatMessage[]): Record<string, unknown> {
  const payloadMessages: { role: string; content: string }[] = []
  if (systemPrompt.value.trim()) {
    payloadMessages.push({ role: 'system', content: systemPrompt.value.trim() })
  }
  // 服务端保存完整历史，但每次请求只携带最近 N 条（默认 20 条）上下文。
  for (const message of history.slice(-contextLimit.value)) {
    payloadMessages.push({ role: message.role, content: message.content })
  }

  const payload: Record<string, unknown> = {
    model: chatModel.value.trim(),
    messages: payloadMessages,
    stream: streamEnabled.value,
  }
  if (typeof temperature.value === 'number') payload.temperature = temperature.value
  if (typeof topP.value === 'number') payload.top_p = topP.value
  if (typeof maxTokens.value === 'number' && maxTokens.value > 0) payload.max_tokens = maxTokens.value
  if (typeof seed.value === 'number') payload.seed = seed.value
  if (typeof frequencyPenalty.value === 'number') payload.frequency_penalty = frequencyPenalty.value
  if (typeof presencePenalty.value === 'number') payload.presence_penalty = presencePenalty.value
  return payload
}

async function sendMessage(): Promise<void> {
  const text = input.value.trim()
  if (!canSend.value || !text) return

  // 空白草稿在发送第一条消息时才创建会话，避免出现 0 条消息的记录。
  let session = activeSession.value
  if (!session) {
    session = createSession()
    activeId.value = session.id
  }

  streamError.value = ''
  input.value = ''

  // 第一条用户消息充当会话标题
  const userMessageCount = session.messages.filter((m) => m.role === 'user').length
  if (userMessageCount === 0) {
    session.title = summarizeTitle(text) || t('tools.chat.untitled')
  }
  session.messages.push({ role: 'user', content: text, createdAt: Date.now() })
  persist()

  // 注意：必须通过数组下标访问拿到 Vue 的响应式代理再修改。
  // 直接修改 push 进去的原始对象不会触发视图更新，表现为"流式输出变成最后一次性出现"。
  const assistantIndex = session.messages.push({ role: 'assistant', content: '', createdAt: Date.now() }) - 1
  const assistantMessage = (): ChatMessage | undefined => session.messages[assistantIndex]

  const controller = new AbortController()
  abortController.value = controller
  streaming.value = true

  try {
    // 历史消息不含刚推入的空助手气泡
    const history = session.messages.slice(0, -1)
    const res = await gatewayFetch('/v1/chat/completions', {
      method: 'POST',
      signal: controller.signal,
      body: JSON.stringify(buildPayload(history)),
    })

    if (!streamEnabled.value) {
      const body: unknown = await res.json()
      const assistant = assistantMessage()
      if (assistant) assistant.content = extractMessageContent(body)
      return
    }

    if (!res.body) throw new GatewayError(t('tools.errors.noStream'), 0)

    const reader = res.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    for (;;) {
      const { done, value } = await reader.read()
      if (done) break
      buffer += decoder.decode(value, { stream: true })

      const lines = buffer.split('\n')
      buffer = lines.pop() ?? ''

      for (const line of lines) {
        const trimmed = line.trim()
        if (!trimmed.startsWith('data:')) continue
        const payload = trimmed.slice(5).trim()
        if (payload === '[DONE]') continue
        const delta = extractDeltaText(payload)
        if (delta) {
          const assistant = assistantMessage()
          if (assistant) assistant.content += delta
          void scrollToBottom()
        }
      }
    }
  } catch (err) {
    if (!(err instanceof Error && err.name === 'AbortError')) {
      streamError.value = localizedError(err)
    }
    // 中止或出错时，保留已生成的部分内容
  } finally {
    streaming.value = false
    abortController.value = null
    const assistant = assistantMessage()
    if (assistant && !assistant.content) {
      // 中断或空响应时不要留下空的助手气泡
      session.messages.splice(assistantIndex, 1)
    }
    persist()
    void scrollToBottom()
  }
}

onMounted(async () => {
  try {
    await loadServerSessions()
  } finally {
    // 没有会话时保持空白草稿；不预先创建 0 条消息的会话。
    if (!sessions.value.some((session) => session.id === activeId.value)) {
      activeId.value = sessions.value[0]?.id ?? ''
    }
    historyReady.value = true
    persistLocal()
  }
})
</script>