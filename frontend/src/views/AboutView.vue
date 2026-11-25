<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-primary-50/20 to-gray-100 text-gray-900 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950 dark:text-white">
    <header class="border-b border-gray-200/70 bg-white/75 backdrop-blur dark:border-dark-800 dark:bg-dark-900/75">
      <div class="mx-auto flex max-w-5xl items-center justify-between px-6 py-4">
        <router-link to="/home" class="flex items-center gap-3">
          <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-9 w-9 rounded-lg object-contain" />
          <span class="text-base font-semibold tracking-tight">{{ siteName }}</span>
        </router-link>
        <router-link
          to="/home"
          class="rounded-lg px-3 py-2 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
        >
          返回首页
        </router-link>
      </div>
    </header>

    <main class="mx-auto max-w-5xl px-6 py-14 sm:py-20">
      <section class="mx-auto max-w-3xl text-center">
        <span class="inline-flex rounded-full border border-primary-200 bg-primary-50 px-3 py-1 text-xs font-semibold text-primary-700 dark:border-primary-800/70 dark:bg-primary-900/25 dark:text-primary-300">
          关于 {{ siteName }}
        </span>
        <h1 class="mt-5 text-4xl font-bold tracking-tight md:text-5xl">
          让开发者更简单地接入多模型能力
        </h1>
        <p class="mx-auto mt-6 max-w-2xl text-base leading-8 text-gray-600 dark:text-dark-300">
          TokenPilot 是面向开发者、独立产品和团队的多模型统一接入平台。我们把 API Base、模型价格、API Key、用量记录和常见开发工具配置放到同一个入口，让 Claude Code、Codex 与 OpenAI 兼容客户端可以更快、更稳地接入你的工作流。
        </p>
        <p class="mt-3 text-sm text-gray-500 dark:text-dark-400">
          TokenPilot 由 Hengshi Boyuan International Co., Limited 运营。
        </p>
        <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
          <router-link to="/docs" class="btn btn-primary px-6 py-3 text-sm shadow-lg shadow-primary-500/25">
            查看配置文档
          </router-link>
          <router-link
            to="/model-plaza"
            class="rounded-lg border border-gray-200 bg-white px-6 py-3 text-sm font-medium text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-300"
          >
            查看模型价格
          </router-link>
        </div>
      </section>

      <section class="mt-16">
        <div class="flex flex-wrap items-center justify-center gap-2">
          <span v-for="tag in focusTags" :key="tag" class="rounded-full border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-600 dark:border-dark-700 dark:bg-dark-800/70 dark:text-dark-300">
            {{ tag }}
          </span>
        </div>
        <div class="mt-8 grid gap-5 md:grid-cols-2">
          <article
            v-for="(feature, index) in features"
            :key="feature.title"
            class="rounded-2xl border border-gray-200/70 bg-white/80 p-6 transition-shadow hover:shadow-lg hover:shadow-primary-500/5 dark:border-dark-700/70 dark:bg-dark-800/60"
          >
            <span class="text-xs font-semibold tabular-nums text-primary-500">0{{ index + 1 }}</span>
            <h2 class="mt-3 text-base font-semibold">{{ feature.title }}</h2>
            <p class="mt-2 text-sm leading-7 text-gray-600 dark:text-dark-300">{{ feature.description }}</p>
          </article>
        </div>
      </section>

      <section class="mt-14 space-y-6">
        <div class="rounded-2xl border border-gray-200/70 bg-white/80 p-6 dark:border-dark-700/70 dark:bg-dark-800/60">
          <h2 class="text-lg font-semibold">运营主体</h2>
          <dl class="mt-5 space-y-4 text-sm">
            <div class="flex items-start justify-between gap-4 border-b border-gray-100 pb-3 dark:border-dark-700">
              <dt class="text-gray-400 dark:text-dark-500">产品名称</dt>
              <dd class="text-right font-medium">TokenPilot</dd>
            </div>
            <div class="flex items-start justify-between gap-4 border-b border-gray-100 pb-3 dark:border-dark-700">
              <dt class="text-gray-400 dark:text-dark-500">运营主体</dt>
              <dd class="text-right font-medium">Hengshi Boyuan International Co., Limited</dd>
            </div>
            <div class="flex items-start justify-between gap-4 border-b border-gray-100 pb-3 dark:border-dark-700">
              <dt class="text-gray-400 dark:text-dark-500">注册地</dt>
              <dd class="text-right font-medium">Hong Kong</dd>
            </div>
            <div class="flex items-start justify-between gap-4">
              <dt class="text-gray-400 dark:text-dark-500">成立日期</dt>
              <dd class="text-right font-medium">2026-07-03</dd>
            </div>
          </dl>
          <p class="mt-5 text-xs leading-6 text-gray-500 dark:text-dark-400">
            Hengshi Boyuan International Co., Limited 负责 TokenPilot 的服务运营、账户、计费与支持。
          </p>
        </div>

        <div class="rounded-2xl border border-primary-100 bg-gradient-to-br from-primary-50/80 to-white p-6 dark:border-primary-900/40 dark:from-primary-900/10 dark:to-dark-800">
          <h2 class="text-lg font-semibold">我们的原则</h2>
          <p class="mt-4 text-sm leading-8 text-gray-600 dark:text-dark-300">
            我们希望 TokenPilot 不只是一个转发入口，而是一个能让开发者放心长期使用的基础服务。因此，文档要讲清楚，价格要看得懂，Key 和用量要能追踪，配置错误要能快速定位。
          </p>
          <p class="mt-4 text-sm leading-8 text-gray-600 dark:text-dark-300">
            这些看起来朴素的事情，正是把模型能力接入真实项目时最重要的部分。
          </p>
          <div class="mt-6 flex flex-wrap gap-3 text-sm">
            <router-link to="/privacy-policy" class="text-primary-600 hover:underline dark:text-primary-400">隐私政策</router-link>
            <router-link to="/user-agreement" class="text-primary-600 hover:underline dark:text-primary-400">用户协议</router-link>
          </div>
        </div>
      </section>
    </main>

    <footer class="relative z-10 px-6">
      <div class="mx-auto max-w-6xl">
        <div class="border-y border-solid border-gray-200 py-12 dark:border-dark-700">
          <div class="grid gap-10 lg:grid-cols-12 lg:gap-8">
            <div class="lg:col-span-5">
              <div class="flex items-center gap-3">
                <div class="h-9 w-9 overflow-hidden rounded-lg">
                  <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
                </div>
                <p class="text-xl font-semibold tracking-tight text-gray-900 dark:text-white">TokenPilot</p>
              </div>
              <p class="mt-3 max-w-sm text-sm leading-6 text-gray-500 dark:text-dark-400">
                统一入口 · 高速稳定 · 可观测可控。为你的产品与团队接入多模型能力。
              </p>
            </div>

            <div class="grid grid-cols-2 gap-8 sm:grid-cols-3 lg:col-span-7">
              <div>
                <h3 class="text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-dark-500">产品</h3>
                <ul class="mt-4 space-y-2.5 text-sm">
                  <li><router-link to="/model-plaza" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">模型价格</router-link></li>
                  <li><router-link to="/docs" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">配置文档</router-link></li>
                  <li><router-link to="/dashboard" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">控制台</router-link></li>
                  <li><router-link to="/keys" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">API Key</router-link></li>
                </ul>
              </div>

              <div>
                <h3 class="text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-dark-500">资源</h3>
                <ul class="mt-4 space-y-2.5 text-sm">
                  <li><a href="https://www.anthropic.com/claude-code" target="_blank" rel="noopener noreferrer" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">Claude Code</a></li>
                  <li><a href="https://openai.com/codex/" target="_blank" rel="noopener noreferrer" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">Codex</a></li>
                  <li><a href="https://platform.openai.com/docs/api-reference" target="_blank" rel="noopener noreferrer" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">OpenAI 兼容 API</a></li>
                </ul>
              </div>

              <div>
                <h3 class="text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-dark-500">公司</h3>
                <ul class="mt-4 space-y-2.5 text-sm">
                  <li><router-link to="/about" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">关于我们</router-link></li>
                  <li><span class="text-gray-500 dark:text-dark-400">Hengshi Boyuan International Co., Limited</span></li>
                  <li><router-link to="/privacy-policy" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">隐私政策</router-link></li>
                  <li><router-link to="/user-agreement" class="text-gray-600 transition-colors hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400">用户协议</router-link></li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        <p class="py-6 text-center text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. 保留所有权利。
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'TokenPilot')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const currentYear = new Date().getFullYear()

const focusTags = ['统一入口', '高速稳定', '价格透明', '可观测可控']

const features = [
  {
    title: '统一多模型入口',
    description: '通过统一 API Base 接入 Claude、GPT 等模型能力，减少在不同平台之间切换配置的成本。',
  },
  {
    title: '面向开发工具优化',
    description: '围绕 Claude Code、Codex 和 OpenAI 兼容客户端整理配置路径、示例和排错流程。',
  },
  {
    title: '价格透明可追踪',
    description: '模型价格、用量记录和 Key 管理放在同一个控制台里，方便团队持续核对成本。',
  },
  {
    title: '稳定、可控、可观测',
    description: '关注接入稳定性、权限边界和调用可见性，让模型能力更适合长期集成到产品和工作流中。',
  },
]
</script>
