import { computed, ref, watch } from 'vue'
import { buildGatewayUrl } from '@/api/url'
import { keysAPI } from '@/api/keys'
import { userChannelsAPI, type UserSupportedModelPricing } from '@/api/channels'
import { useAuthStore } from '@/stores/auth'
import { toolChatPreferencesAPI } from '@/api/toolChat'
import type { ToolChatPreferences } from '@/api/toolChat'
import type { BillingMode } from '@/constants/channel'
import type { ApiKey } from '@/types'

/**
 * 工具页（聊天 / 图像 / 视频）共享的连接设置。
 *
 * - 网关 /v1/* 只认 `Authorization: Bearer sk-xxx`，不认登录会话 Cookie，
 *   所以工具页必须携带一把用户自己的密钥。
 * - 优先让用户从「我的密钥」里直接选择（免去手工粘贴）；选不到时也可手动输入。
 * - 密钥绑定分组，分组带平台属性（openai / grok / anthropic / gemini…），
 *   平台决定该密钥能用哪些能力（图像仅 openai/grok，视频仅 grok）。
 *
 * 模块级单例：三个工具页共享同一份密钥与各自的模型偏好。
 */

const STORAGE_KEYS = {
  apiKey: 'tools_api_key',
  apiKeyId: 'tools_api_key_id',
  chatModel: 'tools_chat_model',
  imageModel: 'tools_image_model',
  videoModel: 'tools_video_model',
  owner: 'tools_settings_owner',
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

function readNumberStorage(key: string): number | null {
  const raw = readStorage(key)
  if (!raw) return null
  const value = Number(raw)
  return Number.isFinite(value) && value > 0 ? value : null
}

const apiKey = ref(readStorage(STORAGE_KEYS.apiKey))
const selectedKeyId = ref<number | null>(readNumberStorage(STORAGE_KEYS.apiKeyId))
const chatModel = ref(readStorage(STORAGE_KEYS.chatModel))
const imageModel = ref(readStorage(STORAGE_KEYS.imageModel))
const videoModel = ref(readStorage(STORAGE_KEYS.videoModel))

let preferencesLoaded = false
let preferencesLoading = false
let preferencesRequestSeq = 0
let preferencesSaveTimer: ReturnType<typeof setTimeout> | null = null
let applyingRemotePreferences = false

watch(apiKey, (v) => writeStorage(STORAGE_KEYS.apiKey, v.trim()))
watch(selectedKeyId, (v) => writeStorage(STORAGE_KEYS.apiKeyId, v === null ? '' : String(v)))
watch(chatModel, (v) => writeStorage(STORAGE_KEYS.chatModel, v.trim()))
watch(imageModel, (v) => writeStorage(STORAGE_KEYS.imageModel, v.trim()))
watch(videoModel, (v) => writeStorage(STORAGE_KEYS.videoModel, v.trim()))
watch([selectedKeyId, chatModel, imageModel, videoModel], () => schedulePreferencesSave())

const userKeys = ref<ApiKey[]>([])
const keysLoading = ref(false)
const keysLoaded = ref(false)
const keysError = ref('')

const gatewayModels = ref<string[]>([])
const modelsLoading = ref(false)
const modelsLoaded = ref(false)
/** /v1/models 失败原因（常见：账户余额为 0 → 网关返回 403） */
const modelsError = ref('')

/** 每个分组支持的模型（来自用户侧渠道接口，不受余额限制） */
export interface CatalogModel {
  name: string
  capability: ToolCapability
  pricing?: UserSupportedModelPricing | null
}

export interface GroupCatalog {
  name: string
  platform: string
  models: CatalogModel[]
}

const groupCatalog = ref<Record<string, GroupCatalog>>({})
const catalogLoading = ref(false)
const catalogLoaded = ref(false)
let keysRequestSeq = 0
let catalogRequestSeq = 0
let keysReloadQueued = false
let catalogReloadQueued = false

/** 当前登录用户 ID；未登录/尚未恢复用户信息时返回空串。 */
function currentToolsOwnerId(): string {
  try {
    const id = useAuthStore().user?.id
    return id === undefined || id === null ? '' : String(id)
  } catch {
    return ''
  }
}

/** 切换账号时清空所有连接选择，避免沿用上一个账号的密钥和分组。 */
function resetToolsConnectionState(): void {
  keysRequestSeq++
  catalogRequestSeq++
  preferencesRequestSeq++
  keysReloadQueued = false
  catalogReloadQueued = false
  preferencesLoaded = false
  preferencesLoading = false
  if (preferencesSaveTimer) {
    clearTimeout(preferencesSaveTimer)
    preferencesSaveTimer = null
  }
  applyingRemotePreferences = false
  apiKey.value = ''
  selectedKeyId.value = null
  chatModel.value = ''
  imageModel.value = ''
  videoModel.value = ''
  gatewayModels.value = []
  modelsLoading.value = false
  modelsLoaded.value = false
  modelsError.value = ''
  userKeys.value = []
  keysLoading.value = false
  keysLoaded.value = false
  keysError.value = ''
  groupCatalog.value = {}
  catalogLoading.value = false
  catalogLoaded.value = false
}

/**
 * 工具页的密钥/分组选择是浏览器级单例，必须绑定登录账号。
 * 账号变化（含首次升级时没有 owner 的旧缓存）一律清空，绝不跨账号复用密钥。
 */
function ensureToolsOwner(): void {
  const owner = currentToolsOwnerId()
  if (!owner) return
  let storedOwner = ''
  try {
    storedOwner = localStorage.getItem(STORAGE_KEYS.owner) ?? ''
  } catch {
    // 忽略存储失败；下面仍会清空内存中的旧选择
  }
  if (storedOwner === owner) return
  resetToolsConnectionState()
  try {
    localStorage.setItem(STORAGE_KEYS.owner, owner)
  } catch {
    // 忽略存储失败
  }
}

/** 应用服务端偏好。只恢复 Key ID，明文密钥仍从当前账号的密钥列表解析。 */
function applyRemotePreferences(preferences: ToolChatPreferences): void {
  applyingRemotePreferences = true
  try {
    selectedKeyId.value = preferences.api_key_id
    const matched =
      preferences.api_key_id === null
        ? undefined
        : userKeys.value.find((key) => key.id === preferences.api_key_id)
    apiKey.value = matched?.key ?? ''
    chatModel.value = preferences.chat_model ?? ''
    imageModel.value = preferences.image_model ?? ''
    videoModel.value = preferences.video_model ?? ''
  } finally {
    applyingRemotePreferences = false
  }
}

function schedulePreferencesSave(): void {
  if (!preferencesLoaded || applyingRemotePreferences || !currentToolsOwnerId()) return
  if (preferencesSaveTimer) clearTimeout(preferencesSaveTimer)
  preferencesSaveTimer = setTimeout(() => {
    preferencesSaveTimer = null
    void savePreferences()
  }, 400)
}

async function savePreferences(): Promise<void> {
  const owner = currentToolsOwnerId()
  if (!owner || !preferencesLoaded) return
  const payload: ToolChatPreferences = {
    api_key_id: selectedKeyId.value,
    chat_model: chatModel.value.trim(),
    image_model: imageModel.value.trim(),
    video_model: videoModel.value.trim(),
  }
  try {
    await toolChatPreferencesAPI.save(payload)
  } catch {
    // 同步失败不阻塞本地使用；下次修改或刷新会重试。
  }
}

/**
 * 从服务端读取账号级偏好。服务端没有记录时，把当前浏览器里的旧选择作为初始值上传。
 * 读取失败时保留本地值，且不标记已加载，避免用本地值覆盖服务端。
 */
async function loadPreferences(force = false): Promise<void> {
  const owner = currentToolsOwnerId()
  if (!owner) return
  if (preferencesLoaded && !force) return
  if (preferencesLoading && !force) return
  const requestSeq = ++preferencesRequestSeq
  preferencesLoading = true
  try {
    const preferences = await toolChatPreferencesAPI.get()
    if (requestSeq !== preferencesRequestSeq || owner !== currentToolsOwnerId()) return
    preferencesLoaded = true
    if (preferences) {
      applyRemotePreferences(preferences)
    } else {
      // 新账号/首次升级：当前本地选择作为种子，随后同步到服务端。
      schedulePreferencesSave()
    }
  } catch {
    // 读取失败保持本地状态；不标记 loaded，后续页面加载会重试。
  } finally {
    if (requestSeq === preferencesRequestSeq) preferencesLoading = false
  }
}

/** 工具能力：与网关的三条业务链路一一对应 */
export type ToolCapability = 'chat' | 'image' | 'video'

/**
 * 模型能力推断。网关 /v1/models 只返回模型 id，没有能力字段，
 * 只能按命名约定归类；调用方始终提供「显示全部模型」的兜底开关。
 */
const VIDEO_MODEL_HINTS = ['video', 'seedance', 'sora', 'veo', 'kling', 'runway', 'hailuo', 'pika']
const IMAGE_MODEL_HINTS = ['image', 'dall-e', 'imagen', 'flux', 'midjourney', 'seedream', 'seededit']

export function inferModelCapability(modelID: string): ToolCapability {
  const name = modelID.trim().toLowerCase()
  if (!name) return 'chat'
  if (VIDEO_MODEL_HINTS.some((hint) => name.includes(hint))) return 'video'
  if (IMAGE_MODEL_HINTS.some((hint) => name.includes(hint))) return 'image'
  return 'chat'
}

/**
 * 用计费模式判定模型能力（比模型名猜测准确）。
 * token/per_request 仍可能是图像模型（按张计费），因此再退回名称判断。
 */
function capabilityFromBillingMode(mode: BillingMode | undefined, modelID: string): ToolCapability {
  if (mode === 'image') return 'image'
  if (mode === 'video') return 'video'
  return inferModelCapability(modelID)
}

/** 平台能力：图像支持 openai/grok/gemini，视频仅 grok（合成分组可路由到 grok） */
const IMAGE_PLATFORMS = ['openai', 'grok', 'gemini']
const VIDEO_PLATFORMS = ['grok', 'composite']

/** 网关错误：带 HTTP 状态码与业务码（code），视图层据此做本地化提示 */
export class GatewayError extends Error {
  readonly status: number
  /** 后端返回的业务错误码，例如 INSUFFICIENT_BALANCE、API_KEY_EXPIRED */
  readonly code: string

  constructor(message: string, status: number, code = '') {
    super(message)
    this.name = 'GatewayError'
    this.status = status
    this.code = code
  }
}

/** 从网关/上游错误响应体中尽力提取可读信息 */
async function toGatewayError(res: Response): Promise<GatewayError> {
  let message = `HTTP ${res.status}`
  let code = ''
  try {
    const body: unknown = await res.json()
    if (body && typeof body === 'object') {
      const rawCode = (body as Record<string, unknown>).code
      if (typeof rawCode === 'string') code = rawCode
      const err = (body as Record<string, unknown>).error
      if (typeof err === 'string' && err) {
        message = err
      } else if (err && typeof err === 'object') {
        const msg = (err as { message?: unknown }).message
        if (typeof msg === 'string' && msg) message = msg
      } else {
        const msg = (body as { message?: unknown }).message
        if (typeof msg === 'string' && msg) message = msg
      }
    }
  } catch {
    // 非 JSON 响应，保留 HTTP 状态码
  }
  return new GatewayError(message, res.status, code)
}

const selectedKey = computed<ApiKey | null>(() => {
  if (selectedKeyId.value === null) return null
  return userKeys.value.find((k) => k.id === selectedKeyId.value) ?? null
})

const selectedKeyPlatform = computed(() => selectedKey.value?.group?.platform ?? '')

const selectedGroupID = computed<number | null>(() => selectedKey.value?.group_id ?? null)

/** 当前分组在渠道目录里的模型清单；目录未覆盖该分组时返回空数组 */
function catalogModelsForGroup(groupID: number | null, tool: ToolCapability): string[] {
  if (groupID === null) return []
  const entry = groupCatalog.value[String(groupID)]
  if (!entry) return []
  return entry.models.filter((model) => model.capability === tool).map((model) => model.name)
}

/** 判断给定平台是否支持某个能力；未知（空平台）时返回 null 表示「不确定」 */
export function platformSupportsTool(platform: string, tool: ToolCapability): boolean | null {
  if (!platform) return null
  if (tool === 'image') return IMAGE_PLATFORMS.includes(platform)
  if (tool === 'video') return VIDEO_PLATFORMS.includes(platform)
  return true
}

/** 当前密钥是否支持某个能力；未知（手动输入/无分组）时返回 null 表示「不确定」 */
function platformSupports(tool: ToolCapability): boolean | null {
  return platformSupportsTool(selectedKeyPlatform.value, tool)
}

/** 按能力过滤模型列表（未知能力的模型归为 chat，避免图像/视频页混入聊天模型） */
function modelsForCapability(tool: ToolCapability, models: string[]): string[] {
  return models.filter((model) => inferModelCapability(model) === tool)
}

export function useToolsSettings() {
  const authStore = useAuthStore()
  ensureToolsOwner()
  void loadPreferences()
  watch(
    () => authStore.user?.id,
    () => {
      if (!currentToolsOwnerId()) {
        resetToolsConnectionState()
        return
      }
      ensureToolsOwner()
      void loadPreferences(true)
      void loadUserKeys(true)
      void loadModelCatalog(true)
    },
  )

  /** 手动输入密钥（同时清掉「从列表选择」的状态） */
  function setApiKey(key: string): void {
    apiKey.value = key.trim()
    selectedKeyId.value = null
  }

  /** 从「我的密钥」中选一把；密钥明文从列表接口取回，避免用户重复粘贴 */
  function selectKey(keyID: number | null): void {
    selectedKeyId.value = keyID
    if (keyID === null) return
    const matched = userKeys.value.find((k) => k.id === keyID)
    if (matched) apiKey.value = matched.key
  }

  /** 清空当前选择：不选分组时不携带任何密钥，也不展示任何模型 */
  function clearSelection(): void {
    selectedKeyId.value = null
    apiKey.value = ''
    gatewayModels.value = []
    modelsLoaded.value = false
    modelsError.value = ''
  }

  async function loadUserKeys(force = false): Promise<ApiKey[]> {
    ensureToolsOwner()
    void loadPreferences()
    if (keysLoading.value) {
      if (force) keysReloadQueued = true
      return userKeys.value
    }
    if (keysLoaded.value && !force) return userKeys.value
    const requestSeq = ++keysRequestSeq
    keysLoading.value = true
    keysError.value = ''
    try {
      const page = await keysAPI.list(1, 100, { sort_by: 'created_at', sort_order: 'desc' })
      // 账号在请求期间切换：丢弃旧账号的响应，避免把上一个账号的密钥写回来。
      if (requestSeq !== keysRequestSeq) return userKeys.value
      const items = Array.isArray(page.items) ? page.items : []
      userKeys.value = items.filter((k) => k.status === 'active')
      keysLoaded.value = true
      // 恢复上次选择：密钥可能被轮换，以服务端最新值为准
      if (selectedKeyId.value !== null) {
        const matched = userKeys.value.find((k) => k.id === selectedKeyId.value)
        if (matched) {
          apiKey.value = matched.key
        } else {
          // 选中的密钥不属于当前账号（账号切换/密钥被删除）：连明文一起清掉。
          clearSelection()
        }
      } else if (apiKey.value.trim()) {
        // 旧版本只存了密钥明文，没有 id：按密钥值反查，让分组下拉能正确回显
        const matched = userKeys.value.find((k) => k.key === apiKey.value.trim())
        if (matched) selectedKeyId.value = matched.id
      }
    } catch (err) {
      if (requestSeq === keysRequestSeq) {
        keysError.value = err instanceof Error ? err.message : String(err)
      }
    } finally {
      if (requestSeq === keysRequestSeq) {
        keysLoading.value = false
        if (keysReloadQueued) {
          keysReloadQueued = false
          void loadUserKeys(true)
        }
      }
    }
    return userKeys.value
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
    const headers = new Headers(init.headers)
    // FormData 必须让浏览器自己带 multipart boundary，不能强制 JSON。
    const isFormData = typeof FormData !== 'undefined' && init.body instanceof FormData
    if (isFormData) {
      headers.delete('Content-Type')
    } else if (!headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }
    for (const [key, value] of Object.entries(authHeaders())) {
      headers.set(key, value)
    }
    const res = await fetch(buildGatewayUrl(path), { ...init, headers })
    if (!res.ok) {
      throw await toGatewayError(res)
    }
    return res
  }

  /**
   * 拉取「分组 → 模型」目录（用户侧渠道接口）。
   * 该接口走登录态、不受网关余额校验影响，且带 billing_mode，
   * 能准确区分聊天 / 图像 / 视频模型，比解析 /v1/models 更可靠。
   */
  async function loadModelCatalog(force = false): Promise<void> {
    if (catalogLoading.value) {
      if (force) catalogReloadQueued = true
      return
    }
    if (catalogLoaded.value && !force) return
    const requestSeq = ++catalogRequestSeq
    catalogLoading.value = true
    try {
      const channels = await userChannelsAPI.getAvailable()
      // 账号在请求期间切换：旧账号的分组目录不能留给新账号。
      if (requestSeq !== catalogRequestSeq) return
      const map: Record<string, GroupCatalog> = {}
      for (const channel of channels ?? []) {
        for (const section of channel.platforms ?? []) {
          const models: CatalogModel[] = (section.supported_models ?? []).map((model) => ({
            name: model.name,
            capability: capabilityFromBillingMode(model.pricing?.billing_mode, model.name),
            pricing: model.pricing,
          }))
          for (const group of section.groups ?? []) {
            const key = String(group.id)
            const existing = map[key]
            if (!existing || existing.models.length === 0) {
              map[key] = {
                name: group.name,
                platform: group.platform || section.platform,
                models,
              }
            }
          }
        }
      }
      groupCatalog.value = map
      catalogLoaded.value = true
    } catch {
      // 渠道接口失败时保持空目录，由调用方退回 /v1/models 或手动输入
    } finally {
      if (requestSeq === catalogRequestSeq) {
        catalogLoading.value = false
        if (catalogReloadQueued) {
          catalogReloadQueued = false
          void loadModelCatalog(true)
        }
      }
    }
  }

  /** 拉取网关可用模型列表（OpenAI /v1/models 格式），密钥缺失时返回空数组 */
  async function fetchGatewayModels(force = false): Promise<string[]> {
    if (!hasApiKey.value) return []
    if (modelsLoaded.value && !force) return gatewayModels.value
    modelsLoading.value = true
    try {
      const res = await gatewayFetch('/v1/models')
      const body: unknown = await res.json()
      const data = (body as { data?: unknown })?.data
      if (!Array.isArray(data)) return gatewayModels.value
      modelsError.value = ''
      gatewayModels.value = data
        .map((m) => {
          const id = (m as { id?: unknown })?.id
          return typeof id === 'string' ? id : ''
        })
        .filter(Boolean)
      modelsLoaded.value = true
    } catch (err) {
      // 模型列表拉取失败不影响手工填写，但要记录原因供界面提示
      modelsError.value = err instanceof Error ? err.message : ''
    } finally {
      modelsLoading.value = false
    }
    return gatewayModels.value
  }

  return {
    apiKey,
    hasApiKey,
    setApiKey,
    selectKey,
    clearSelection,
    selectedKeyId,
    selectedKey,
    selectedKeyPlatform,
    selectedGroupID,
    groupCatalog,
    catalogLoading,
    loadModelCatalog,
    catalogModelsForGroup,
    userKeys,
    keysLoading,
    keysError,
    loadUserKeys,
    authHeaders,
    gatewayFetch,
    fetchGatewayModels,
    gatewayModels,
    modelsError,
    modelsLoading,
    platformSupports,
    platformSupportsTool,
    modelsForCapability,
    chatModel,
    imageModel,
    videoModel,
  }
}