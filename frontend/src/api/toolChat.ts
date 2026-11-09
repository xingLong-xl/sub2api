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
