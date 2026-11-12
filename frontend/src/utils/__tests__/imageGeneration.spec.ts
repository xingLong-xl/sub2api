import { describe, expect, it } from 'vitest'
import {
  MAX_IMAGE_BATCH,
  buildImageGenerationPayload,
  detectImageModelFamily,
  imageExtensionFromMime,
  imageMimeFromFormat,
  imageSizePreset,
  splitImagePrompts,
} from '@/utils/imageGeneration'

describe('imageGeneration', () => {
  it('detects the image model family', () => {
    expect(detectImageModelFamily('gpt-image-2.5-flare')).toBe('openai')
    expect(detectImageModelFamily('gpt-image-2.6-flare')).toBe('openai')
    expect(detectImageModelFamily('dall-e-3')).toBe('openai')
    expect(detectImageModelFamily('grok-imagine-image-quality')).toBe('grok')
    expect(detectImageModelFamily('gemini-3.1-flash-image')).toBe('gemini')
    expect(detectImageModelFamily('gemini-3-pro-image-preview')).toBe('gemini')
    expect(detectImageModelFamily('custom-image-model')).toBe('generic')
  })

  it('splits semicolon prompts and ignores empty segments', () => {
    expect(splitImagePrompts(' a cat ; ; a dog;  ')).toEqual(['a cat', 'a dog'])
    expect(splitImagePrompts('   ')).toEqual([])
  })

  it('keeps OpenAI-only parameters away from Grok and Gemini', () => {
    const common = {
      model: 'gpt-image-2',
      prompt: 'draw',
      count: 2,
      size: '2048x2048',
      quality: 'high',
      style: 'vivid',
      responseFormat: 'b64_json',
      background: 'opaque',
      outputFormat: 'webp',
      moderation: 'low',
      outputCompression: '80',
      partialImages: '2',
      inputFidelity: 'high',
    } as const

    const openai = buildImageGenerationPayload({ ...common, family: 'openai' })
    expect(openai).toMatchObject({
      model: 'gpt-image-2',
      n: 2,
      size: '2048x2048',
      quality: 'high',
      output_format: 'webp',
      output_compression: 80,
      input_fidelity: 'high',
    })

    expect(openai).not.toHaveProperty('partial_images')

    const grok = buildImageGenerationPayload({ ...common, family: 'grok' })
    expect(grok).toMatchObject({ model: 'gpt-image-2', n: 2, size: '2048x2048' })
    expect(grok).not.toHaveProperty('quality')
    expect(grok).not.toHaveProperty('background')
    expect(grok).not.toHaveProperty('output_format')

    const gemini = buildImageGenerationPayload({ ...common, family: 'gemini' })
    expect(gemini).not.toHaveProperty('quality')
    expect(gemini).not.toHaveProperty('output_compression')
  })

  it('clamps count to the batch limit and omits auto size', () => {
    const payload = buildImageGenerationPayload({
      family: 'openai',
      model: 'gpt-image-2',
      prompt: 'draw',
      count: 99,
      size: 'auto',
      quality: '',
      style: '',
      responseFormat: '',
      background: '',
      outputFormat: '',
      moderation: '',
      outputCompression: '',
      partialImages: '',
      inputFidelity: '',
    })
    expect(payload.n).toBe(MAX_IMAGE_BATCH)
    expect(payload).not.toHaveProperty('size')
  })

  it('maps output formats and size presets', () => {
    expect(imageMimeFromFormat('jpeg')).toBe('image/jpeg')
    expect(imageMimeFromFormat('webp')).toBe('image/webp')
    expect(imageExtensionFromMime('image/jpeg')).toBe('jpg')
    expect(imageSizePreset('2048x2048')).toMatchObject({ ratio: '1:1', resolution: '2K' })
  })
})
