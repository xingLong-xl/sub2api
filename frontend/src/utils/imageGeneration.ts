export type ImageModelFamily = 'openai' | 'grok' | 'gemini' | 'generic'
export type ImageMode = 'generate' | 'edit'

export const MAX_IMAGE_BATCH = 10
export const MAX_IMAGE_UPLOAD_BYTES = 20 * 1024 * 1024

export interface ImageSizePreset {
  value: string
  label: string
  ratio: string
  resolution: string
}

/** 尺寸同时是「比例 + 分辨率」的可视化选择；真正发给上游的字段按模型族决定。 */
export const IMAGE_SIZE_PRESETS: ImageSizePreset[] = [
  { value: 'auto', label: 'Auto', ratio: '', resolution: '' },
  { value: '1024x1024', label: '1024 × 1024 · 1:1 · 1K', ratio: '1:1', resolution: '1K' },
  { value: '1536x1024', label: '1536 × 1024 · 3:2', ratio: '3:2', resolution: '2K' },
  { value: '1024x1536', label: '1024 × 1536 · 2:3', ratio: '2:3', resolution: '2K' },
  { value: '2048x2048', label: '2048 × 2048 · 1:1 · 2K', ratio: '1:1', resolution: '2K' },
  { value: '2048x1152', label: '2048 × 1152 · 16:9 · 2K', ratio: '16:9', resolution: '2K' },
  { value: '1152x2048', label: '1152 × 2048 · 9:16 · 2K', ratio: '9:16', resolution: '2K' },
]

export function detectImageModelFamily(model: string): ImageModelFamily {
  const normalized = model.trim().toLowerCase()
  if (normalized.startsWith('gpt-image-') || normalized.startsWith('dall-e-')) return 'openai'
  if (normalized.startsWith('grok-imagine')) return 'grok'
  if (normalized.startsWith('gemini-') && (normalized.endsWith('-image') || normalized.includes('-image-'))) return 'gemini'
  return 'generic'
}

/** 多提示词用英文分号分隔；空片段会被忽略，避免一次发送空 prompt。 */
export function splitImagePrompts(prompt: string): string[] {
  return prompt
    .split(';')
    .map((part) => part.trim())
    .filter(Boolean)
}

export function imageSizePreset(value: string): ImageSizePreset {
  return IMAGE_SIZE_PRESETS.find((item) => item.value === value) ?? IMAGE_SIZE_PRESETS[0]
}

export interface ImagePayloadOptions {
  family: ImageModelFamily
  model: string
  prompt: string
  count: number
  size: string
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

/**
 * 构造 /v1/images/generations 的 JSON payload。
 *
 * OpenAI 图片模型接受完整的原生参数；Grok/Gemini 只发送它们共有的基础字段，
 * 避免把 quality/background 等 OpenAI 专属字段发送给不支持它们的上游。
 */
export function buildImageGenerationPayload(options: ImagePayloadOptions): Record<string, unknown> {
  const payload: Record<string, unknown> = {
    model: options.model.trim(),
    prompt: options.prompt.trim(),
    n: Math.max(1, Math.min(MAX_IMAGE_BATCH, Math.floor(options.count) || 1)),
  }

  if (options.size && options.size !== 'auto') payload.size = options.size

  if (options.family === 'openai') {
    if (options.quality) payload.quality = options.quality
    if (options.style) payload.style = options.style
    if (options.responseFormat) payload.response_format = options.responseFormat
    if (options.background) payload.background = options.background
    if (options.outputFormat) payload.output_format = options.outputFormat
    if (options.moderation) payload.moderation = options.moderation
    const compressedFormat = options.outputFormat === 'jpeg' || options.outputFormat === 'webp'
    if (compressedFormat && options.outputCompression.trim()) {
      payload.output_compression = Number(options.outputCompression)
    }
    // partial_images 只对 stream=true 有意义；当前图像页暂未提供流式预览开关。
    if (options.inputFidelity) payload.input_fidelity = options.inputFidelity
  }

  return payload
}

/** 图片结果可能是 png/jpeg/webp；data URL 的 MIME 必须和真实字节一致。 */
export function imageMimeFromFormat(format: unknown): string {
  switch (String(format ?? '').toLowerCase()) {
    case 'jpeg':
    case 'jpg':
      return 'image/jpeg'
    case 'webp':
      return 'image/webp'
    default:
      return 'image/png'
  }
}

export function imageExtensionFromMime(mime: string): string {
  if (mime === 'image/jpeg') return 'jpg'
  if (mime === 'image/webp') return 'webp'
  return 'png'
}
