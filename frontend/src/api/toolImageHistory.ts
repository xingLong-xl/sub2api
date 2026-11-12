import { apiClient } from './client'

export interface ToolImageHistoryImage {
  src: string
  isDataUrl: boolean
  mime?: string
  width?: number
  height?: number
  bytes?: number
}

export interface ToolImageHistoryParams {
  mode: 'generate' | 'edit'
  size: string
  count: number
  quality: string
  style: string
  responseFormat: string
  background: string
  outputFormat: string
  moderation: string
  outputCompression: string
  partialImages: string
  inputFidelity: string
}

export interface ToolImageHistoryItem {
  id: string
  prompt: string
  model: string
  mode: 'generate' | 'edit'
  images: ToolImageHistoryImage[]
  params: ToolImageHistoryParams
  duration_ms?: number
  cost_estimate?: number
  createdAt: number
}

interface ToolImageHistoryResponse {
  items: ToolImageHistoryItem[]
  max_items?: number
  max_bytes?: number
  retention_days?: number
}

export interface ToolImageHistoryResult {
  items: ToolImageHistoryItem[]
  maxItems: number
  maxBytes: number
  retentionDays: number
}

interface ToolImageHistorySaveResponse {
  saved: boolean
  max_items?: number
  retention_days?: number
}

/** 图片生成历史按登录用户保存在服务端，跨浏览器共享。 */
export const toolImageHistoryAPI = {
  async list(): Promise<ToolImageHistoryResult> {
    const { data } = await apiClient.get<ToolImageHistoryResponse>('/tools/image/history', {
      timeout: 60000,
    })
    return {
      items: Array.isArray(data.items) ? data.items : [],
      maxItems: typeof data.max_items === 'number' && data.max_items > 0 ? data.max_items : 20,
      maxBytes: typeof data.max_bytes === 'number' && data.max_bytes > 0 ? data.max_bytes : 48 * 1024 * 1024,
      retentionDays:
        typeof data.retention_days === 'number' && data.retention_days > 0
          ? data.retention_days
          : 7,
    }
  },

  async append(item: ToolImageHistoryItem): Promise<ToolImageHistorySaveResponse> {
    const { data } = await apiClient.post<ToolImageHistorySaveResponse>('/tools/image/history', item, {
      timeout: 60000,
    })
    return data
  },

  async update(item: ToolImageHistoryItem): Promise<ToolImageHistorySaveResponse> {
    const { data } = await apiClient.put<ToolImageHistorySaveResponse>(
      `/tools/image/history/${encodeURIComponent(item.id)}`,
      item,
      { timeout: 60000 },
    )
    return data
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/tools/image/history/${encodeURIComponent(id)}`, { timeout: 60000 })
  },
}
