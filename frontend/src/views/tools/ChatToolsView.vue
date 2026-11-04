<template>
  <div class="space-y-4">
    <!-- 页头 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ t('tools.chat.title') }}
        </h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
          {{ t('tools.chat.description') }}
        </p>
      </div>
    </div>

    <!-- 密钥面板 -->
    <ToolsApiKeyPanel />

    <!-- 模型选择 -->
    <div class="rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800">
      <div class="flex flex-wrap items-center gap-3">
        <label class="text-sm font-medium text-gray-700 dark:text-dark-200" for="tools-chat-model">
          {{ t('tools.chat.model') }}
        </label>
        <input
          id="tools-chat-model"
          v-model="chatModel"
          list="tools-chat-models"
          type="text"
          :placeholder="t('tools.chat.modelPlaceholder')"
          class="input h-[38px] min-w-[240px] flex-1 sm:max-w-md"
          autocomplete="off"
          spellcheck="false"
        />
        <datalist id="tools-chat-models">
          <option v-for="m in availableModels" :key="m" :value="m" />
        </datalist>
        <button
          type="button"
          class="btn btn-secondary h-[38px] text-sm"
          :disabled="loadingModels"
          @click="loadModels"
        >
          <Icon name="refresh" size="sm" :class="{ 'animate-spin': loadingModels }" />
          {{ t('tools.chat.refreshModels') }}
        </button>
      </div>
      <p v-if="availableModels.length" class="mt-2 text-xs text-gray-400 dark:text-dark-500">
        {{ t('tools.chat.modelCount', { count: availableModels.length }) }}
      </p>
    </div>

    <!-- 对话窗口 -->
    <div class="flex h-[calc(100vh-19rem)] min-h-[420px] flex-col overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800">
      <!-- 消息区 -->
      <div ref="messagesContainer" class="flex-1 space-y-4 overflow-y-auto p-4">
        <div v-if="messages.length === 0" class="flex h-full flex-col items-center justify-center gap-2 text-gray-400 dark:text-dark-500">
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
            class="group relative max-w-[85%] rounded-2xl px-4 py-2.5 text-sm leading-relaxed"
            :class="message.role === 'user'
              ? 'rounded-br-md bg-primary-600 text-white'
              : 'rounded-bl-md border border-gray-200 bg-gray-50 text-gray-900 dark:border-dark-700 dark:bg-dark-700/50 dark:text-dark-100'"
          >
            <!-- 用户消息：纯文本 -->
            <div v-if="message.role === 'user'" class="whitespace-pre-wrap break-words">{{ message.content }}</div>
            <!-- 助手消息：Markdown 渲染 -->
            <template v-else>
              <div v-if="message.content" class="markdown-body prose prose-sm max-w-none break-words dark:prose-invert" v-html="renderMarkdown(message.content)"></div>
              <div v-else class="flex items-center gap-1.5 py-1 text-gray-400 dark:text-dark-500">
                <span class="h-1.5 w-1.5 animate-bounce rounded-full bg-current [animation-delay:0ms]"></span>
                <span class="h-1.5 w-1.5 animate-bounce rounded-full bg-current [animation-delay:150ms]"></span>
                <span class="h-1.5 w-1.5 animate-bounce rounded-full bg-current [animation-delay:300ms]"></span>
              </div>
              <span v-if="message.content && !streaming" class="absolute -bottom-5 left-1 hidden cursor-pointer items-center gap-1 text-xs text-gray-400 group-hover:flex hover:text-gray-600 dark:hover:text-dark-200" @click="copyToClipboard(message.content)">
                <Icon name="copy" size="xs" />
                {{ t('common.copy') }}
              </span>
            </template>
          </div>
        </div>

        <!-- 流式中错误提示 -->
        <div v-if="streamError" class="flex justify-start">
          <div class="max-w-[85%] rounded-2xl rounded-bl-md border border-red-200 bg-red-50 px-4 py-2.5 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-900/20 dark:text-red-400">
            {{ streamError }}
          </div>
        </div>
      </div>

      <!-- 输入区 -->
      <div class="border-t border-gray-200 p-3 dark:border-dark-700">
        <div class="flex items-end gap-2">
          <textarea
            v-model="input"
            rows="1"
            :placeholder="t('tools.chat.inputPlaceholder')"
            class="input max-h-40 flex-1 resize-none"
            :disabled="streaming"
            @keydown.enter.exact.prevent="sendMessage"
          ></textarea>
          <button
            v-if="!streaming"
            type="button"
            class="btn btn-primary h-[38px] px-4"
            :disabled="!canSend"
            @click="sendMessage"
          >
            {{ t('tools.chat.send') }}
          </button>
          <button
            v-else
            type="button"
            class="btn btn-secondary h-[38px] px-4"
            @click="stopStreaming"
          >
            <Icon name="x" size="sm" />
            {{ t('tools.chat.stop') }}
          </button>
          <button
            v-if="messages.length && !streaming"
            type="button"
            class="btn btn-secondary h-[38px] px-3"
            :title="t('tools.chat.clear')"
            @click="clearConversation"
          >
            <Icon name="trash" size="sm" />
          </button>
        </div>
        <p class="mt-1.5 text-xs text-gray-400 dark:text-dark-500">
          {{ t('tools.chat.enterHint') }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import Icon from '@/components/icons/Icon.vue'
import ToolsApiKeyPanel from '@/components/tools/ToolsApiKeyPanel.vue'
import { useClipboard } from '@/composables/useClipboard'
import { GatewayError, useToolsSettings } from '@/composables/useToolsSettings'

interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
}

