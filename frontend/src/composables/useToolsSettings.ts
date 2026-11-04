import { computed, ref, watch } from 'vue'
import { buildGatewayUrl } from '@/api/url'

/**
 * 工具（聊天 / 图像 / 视频）页共享设置：
 * - API 密钥：调用本站网关 /v1/* 端点必须携带用户自己的 sk- 密钥，
 *   网关不认登录会话 Cookie，因此需要用户在工具页填一次密钥（存 localStorage）。
 * - 各工具记住的模型名：分别持久化，跨页面刷新保留。
 *
 * 模块级单例：三个工具页共享同一份密钥与模型偏好。
 */

const STORAGE_KEYS = {
  apiKey: 'tools_api_key',
  chatModel: 'tools_chat_model',
  imageModel: 'tools_image_model',
  videoModel: 'tools_video_model',
} as const

function readStorage(key: string): string {
  try {
    return localStorage.getItem(key) ?? ''
  } catch {
    return ''
  }
}

function writeStorage(key: string, value: string): void {
  try {
    if (value) {
      localStorage.setItem(key, value)
    } else {
      localStorage.removeItem(key)
    }
  } catch {
    // localStorage 不可用时静默降级为仅内存保存
  }
}

const apiKey = ref(readStorage(STORAGE_KEYS.apiKey))
const chatModel = ref(readStorage(STORAGE_KEYS.chatModel))
const imageModel = ref(readStorage(STORAGE_KEYS.imageModel))
const videoModel = ref(readStorage(STORAGE_KEYS.videoModel))

watch(apiKey, (v) => writeStorage(STORAGE_KEYS.apiKey, v.trim()))
watch(chatModel, (v) => writeStorage(STORAGE_KEYS.chatModel, v.trim()))
watch(imageModel, (v) => writeStorage(STORAGE_KEYS.imageModel, v.trim()))
watch(videoModel, (v) => writeStorage(STORAGE_KEYS.videoModel, v.trim()))

/** 网关错误：带 HTTP 状态码，视图层据此做本地化提示 */
export class GatewayError extends Error {
  readonly status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'GatewayError'
    this.status = status
  }
}

/** 从网关/上游错误响应体中尽力提取可读信息 */
async function toGatewayError(res: Response): Promise<GatewayError> {
  let message = `HTTP ${res.status}`
  try {
    const body: unknown = await res.json()
    if (body && typeof body === 'object') {
      const err = (body as Record<string, unknown>).error
      if (typeof err === 'string' && err) {
        message = err
      } else if (err && typeof err === 'object') {
        const msg = (err as Record<string, unknown>).message
        if (typeof msg === 'string' && msg) message = msg
      } else {
        const msg = (body as Record<string, unknown>).message
        if (typeof msg === 'string' && msg) message = msg
      }
    }
  } catch {
    // 非 JSON 响应，保留 HTTP 状态码
  }
  return new GatewayError(message, res.status)
}

export function useToolsSettings() {
  function setApiKey(key: string): void {
    apiKey.value = key.trim()
  }

  function authHeaders(): Record<string, string> {
    return { Authorization: `Bearer ${apiKey.value.trim()}` }
  }

  const hasApiKey = computed(() => apiKey.value.trim() !== '')

  /**
   * 网关请求统一入口：自动附带密钥，非 2xx 响应统一转成 GatewayError。
   * 返回原始 Response，便于调用方按需读流（聊天 SSE）或读 JSON。
   */
  async function gatewayFetch(path: string, init: RequestInit = {}): Promise<Response> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...authHeaders(),
    }
    const res = await fetch(buildGatewayUrl(path), { ...init, headers })
    if (!res.ok) {
      throw await toGatewayError(res)
    }
    return res
  }

  /** 拉取网关可用模型列表（OpenAI /v1/models 格式），密钥缺失时返回空数组 */
  async function fetchGatewayModels(): Promise<string[]> {
    if (!hasApiKey.value) return []
    try {
      const res = await gatewayFetch('/v1/models')
      const body: unknown = await res.json()
      const data = (body as { data?: unknown })?.data
      if (!Array.isArray(data)) return []
      return data
        .map((m) => {
          const id = (m as { id?: unknown })?.id
          return typeof id === 'string' ? id : ''
        })
        .filter(Boolean)
    } catch {
      return []
    }
  }

  return {
    apiKey,
    hasApiKey,
    setApiKey,
    authHeaders,
    gatewayFetch,
    fetchGatewayModels,
    chatModel,
    imageModel,
    videoModel,
  }
}
