import { apiClient } from './client'

export interface ToolChatMessage {
  role: 'user' | 'assistant'
  content: string
  createdAt?: number
}

export interface ToolChatSession {
  id: string
  title: string
  createdAt: number
  messages: ToolChatMessage[]
}

interface ToolChatSessionsResponse {
  sessions: ToolChatSession[]
  context_limit?: number
  retention_days?: number
}

export interface ToolChatSessionsResult {
  sessions: ToolChatSession[]
  contextLimit: number
  retentionDays: number
}

export interface ToolChatPreferences {
  api_key_id: number | null
  chat_model: string
  image_model: string
  video_model: string
}

interface ToolChatPreferencesResponse {
  preferences: ToolChatPreferences | null
}

/** 账号级工具连接偏好，用于跨浏览器同步分组/模型选择；不包含密钥明文。 */
export const toolChatPreferencesAPI = {
  async get(): Promise<ToolChatPreferences | null> {
    const { data } = await apiClient.get<ToolChatPreferencesResponse>('/tools/chat/preferences')
    return data.preferences ?? null
  },

  async save(preferences: ToolChatPreferences): Promise<void> {
    await apiClient.put('/tools/chat/preferences', preferences)
  },
}

/** 服务端聊天历史。历史按登录用户隔离，保存在 Redis。 */
export const toolChatAPI = {
  async list(): Promise<ToolChatSessionsResult> {
    const { data } = await apiClient.get<ToolChatSessionsResponse>('/tools/chat/sessions')
    return {
      sessions: Array.isArray(data.sessions) ? data.sessions : [],
      contextLimit: typeof data.context_limit === 'number' && data.context_limit > 0 ? data.context_limit : 20,
      retentionDays: typeof data.retention_days === 'number' && data.retention_days > 0 ? data.retention_days : 7,
    }
  },

  async save(sessions: ToolChatSession[]): Promise<void> {
    await apiClient.put('/tools/chat/sessions', { sessions })
  },
}