const { t } = useI18n()
const { hasApiKey, chatModel, gatewayFetch, fetchGatewayModels } = useToolsSettings()
const { copyToClipboard } = useClipboard()

const messages = ref<ChatMessage[]>([])
const input = ref('')
const streaming = ref(false)
const streamError = ref('')
const abortController = ref<AbortController | null>(null)

const availableModels = ref<string[]>([])
const loadingModels = ref(false)

const messagesContainer = ref<HTMLElement | null>(null)

const canSend = computed(() => !streaming.value && input.value.trim() !== '' && chatModel.value.trim() !== '' && hasApiKey.value)

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

watch(() => messages.value.length, () => scrollToBottom())

async function loadModels(): Promise<void> {
  if (!hasApiKey.value || loadingModels.value) return
  loadingModels.value = true
  try {
    availableModels.value = await fetchGatewayModels()
  } finally {
    loadingModels.value = false
  }
}

function stopStreaming(): void {
  abortController.value?.abort()
}

function clearConversation(): void {
  if (streaming.value) stopStreaming()
  messages.value = []
  streamError.value = ''
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

async function sendMessage(): Promise<void> {
  const text = input.value.trim()
  if (!canSend.value || !text) return

  streamError.value = ''
  input.value = ''
  messages.value.push({ role: 'user', content: text })

  const assistantMessage: ChatMessage = { role: 'assistant', content: '' }
  messages.value.push(assistantMessage)

  const controller = new AbortController()
  abortController.value = controller
  streaming.value = true

  try {
    const res = await gatewayFetch('/v1/chat/completions', {
      method: 'POST',
      signal: controller.signal,
      body: JSON.stringify({
        model: chatModel.value.trim(),
        messages: messages.value
          .slice(0, -1)
          .map((m) => ({ role: m.role, content: m.content })),
        stream: true,
      }),
    })

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
          assistantMessage.content += delta
          scrollToBottom()
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
    if (!assistantMessage.content && streamError.value) {
      // 完全没有产出则移除空的助手气泡
      messages.value = messages.value.filter((m) => m !== assistantMessage)
    }
    scrollToBottom()
  }
}

onMounted(() => {
  if (hasApiKey.value) loadModels()
})
</script>
