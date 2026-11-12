<template>
  <AppLayout>
    <div class="-m-[11px] flex h-[calc(100dvh-77px)] flex-col overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800 md:-m-[19px] lg:-m-[27px]">
      <div
        class="grid min-h-0 flex-1 grid-cols-1 divide-y divide-gray-200 overflow-y-auto dark:divide-dark-700 xl:divide-x xl:divide-y-0 xl:overflow-hidden"
        :class="paramsCollapsed
          ? (historyCollapsed ? 'xl:grid-cols-[0px_minmax(0,1fr)_0px]' : 'xl:grid-cols-[0px_minmax(0,1fr)_300px]')
          : (historyCollapsed ? 'xl:grid-cols-[380px_minmax(0,1fr)_0px]' : 'xl:grid-cols-[380px_minmax(0,1fr)_300px]')"
      >
        <!-- 第一栏：生成参数 -->
        <ToolsPanel v-if="!paramsCollapsed" class="xl:col-start-1" :title="t('tools.image.params')" icon="cog">
          <template #actions>
            <button
              type="button"
              class="btn btn-ghost !px-2 !py-1 text-xs text-gray-500"
              :disabled="generating"
              @click="restoreDefaults"
            >
              {{ t('tools.image.restoreDefaults') }}
            </button>
          </template>

          <fieldset :disabled="generating" class="min-w-0 space-y-3 border-0 p-0">
            <ToolsConnectionBar
              v-model="imageModel"
              tool="image"
              full-width
              class="border-b border-gray-100 pb-3 dark:border-dark-700"
            />

            <div class="grid grid-cols-2 gap-1 rounded-lg bg-gray-100 p-1 dark:bg-dark-700/60">
              <button
                type="button"
                class="rounded-md px-2 py-1.5 text-xs font-medium transition-colors"
                :class="mode === 'generate'
                  ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-800 dark:text-dark-100'
                  : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'"
                @click="mode = 'generate'"
              >
                {{ t('tools.image.modeGenerate') }}
              </button>
              <button
                type="button"
                class="rounded-md px-2 py-1.5 text-xs font-medium transition-colors"
                :class="mode === 'edit'
                  ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-800 dark:text-dark-100'
                  : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'"
                @click="mode = 'edit'"
              >
                {{ t('tools.image.modeEdit') }}
              </button>
            </div>

            <p class="-mt-1 text-[11px] leading-4 text-gray-400 dark:text-dark-500">
              {{ t('tools.image.modeHint') }}
            </p>

            <div>
              <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-image-prompt">
                {{ t('tools.image.prompt') }}
              </label>
              <textarea
                id="tools-image-prompt"
                v-model="prompt"
                rows="4"
                :placeholder="mode === 'edit' ? t('tools.image.editPromptPlaceholder') : t('tools.image.promptPlaceholder')"
                class="input resize-none text-sm"
              ></textarea>
            </div>

            <div v-if="mode === 'edit'" class="grid grid-cols-2 gap-2">
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300">
                  {{ t('tools.image.sourceImage') }}
                </label>
                <input
                  ref="sourceInput"
                  type="file"
                  accept="image/png,image/jpeg,image/webp"
                  class="hidden"
                  @change="onFileChange('source', $event)"
                />
                <div class="group relative overflow-hidden rounded-lg">
                  <button
                    type="button"
                    class="flex h-24 w-full flex-col items-center justify-center gap-1 rounded-lg border border-dashed border-gray-300 bg-gray-50/60 text-[11px] text-gray-500 transition-colors hover:border-primary-300 hover:text-primary-600 dark:border-dark-600 dark:bg-dark-700/30 dark:text-dark-400 dark:hover:border-primary-600 dark:hover:text-primary-400"
                    @click="sourceInput?.click()"
                  >
                    <img v-if="sourcePreview" :src="sourcePreview" alt="" class="absolute inset-0 h-full w-full object-cover" />
                    <template v-else>
                      <Icon name="upload" size="sm" />
                      <span>{{ t('tools.image.chooseImage') }}</span>
                    </template>
                  </button>
                  <button
                    v-if="sourcePreview"
                    type="button"
                    class="absolute right-1.5 top-1.5 inline-flex h-5 w-5 items-center justify-center rounded-full bg-black/60 text-white opacity-0 transition-opacity group-hover:opacity-100"
                    :title="t('tools.image.removeImage')"
                    @click.stop="clearFile('source')"
                  >
                    <Icon name="x" size="xs" />
                  </button>
                </div>
              </div>

              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300">
                  {{ t('tools.image.maskImage') }}
                  <span class="font-normal text-gray-400">· {{ t('tools.image.optional') }}</span>
                </label>
                <input
                  ref="maskInput"
                  type="file"
                  accept="image/png,image/jpeg,image/webp"
                  class="hidden"
                  @change="onFileChange('mask', $event)"
                />
                <div class="group relative overflow-hidden rounded-lg">
                  <button
                    type="button"
                    class="flex h-24 w-full flex-col items-center justify-center gap-1 rounded-lg border border-dashed border-gray-300 bg-gray-50/60 text-[11px] text-gray-500 transition-colors hover:border-primary-300 hover:text-primary-600 dark:border-dark-600 dark:bg-dark-700/30 dark:text-dark-400 dark:hover:border-primary-600 dark:hover:text-primary-400"
                    @click="maskInput?.click()"
                  >
                    <img v-if="maskPreview" :src="maskPreview" alt="" class="absolute inset-0 h-full w-full object-cover" />
                    <template v-else>
                      <Icon name="upload" size="sm" />
                      <span>{{ t('tools.image.chooseMask') }}</span>
                    </template>
                  </button>
                  <button
                    v-if="maskPreview"
                    type="button"
                    class="absolute right-1.5 top-1.5 inline-flex h-5 w-5 items-center justify-center rounded-full bg-black/60 text-white opacity-0 transition-opacity group-hover:opacity-100"
                    :title="t('tools.image.removeImage')"
                    @click.stop="clearFile('mask')"
                  >
                    <Icon name="x" size="xs" />
                  </button>
                </div>
              </div>
            </div>

            <p v-if="mode === 'edit'" class="text-[11px] leading-4 text-gray-400 dark:text-dark-500">
              {{ t('tools.image.imageUploadHint') }}
            </p>

            <div class="grid grid-cols-[minmax(0,1fr)_92px] gap-2">
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-image-size">
                  {{ t('tools.image.size') }}
                </label>
                <select id="tools-image-size" v-model="size" class="input !py-1.5 text-sm" :title="sizeLabel">
                  <option v-for="preset in IMAGE_SIZE_PRESETS" :key="preset.value" :value="preset.value">
                    {{ preset.value === 'auto' ? t('tools.image.sizeAuto') : preset.label }}
                  </option>
                </select>
              </div>
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-image-count">
                  {{ promptSegments.length > 1 ? t('tools.image.countPerPrompt') : t('tools.image.count') }}
                </label>
                <div class="relative">
                  <input
                    id="tools-image-count"
                    v-model.number="count"
                    type="number"
                    min="1"
                    :max="MAX_IMAGE_BATCH"
                    step="1"
                    class="input !py-1.5 !pr-7 text-sm [appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
                    :aria-label="promptSegments.length > 1 ? t('tools.image.countPerPrompt') : t('tools.image.count')"
                    @change="normalizeCount"
                  />
                  <div class="absolute right-1 top-1/2 flex -translate-y-1/2 flex-col">
                    <button
                      type="button"
                      class="flex h-3 w-4 items-center justify-center text-gray-900 transition-opacity hover:opacity-60 disabled:cursor-not-allowed disabled:opacity-20 dark:text-gray-100"
                      :disabled="effectiveCount >= MAX_IMAGE_BATCH"
                      :aria-label="t('tools.image.increaseCount')"
                      :title="t('tools.image.increaseCount')"
                      @mousedown.prevent
                      @click="incrementCount"
                    >
                      <svg viewBox="0 0 8 5" class="h-1.5 w-2.5 fill-current" aria-hidden="true">
                        <path d="M4 0L8 5H0z" />
                      </svg>
                    </button>
                    <button
                      type="button"
                      class="flex h-3 w-4 items-center justify-center text-gray-900 transition-opacity hover:opacity-60 disabled:cursor-not-allowed disabled:opacity-20 dark:text-gray-100"
                      :disabled="effectiveCount <= 1"
                      :aria-label="t('tools.image.decreaseCount')"
                      :title="t('tools.image.decreaseCount')"
                      @mousedown.prevent
                      @click="decrementCount"
                    >
                      <svg viewBox="0 0 8 5" class="h-1.5 w-2.5 fill-current" aria-hidden="true">
                        <path d="M0 0h8L4 5z" />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="family === 'openai'" class="grid grid-cols-2 gap-2">
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-image-output">
                  {{ t('tools.image.outputFormat') }}
                </label>
                <select id="tools-image-output" v-model="outputFormat" class="input !py-1.5 text-sm">
                  <option value="">{{ t('tools.image.qualityAuto') }}</option>
                  <option value="png">png</option>
                  <option value="jpeg">jpeg</option>
                  <option value="webp">webp</option>
                </select>
              </div>
              <div>
                <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-dark-300" for="tools-image-quality">
                  {{ t('tools.image.quality') }}
                </label>
                <select id="tools-image-quality" v-model="quality" class="input !py-1.5 text-sm">
                  <option value="">{{ t('tools.image.qualityAuto') }}</option>
                  <option value="low">{{ t('tools.image.qualityLow') }}</option>
                  <option value="medium">{{ t('tools.image.qualityMedium') }}</option>
                  <option value="high">{{ t('tools.image.qualityHigh') }}</option>
                  <option value="xhigh">{{ t('tools.image.qualityXHigh') }}</option>
                </select>
              </div>
            </div>

            <p v-if="family !== 'openai'" class="rounded-lg bg-gray-50 px-2.5 py-2 text-[11px] leading-4 text-gray-500 dark:bg-dark-700/30 dark:text-dark-400">
              {{ t('tools.image.basicParamsHint') }}
            </p>

            <div v-if="family === 'openai'" class="rounded-lg border border-gray-200 dark:border-dark-700">
              <button
                type="button"
                class="flex w-full items-center justify-between px-2.5 py-2 text-xs font-medium text-gray-600 dark:text-dark-300"
                @click="showAdvanced = !showAdvanced"
              >
                <span>{{ t('tools.image.advanced') }}</span>
                <Icon name="chevronDown" size="sm" :class="['transition-transform', showAdvanced && 'rotate-180']" />
              </button>
              <div v-if="showAdvanced" class="space-y-2.5 border-t border-gray-200 px-2.5 py-2.5 dark:border-dark-700">
                <div class="grid grid-cols-2 gap-2">
                  <div>
                    <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-background">
                      {{ t('tools.image.transparentBackground') }}
                    </label>
                    <select id="tools-image-background" v-model="background" class="input !py-1.5 text-sm">
                      <option value="">{{ t('tools.image.backgroundNo') }}</option>
                      <option value="transparent">{{ t('tools.image.backgroundYes') }}</option>
                    </select>
                  </div>
                  <div>
                    <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-style">
                      {{ t('tools.image.style') }}
                    </label>
                    <select id="tools-image-style" v-model="style" class="input !py-1.5 text-sm">
                      <option value="">{{ t('tools.image.qualityAuto') }}</option>
                      <option value="natural">{{ t('tools.image.styleNatural') }}</option>
                      <option value="vivid">{{ t('tools.image.styleVivid') }}</option>
                    </select>
                  </div>
                </div>
                <div class="grid grid-cols-2 gap-2">
                  <div>
                    <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-format">
                      {{ t('tools.image.responseFormat') }}
                    </label>
                    <select id="tools-image-format" v-model="responseFormat" class="input !py-1.5 text-sm">
                      <option value="">{{ t('tools.image.qualityAuto') }}</option>
                      <option value="url">{{ t('tools.image.formatUrl') }}</option>
                      <option value="b64_json">{{ t('tools.image.formatB64') }}</option>
                    </select>
                  </div>
                  <div>
                    <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-moderation">
                      {{ t('tools.image.moderation') }}
                    </label>
                    <select id="tools-image-moderation" v-model="moderation" class="input !py-1.5 text-sm">
                      <option value="">{{ t('tools.image.moderationAuto') }}</option>
                      <option value="low">{{ t('tools.image.moderationLow') }}</option>
                    </select>
                  </div>
                </div>
                <div v-if="canUseOutputCompression">
                  <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-compression">
                    {{ t('tools.image.outputCompression') }}
                    <span class="font-normal text-gray-400">0–100</span>
                  </label>
                  <input
                    id="tools-image-compression"
                    v-model="outputCompression"
                    type="number"
                    min="0"
                    max="100"
                    step="1"
                    class="input !py-1.5 text-sm"
                    placeholder="100"
                    @change="normalizeOutputCompression"
                  />
                  <p class="mt-1 text-[11px] leading-4 text-gray-400 dark:text-dark-500">
                    {{ t('tools.image.outputCompressionHint') }}
                  </p>
                </div>
                <div v-if="mode === 'edit'">
                  <label class="mb-1 block text-xs text-gray-500 dark:text-dark-400" for="tools-image-fidelity">
                    {{ t('tools.image.inputFidelity') }}
                  </label>
                  <select id="tools-image-fidelity" v-model="inputFidelity" class="input !py-1.5 text-sm">
                    <option value="">{{ t('tools.image.qualityAuto') }}</option>
                    <option value="high">{{ t('tools.image.fidelityHigh') }}</option>
                    <option value="low">{{ t('tools.image.fidelityLow') }}</option>
                  </select>
                </div>
              </div>
            </div>

            <button
              type="button"
              class="btn btn-primary w-full"
              :disabled="!canGenerate"
              @click="generate"
            >
              <Icon name="sparkles" size="sm" />
              {{ generateLabel }}
            </button>

            <p v-if="error" class="rounded-lg bg-red-50 px-2.5 py-2 text-xs text-red-600 dark:bg-red-900/20 dark:text-red-400">
              {{ error }}
            </p>
          </fieldset>
        </ToolsPanel>

        <!-- 第二栏：生成结果 -->
        <ToolsPanel
          class="xl:col-start-2"
          :title="t('tools.image.results')"
          icon="grid"
          :badge="activeRecord ? activeRecord.images.length : null"
          body-class="!overflow-hidden !p-0"
        >
          <template #leading>
            <button
              type="button"
              class="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-md text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-300 dark:hover:bg-dark-700 dark:hover:text-dark-200"
              :title="paramsCollapsed ? t('tools.image.expandParams') : t('tools.image.collapseParams')"
              :aria-label="paramsCollapsed ? t('tools.image.expandParams') : t('tools.image.collapseParams')"
              @click="paramsCollapsed = !paramsCollapsed"
            >
              <Icon :name="paramsCollapsed ? 'panelLeftOpen' : 'panelLeftClose'" size="sm" :stroke-width="2" />
            </button>
          </template>
          <template #trailing>
            <button
              type="button"
              class="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-md text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-300 dark:hover:bg-dark-700 dark:hover:text-dark-200"
              :title="historyCollapsed ? t('tools.image.expandHistory') : t('tools.image.collapseHistory')"
              :aria-label="historyCollapsed ? t('tools.image.expandHistory') : t('tools.image.collapseHistory')"
              @click="historyCollapsed = !historyCollapsed"
            >
              <Icon :name="historyCollapsed ? 'panelRightOpen' : 'panelRightClose'" size="sm" :stroke-width="2" />
            </button>
          </template>
          <template #actions>
            <button
              v-if="activeRecord"
              type="button"
              class="btn btn-secondary !px-2 !py-1 text-xs"
              @click="activeRecord = null"
            >
              {{ t('tools.common.clear') }}
            </button>
          </template>

          <div class="flex h-full min-h-0 flex-col">
            <div class="flex min-h-0 flex-1 p-3">
              <div class="flex min-h-0 flex-1 overflow-hidden rounded-xl border border-gray-200 bg-gray-50/40 dark:border-dark-700 dark:bg-dark-900/30">
                <div class="flex min-h-0 flex-1 overflow-auto p-3">
                  <div :class="activeRecord && activeRecord.images.length > 1 ? 'h-full w-full' : 'm-auto w-full'">
          <div v-if="generating" class="mx-auto flex min-h-full w-full flex-col items-center justify-center px-8 py-12 text-center">
            <div class="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-300">
              <Icon name="sparkles" size="lg" class="animate-pulse" />
            </div>
            <p class="mt-4 text-base font-medium text-gray-800 dark:text-dark-100">
              {{ t('tools.image.progressTitle') }}
            </p>
            <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">
              {{ t('tools.image.progressDetail', { done: progress.done, total: progress.total, count: totalImages }) }}
            </p>
            <div class="mt-6 h-2 w-full max-w-xl overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
              <div
                class="h-full rounded-full bg-primary-500 transition-all duration-500 dark:bg-primary-400"
                :class="progress.done === 0 ? 'animate-pulse' : ''"
                :style="{ width: progressPercent + '%' }"
              ></div>
            </div>
            <p class="mt-2.5 text-[11px] tabular-nums text-gray-400 dark:text-dark-500">
              {{ t('tools.image.progressElapsed', { seconds: elapsedSeconds }) }}
            </p>
          </div>

          <div
            v-else-if="activeRecord && activeRecord.images.length"
            :class="activeRecord.images.length === 1
              ? 'mx-auto w-full max-w-4xl'
              : 'grid h-full w-full auto-rows-fr grid-cols-2 gap-3 lg:grid-cols-3'"
          >
            <figure
              v-for="(image, index) in activeRecord.images"
              :key="index"
              class="group overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-dark-700 dark:bg-dark-900"
              :class="activeRecord.images.length === 1 ? 'w-full' : 'flex h-full min-h-0 flex-col'"
            >
              <img
                :src="image.src"
                :alt="activeRecord.prompt"
                loading="lazy"
                :class="activeRecord.images.length === 1
                  ? 'h-auto w-full cursor-zoom-in object-contain'
                  : 'min-h-0 w-full flex-1 cursor-zoom-in object-contain'"
                @load="onImageLoaded(image, $event)"
                @dblclick.stop="openImagePreview(image, activeRecord?.model)"
              />
              <div class="flex shrink-0 items-center justify-end gap-1 border-t border-gray-200 bg-white/95 px-2 py-1.5 dark:border-dark-700 dark:bg-dark-800/95">
                <button
                  type="button"
                  class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-[11px] text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-700 dark:hover:text-dark-200"
                  :title="t('tools.common.download')"
                  @click="downloadImage(image, activeRecord.model, index)"
                >
                  <Icon name="download" size="xs" />
                  {{ t('tools.common.download') }}
                </button>
                <button
                  type="button"
                  class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-[11px] text-primary-600 transition-colors hover:bg-primary-50 hover:text-primary-700 dark:text-primary-400 dark:hover:bg-primary-900/20 dark:hover:text-primary-300"
                  :title="t('tools.image.loadToEdit')"
                  @click="loadImageForEdit(image)"
                >
                  <Icon name="edit" size="xs" />
                  {{ t('tools.image.loadToEdit') }}
                </button>
              </div>
            </figure>
          </div>

          <div v-else class="flex h-full min-h-[240px] flex-col items-center justify-center gap-2 text-gray-400 dark:text-dark-500">
            <Icon name="sparkles" size="xl" />
            <p class="text-sm">{{ t('tools.image.empty') }}</p>
          </div>

                  </div>
                </div>
              </div>
            </div>

            <div v-if="activeRecord && !generating">
              <div class="h-px w-full bg-gray-200 dark:bg-dark-700"></div>
              <div class="px-3 py-2 sm:px-4">
                <div
                  class="flex h-[50px] items-center justify-center gap-3 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 px-4 text-xs text-gray-500 shadow-sm dark:border-dark-700 dark:bg-dark-700/40 dark:text-dark-400"
                >
                <span class="shrink-0 text-gray-400 dark:text-dark-500">{{ t('tools.image.summaryTitle') }}</span>
                <span class="shrink-0">·</span>
                <span class="shrink-0">{{ t('tools.image.summaryCount', { count: activeRecord.images.length }) }}</span>
                <span class="shrink-0 text-gray-300 dark:text-dark-600">|</span>
                <span class="truncate" :title="summaryResolution">{{ summaryResolution }}</span>
                <span class="shrink-0 text-gray-300 dark:text-dark-600">|</span>
                <span class="shrink-0" :title="summaryCost">{{ t('tools.image.summaryCost', { cost: summaryCost }) }}</span>
                <span class="shrink-0 text-gray-300 dark:text-dark-600">|</span>
                <span class="shrink-0">{{ t('tools.image.summaryDuration', { duration: summaryDuration }) }}</span>
                </div>
              </div>
            </div>
          </div>
        </ToolsPanel>

        <!-- 第三栏：历史记录（暂为会话内） -->
        <ToolsPanel v-if="!historyCollapsed" class="xl:col-start-3" :title="t('tools.common.history')" icon="clock" :badge="history.length">
          <p class="mb-2 text-xs text-gray-400 dark:text-dark-500">{{ t('tools.image.historyHint', { days: historyRetentionDays, count: historyMaxItems, size: Math.round(historyMaxBytes / 1024 / 1024) }) }}</p>
          <p v-if="history.length === 0" class="py-6 text-center text-xs text-gray-400 dark:text-dark-500">
            {{ t('tools.common.emptyHistory') }}
          </p>
          <div v-else class="space-y-2">
            <div
              v-for="(record, index) in history"
              :key="index"
              role="button"
              tabindex="0"
              class="w-full cursor-pointer rounded-lg border border-gray-200 p-2 text-left transition-colors hover:border-primary-300 hover:bg-primary-50/40 dark:border-dark-700 dark:hover:border-primary-700 dark:hover:bg-primary-900/10"
              :class="record === activeRecord ? 'border-primary-400 bg-primary-50/60 dark:border-primary-600 dark:bg-primary-900/20' : ''"
              @click="activeRecord = record"
              @keydown.enter.prevent="activeRecord = record"
              @keydown.space.prevent="activeRecord = record"
            >
              <div class="flex items-center justify-between gap-2 text-[11px] text-gray-400 dark:text-dark-500">
                <span>{{ formatTime(record.createdAt) }}</span>
                <span class="truncate">{{ record.model }}</span>
              </div>
              <div class="mt-1 flex items-center gap-1.5 text-[10px] text-gray-400 dark:text-dark-500">
                <span class="rounded bg-gray-100 px-1 py-0.5 dark:bg-dark-700">
                  {{ record.mode === 'edit' ? t('tools.image.modeEdit') : t('tools.image.modeGenerate') }}
                </span>
                <span v-if="record.params.size !== 'auto'">{{ record.params.size }}</span>
              </div>
              <p class="mt-1 truncate text-xs text-gray-700 dark:text-dark-200">{{ record.prompt }}</p>
              <div v-if="record !== activeRecord" class="mt-1.5 flex items-start gap-1.5">
                <div
                  v-for="(image, imageIndex) in record.images.slice(0, 4)"
                  :key="imageIndex"
                  class="w-12 shrink-0"
                >
                  <div class="group/image relative">
                    <img
                      :src="image.src"
                      class="h-9 w-9 cursor-zoom-in rounded object-cover"
                      alt=""
                      :title="imageInfoLabel(image, record.params.size)"
                      @load="onImageLoaded(image, $event)"
                      @dblclick.stop="openImagePreview(image, record.model)"
                    />
                    <button
                      type="button"
                      class="absolute -right-1 -top-1 inline-flex h-4 w-4 items-center justify-center rounded-full bg-black/70 text-white opacity-100 transition-colors hover:bg-red-600 sm:opacity-0 sm:group-hover/image:opacity-100"
                      :title="t('tools.image.deleteImage')"
                      :aria-label="t('tools.image.deleteImage')"
                      @click.stop="deleteHistoryImage(record, imageIndex)"
                    >
                      <Icon name="x" size="xs" />
                    </button>
                  </div>
                  <p class="mt-0.5 truncate text-[8px] leading-3 text-gray-400 dark:text-dark-500" :title="imageDimensionsLabel(image, record.params.size)">
                    {{ imageDimensionsLabel(image, record.params.size) }}
                  </p>
                </div>
                <span class="ml-auto pt-0.5 text-[11px] text-gray-400 dark:text-dark-500">
                  {{ t('tools.image.imageCount', { count: record.images.length }) }}
                </span>
              </div>
              <div v-else class="mt-2 space-y-1.5 border-t border-gray-200/80 pt-2 dark:border-dark-600/80">
                <div
                  v-for="(image, imageIndex) in record.images"
                  :key="imageIndex"
                  class="flex items-center gap-2 rounded-md bg-white/60 p-1.5 dark:bg-dark-800/50"
                >
                  <img
                    :src="image.src"
                    class="h-10 w-10 shrink-0 cursor-zoom-in rounded object-cover"
                    alt=""
                    @load="onImageLoaded(image, $event)"
                    @dblclick.stop="openImagePreview(image, record.model)"
                  />
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-[11px] font-medium text-gray-700 dark:text-dark-200">
                      {{ imageDimensionsLabel(image, record.params.size) }}
                    </p>
                    <p class="mt-0.5 text-[10px] text-gray-400 dark:text-dark-500">
                      {{ imageFormatLabel(image) }} · {{ imageBytesLabel(image) }}
                    </p>
                  </div>
                  <button
                    type="button"
                    class="inline-flex h-7 w-7 shrink-0 items-center justify-center rounded-md text-gray-400 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400"
                    :title="t('tools.image.deleteImage')"
                    :aria-label="t('tools.image.deleteImage')"
                    @click.stop="deleteHistoryImage(record, imageIndex)"
                  >
                    <Icon name="trash" size="xs" />
                  </button>
                </div>
                <div class="flex justify-end pt-0.5">
                  <span
                    class="inline-flex items-center gap-1 text-[11px] text-primary-600 dark:text-primary-400"
                    @click.stop="reuse(record)"
                  >
                    <Icon name="refresh" size="xs" />
                    {{ t('tools.common.reuse') }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </ToolsPanel>
      </div>
    </div>
  </AppLayout>

  <Teleport to="body">
    <div
      v-if="previewImage"
      class="fixed inset-0 z-[100] flex flex-col bg-gray-800/70 backdrop-blur-sm"
      @click.self="closeImagePreview"
    >
      <div class="flex flex-wrap items-center justify-between gap-2 px-3 py-2.5 text-white sm:px-4">
        <div class="min-w-0">
          <p class="truncate text-sm font-medium">{{ imageDimensionsLabel(previewImage) }}</p>
          <p class="mt-0.5 text-[11px] text-white/60">
            {{ imageFormatLabel(previewImage) }} · {{ imageBytesLabel(previewImage) }}
          </p>
        </div>
        <div class="flex flex-wrap items-center gap-1.5">
          <button
            type="button"
            class="rounded-md px-2.5 py-1.5 text-xs text-white/80 transition-colors hover:bg-white/10 hover:text-white"
            @click="togglePreviewActualSize"
          >
            {{ previewActualSize ? t('tools.image.fitScreen') : t('tools.image.actualSize') }}
          </button>
          <button
            type="button"
            class="inline-flex items-center gap-1 rounded-md px-2.5 py-1.5 text-xs text-white/80 transition-colors hover:bg-white/10 hover:text-white"
            @click="downloadImage(previewImage, previewModel, 0)"
          >
            <Icon name="download" size="xs" />
            {{ t('tools.common.download') }}
          </button>
          <button
            type="button"
            class="inline-flex items-center gap-1 rounded-md px-2.5 py-1.5 text-xs text-white/80 transition-colors hover:bg-white/10 hover:text-white"
            @click="loadImageForEdit(previewImage)"
          >
            <Icon name="edit" size="xs" />
            {{ t('tools.image.loadToEdit') }}
          </button>
          <button
            type="button"
            class="inline-flex h-8 w-8 items-center justify-center rounded-md text-white/70 transition-colors hover:bg-white/10 hover:text-white"
            :aria-label="t('tools.image.closePreview')"
            :title="t('tools.image.closePreview')"
            @click="closeImagePreview"
          >
            <Icon name="x" size="sm" />
          </button>
        </div>
      </div>
      <div
        class="min-h-0 flex-1 overflow-auto p-3 sm:p-5"
        :class="previewActualSize ? 'block' : 'flex items-center justify-center'"
        @click.self="closeImagePreview"
      >
        <img
          :src="previewImage.src"
          alt=""
          class="mx-auto block cursor-zoom-in"
          :class="previewActualSize ? 'max-h-none max-w-none' : 'max-h-full max-w-full object-contain'"
          @load="onPreviewLoaded"
          @dblclick="togglePreviewActualSize"
        />
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ToolsConnectionBar from '@/components/tools/ToolsConnectionBar.vue'
import ToolsPanel from '@/components/tools/ToolsPanel.vue'
import { GatewayError, useToolsSettings } from '@/composables/useToolsSettings'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import {
  toolImageHistoryAPI,
  type ToolImageHistoryItem,
} from '@/api/toolImageHistory'
import {
  IMAGE_SIZE_PRESETS,
  MAX_IMAGE_BATCH,
  MAX_IMAGE_UPLOAD_BYTES,
  buildImageGenerationPayload,
  detectImageModelFamily,
  imageExtensionFromMime,
  imageMimeFromFormat,
  imageSizePreset,
  splitImagePrompts,
  type ImageMode,
} from '@/utils/imageGeneration'

interface GeneratedImage {
  /** data: URL（b64_json）或上游 URL */
  src: string
  /** 是否为 data URL，决定下载方式 */
  isDataUrl: boolean
  mime: string
  width?: number
  height?: number
  bytes?: number
}

interface HistoryParams {
  mode: ImageMode
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

interface GenerationRecord {
  id: string
  prompt: string
  model: string
  mode: ImageMode
  images: GeneratedImage[]
  params: HistoryParams
  durationMs?: number
  costEstimate?: number
  createdAt: Date
}

const { t } = useI18n()
const { hasApiKey, imageModel, gatewayFetch, selectedGroupID, groupCatalog } = useToolsSettings()
const authStore = useAuthStore()
const appStore = useAppStore()

const prompt = ref('')
const mode = ref<ImageMode>('generate')
const size = ref('auto')
const count = ref(1)
const quality = ref('')
const style = ref('')
const responseFormat = ref('')
const background = ref('')
const outputFormat = ref('png')
const moderation = ref('')
const outputCompression = ref('100')
const partialImages = ref('')
const inputFidelity = ref('')
const showAdvanced = ref(false)

const generating = ref(false)
const error = ref('')
const progress = ref({ done: 0, total: 0 })
const elapsedSeconds = ref(0)
let elapsedTimer: ReturnType<typeof setInterval> | null = null

const progressPercent = computed(() => {
  if (!progress.value.total) return 8
  return Math.min(100, Math.max(8, Math.round((progress.value.done / progress.value.total) * 100)))
})

function startGenerationTimer(): void {
  stopGenerationTimer()
  elapsedSeconds.value = 0
  elapsedTimer = setInterval(() => {
    elapsedSeconds.value += 1
  }, 1000)
}

function stopGenerationTimer(): void {
  if (elapsedTimer !== null) {
    clearInterval(elapsedTimer)
    elapsedTimer = null
  }
}

const sourceInput = ref<HTMLInputElement | null>(null)
const maskInput = ref<HTMLInputElement | null>(null)
const sourceFile = ref<File | null>(null)
const maskFile = ref<File | null>(null)
const sourcePreview = ref('')
const maskPreview = ref('')

// 左右两栏折叠状态（记忆到本地，刷新后保持）
const PARAMS_COLLAPSED_KEY = 'tools_image_params_collapsed'
const HISTORY_COLLAPSED_KEY = 'tools_image_history_collapsed'
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
const paramsCollapsed = ref(readCollapsedFlag(PARAMS_COLLAPSED_KEY))
watch(paramsCollapsed, (value) => saveCollapsedFlag(PARAMS_COLLAPSED_KEY, value))
const historyCollapsed = ref(readCollapsedFlag(HISTORY_COLLAPSED_KEY))
watch(historyCollapsed, (value) => saveCollapsedFlag(HISTORY_COLLAPSED_KEY, value))

const history = ref<GenerationRecord[]>([])
const activeRecord = ref<GenerationRecord | null>(null)
const previewImage = ref<GeneratedImage | null>(null)
const previewActualSize = ref(false)
const previewModel = ref("")
const historyMaxItems = ref(20)
const historyRetentionDays = ref(7)
const historyMaxBytes = ref(48 * 1024 * 1024)

const HISTORY_STORAGE_KEY = 'tools_image_history_v1'
const HISTORY_OWNER_KEY = 'tools_image_history_owner'

function currentOwnerId(): string {
  const id = authStore.user?.id
  return id === undefined || id === null ? '' : String(id)
}

function newHistoryId(): string {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  return `img_${Date.now()}_${Math.random().toString(36).slice(2, 10)}`
}

function toApiItem(record: GenerationRecord): ToolImageHistoryItem {
  return {
    id: record.id,
    prompt: record.prompt,
    model: record.model,
    mode: record.mode,
    images: record.images.map((image) => ({
      src: image.src,
      isDataUrl: image.isDataUrl,
      mime: image.mime,
      width: image.width,
      height: image.height,
      bytes: image.bytes,
    })),
    params: record.params,
    duration_ms: record.durationMs,
    cost_estimate: record.costEstimate,
    createdAt: record.createdAt.getTime(),
  }
}

function fromApiItem(item: ToolImageHistoryItem): GenerationRecord {
  return {
    id: item.id,
    prompt: item.prompt,
    model: item.model,
    mode: item.mode === 'edit' ? 'edit' : 'generate',
    images: Array.isArray(item.images)
      ? item.images.map((image) => ({
          src: image.src,
          isDataUrl: image.isDataUrl,
          mime: image.mime || imageMimeFromFormat(item.params?.outputFormat),
          width: image.width,
          height: image.height,
          bytes: image.bytes,
        }))
      : [],
    params: {
      ...item.params,
      mode: item.params?.mode === 'edit' ? 'edit' : 'generate',
    },
    durationMs: item.duration_ms,
    costEstimate: item.cost_estimate,
    createdAt: new Date(item.createdAt),
  }
}

function loadLocalHistory(): GenerationRecord[] {
  const owner = currentOwnerId()
  try {
    const storedOwner = localStorage.getItem(HISTORY_OWNER_KEY) ?? ''
    if (storedOwner && owner && storedOwner !== owner) return []
    if (!storedOwner && owner) localStorage.setItem(HISTORY_OWNER_KEY, owner)
    const raw = localStorage.getItem(HISTORY_STORAGE_KEY)
    if (!raw) return []
    const parsed: unknown = JSON.parse(raw)
    if (!Array.isArray(parsed)) return []
    const records: GenerationRecord[] = []
    for (const value of parsed) {
      const item = value as Partial<ToolImageHistoryItem>
      if (typeof item?.id !== 'string' || !Array.isArray(item.images)) continue
      records.push(fromApiItem(item as ToolImageHistoryItem))
    }
    return records.slice(0, historyMaxItems.value)
  } catch {
    return []
  }
}

function persistLocalHistory(): void {
  const owner = currentOwnerId()
  try {
    if (owner) localStorage.setItem(HISTORY_OWNER_KEY, owner)
    localStorage.setItem(
      HISTORY_STORAGE_KEY,
      JSON.stringify(history.value.slice(0, historyMaxItems.value).map(toApiItem)),
    )
  } catch {
    // localStorage 容量不足时仅保留服务端历史，不影响当前页面。
  }
}

function mergeHistory(local: GenerationRecord[], remote: GenerationRecord[]): GenerationRecord[] {
  const merged = new Map<string, GenerationRecord>()
  for (const record of remote) merged.set(record.id, record)
  for (const record of local) {
    if (!merged.has(record.id)) merged.set(record.id, record)
  }
  return [...merged.values()]
    .sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime())
    .slice(0, historyMaxItems.value)
}

async function uploadMissingHistory(local: GenerationRecord[], remote: GenerationRecord[]): Promise<void> {
  const existing = new Set(remote.map((record) => record.id))
  for (const record of local) {
    if (existing.has(record.id)) continue
    try {
      await toolImageHistoryAPI.append(toApiItem(record))
    } catch (error) {
      console.warn('Failed to migrate image history to server', error)
      break
    }
  }
}

async function loadHistory(): Promise<void> {
  const local = loadLocalHistory()
  history.value = mergeHistory(history.value, local)
  if (!activeRecord.value) activeRecord.value = history.value[0] ?? null
  try {
    const result = await toolImageHistoryAPI.list()
    historyMaxItems.value = result.maxItems
    historyMaxBytes.value = result.maxBytes
    historyRetentionDays.value = result.retentionDays
    const remote = result.items.map(fromApiItem)
    // 请求期间可能刚好生成了一张图，合并当前内存而不是只合并请求前的 local。
    history.value = mergeHistory(history.value, remote)
    if (!activeRecord.value || !history.value.some((record) => record.id === activeRecord.value?.id)) {
      activeRecord.value = history.value[0] ?? null
    }
    persistLocalHistory()
    // 首次上线或离线期间产生的本地记录，补写一次服务端。
    await uploadMissingHistory(history.value, remote)
  } catch (error) {
    console.warn('Failed to load image history from server', error)
  }
}

const family = computed(() => detectImageModelFamily(imageModel.value))
const sizeLabel = computed(() => (size.value === 'auto' ? t('tools.image.sizeAuto') : imageSizePreset(size.value).label))
const canUseOutputCompression = computed(() => outputFormat.value === 'jpeg' || outputFormat.value === 'webp')
const promptSegments = computed(() => splitImagePrompts(prompt.value))
const effectiveCount = computed(() => {
  const value = Math.floor(Number(count.value))
  if (!Number.isFinite(value)) return 1
  return Math.min(MAX_IMAGE_BATCH, Math.max(1, value))
})
const totalImages = computed(() => promptSegments.value.length * effectiveCount.value)
const hasBatchOverflow = computed(() => totalImages.value > MAX_IMAGE_BATCH)
const activeModelPricing = computed(() => {
  const groupID = selectedGroupID.value
  if (groupID === null) return null
  const entry = groupCatalog.value[String(groupID)]
  const model = entry?.models.find((item) => item.name === imageModel.value.trim())
  return model?.pricing ?? null
})
const estimatedUnitCost = computed(() => {
  const pricing = activeModelPricing.value
  if (!pricing) return null
  if (pricing.billing_mode !== 'image' && pricing.billing_mode !== 'per_request') return null
  const price = pricing.per_request_price
  return typeof price === 'number' && price > 0 ? price : null
})
const summaryResolution = computed(() => {
  if (!activeRecord.value) return '—'
  const values = activeRecord.value.images.map((image) => imageDimensionsLabel(image, activeRecord.value?.params.size))
  const unique = [...new Set(values.filter((value) => value && value !== '—'))]
  if (unique.length === 0) return '—'
  return unique.length === 1 ? unique[0] : t('tools.image.mixedResolution')
})
const summaryCost = computed(() => {
  const record = activeRecord.value
  if (!record) return '—'
  return typeof record.costEstimate === 'number' ? `~$${record.costEstimate.toFixed(4)}` : '—'
})
const summaryDuration = computed(() => {
  const duration = activeRecord.value?.durationMs
  return typeof duration === 'number' && duration > 0 ? `${(duration / 1000).toFixed(1)}s` : '—'
})
const canGenerate = computed(
  () =>
    !generating.value &&
    promptSegments.value.length > 0 &&
    imageModel.value.trim() !== '' &&
    hasApiKey.value &&
    !hasBatchOverflow.value &&
    (mode.value === 'generate' || sourceFile.value !== null),
)
const generateLabel = computed(() => {
  if (!generating.value) return t('tools.image.generate')
  if (progress.value.total > 1) {
    return t('tools.image.generatingProgress', { done: progress.value.done, total: progress.value.total })
  }
  return t('tools.image.generating')
})

function normalizeCount(): void {
  const value = Math.floor(Number(count.value))
  count.value = Number.isFinite(value) ? Math.min(MAX_IMAGE_BATCH, Math.max(1, value)) : 1
}

function incrementCount(): void {
  if (effectiveCount.value >= MAX_IMAGE_BATCH) return
  count.value = effectiveCount.value + 1
}

function decrementCount(): void {
  if (effectiveCount.value <= 1) return
  count.value = effectiveCount.value - 1
}

function normalizeOutputCompression(): void {
  const value = Math.floor(Number(outputCompression.value))
  outputCompression.value = String(Number.isFinite(value) ? Math.min(100, Math.max(0, value)) : 100)
}

function localizedError(err: unknown): string {
  if (err instanceof GatewayError) {
    if (err.code === 'INSUFFICIENT_BALANCE') return t('tools.errors.insufficientBalance')
    if (err.status === 401) return t('tools.errors.unauthorized')
    if (err.status === 403) return t('tools.errors.forbidden')
    if (err.status === 429) return t('tools.errors.rateLimited')
    return err.message
  }
  return err instanceof Error ? err.message : t('tools.errors.unknown')
}

function formatTime(date: Date): string {
  return date.toLocaleTimeString()
}

function dimensionsFromSize(value: unknown): { width: number; height: number } | null {
  if (typeof value !== 'string') return null
  const match = value.trim().match(/^(\d+)\s*[x×]\s*(\d+)$/i)
  if (!match) return null
  const width = Number(match[1])
  const height = Number(match[2])
  return Number.isFinite(width) && Number.isFinite(height) && width > 0 && height > 0
    ? { width, height }
    : null
}

function base64Bytes(value: string): number | undefined {
  if (!value) return undefined
  const padding = value.endsWith('==') ? 2 : value.endsWith('=') ? 1 : 0
  return Math.max(0, Math.floor((value.length * 3) / 4) - padding)
}

/** 解析 /v1/images/generations 或 /v1/images/edits 的图片响应。 */
function parseImages(body: unknown, fallbackSize = size.value): GeneratedImage[] {
  const data = (body as { data?: unknown })?.data
  if (!Array.isArray(data)) return []
  const images: GeneratedImage[] = []
  for (const item of data) {
    const b64 = (item as { b64_json?: unknown })?.b64_json
    const url = (item as { url?: unknown })?.url
    const format = (item as { output_format?: unknown })?.output_format ?? outputFormat.value
    const mime = imageMimeFromFormat(format)
    const itemWidth = Number((item as { width?: unknown })?.width)
    const itemHeight = Number((item as { height?: unknown })?.height)
    const fromSize = dimensionsFromSize((item as { size?: unknown })?.size) ?? dimensionsFromSize(fallbackSize)
    const width = Number.isFinite(itemWidth) && itemWidth > 0 ? itemWidth : fromSize?.width
    const height = Number.isFinite(itemHeight) && itemHeight > 0 ? itemHeight : fromSize?.height
    if (typeof b64 === 'string' && b64) {
      images.push({
        src: `data:${mime};base64,${b64}`,
        isDataUrl: true,
        mime,
        width,
        height,
        bytes: base64Bytes(b64),
      })
    } else if (typeof url === 'string' && url) {
      images.push({ src: url, isDataUrl: false, mime, width, height })
    }
  }
  return images
}

function currentParams(): HistoryParams {
  return {
    mode: mode.value,
    size: size.value,
    count: effectiveCount.value,
    quality: quality.value,
    style: style.value,
    responseFormat: responseFormat.value,
    background: background.value,
    outputFormat: outputFormat.value,
    moderation: moderation.value,
    outputCompression: outputCompression.value,
    partialImages: partialImages.value,
    inputFidelity: inputFidelity.value,
  }
}

function buildPayload(promptText: string): Record<string, unknown> {
  return buildImageGenerationPayload({
    family: family.value,
    model: imageModel.value,
    prompt: promptText,
    count: effectiveCount.value,
    size: size.value,
    quality: quality.value,
    style: style.value,
    responseFormat: responseFormat.value,
    background: background.value,
    outputFormat: outputFormat.value,
    moderation: moderation.value,
    outputCompression: outputCompression.value,
    partialImages: partialImages.value,
    inputFidelity: inputFidelity.value,
  })
}

function buildEditForm(payload: Record<string, unknown>): FormData {
  const form = new FormData()
  for (const [key, value] of Object.entries(payload)) {
    if (value === undefined || value === null || value === '') continue
    if (typeof value === 'number' || typeof value === 'string' || typeof value === 'boolean') {
      form.append(key, String(value))
    }
  }
  if (sourceFile.value) form.append('image', sourceFile.value)
  if (maskFile.value) form.append('mask', maskFile.value)
  return form
}

function saveRecord(images: GeneratedImage[], durationMs: number): void {
  const record: GenerationRecord = {
    id: newHistoryId(),
    prompt: promptSegments.value.join('; '),
    model: imageModel.value.trim(),
    mode: mode.value,
    images,
    params: currentParams(),
    durationMs,
    costEstimate:
      estimatedUnitCost.value !== null ? estimatedUnitCost.value * images.length : undefined,
    createdAt: new Date(),
  }
  history.value.unshift(record)
  history.value = history.value.slice(0, historyMaxItems.value)
  activeRecord.value = record
  persistLocalHistory()
  void toolImageHistoryAPI.append(toApiItem(record)).catch((err) => {
    console.warn('Failed to save image history to server', err)
    const reason = extractApiErrorMessage(err, '').toLowerCase()
    appStore.showError(
      reason.includes('too large')
        ? t('tools.image.imageHistoryTooLarge')
        : t('tools.image.imageHistorySaveFailed'),
    )
  })
}

async function generate(): Promise<void> {
  if (!canGenerate.value) return
  const segments = promptSegments.value
  if (totalImages.value > MAX_IMAGE_BATCH) {
    error.value = t('tools.image.batchTooLarge', { max: MAX_IMAGE_BATCH })
    return
  }
  if (mode.value === 'edit' && !sourceFile.value) {
    error.value = t('tools.image.sourceImageRequired')
    return
  }

  generating.value = true
  startGenerationTimer()
  const startedAt = performance.now()
  error.value = ''
  progress.value = { done: 0, total: segments.length }
  const images: GeneratedImage[] = []
  try {
    for (const [index, segment] of segments.entries()) {
      try {
        const payload = buildPayload(segment)
        let res: Response
        if (mode.value === 'edit') {
          res = await gatewayFetch('/v1/images/edits', {
            method: 'POST',
            body: buildEditForm(payload),
          })
        } else {
          res = await gatewayFetch('/v1/images/generations', {
            method: 'POST',
            body: JSON.stringify(payload),
          })
        }
        const body: unknown = await res.json()
        images.push(...parseImages(body, size.value))
        progress.value = { done: index + 1, total: segments.length }
      } catch (err) {
        // 多段提示词中某一段失败时，保留已经生成成功的图片。
        error.value = localizedError(err)
        break
      }
    }
    if (images.length > 0) {
      saveRecord(images, Math.max(0, Math.round(performance.now() - startedAt)))
    } else if (!error.value) {
      error.value = t('tools.image.noResults')
    }
  } finally {
    generating.value = false
    stopGenerationTimer()
  }
}

/** 把历史记录里的参数回填到表单，方便微调后重试。 */
function reuse(record: GenerationRecord): void {
  mode.value = record.mode === 'edit' ? 'edit' : 'generate'
  prompt.value = record.prompt
  imageModel.value = record.model
  size.value = record.params.size
  count.value = record.params.count
  quality.value = record.params.quality
  style.value = record.params.style
  responseFormat.value = record.params.responseFormat
  background.value = record.params.background
  outputFormat.value = record.params.outputFormat
  moderation.value = record.params.moderation
  outputCompression.value = record.params.outputCompression
  partialImages.value = record.params.partialImages
  inputFidelity.value = record.params.inputFidelity
  if (record.mode === 'edit') {
    // File 不能从历史恢复，必须重新选择原图。
    clearFile('source')
    clearFile('mask')
  }
  error.value = ''
  activeRecord.value = record
}

function restoreDefaults(): void {
  mode.value = 'generate'
  size.value = 'auto'
  count.value = 1
  quality.value = ''
  style.value = ''
  responseFormat.value = ''
  background.value = ''
  outputFormat.value = 'png'
  moderation.value = ''
  outputCompression.value = '100'
  partialImages.value = ''
  inputFidelity.value = ''
  showAdvanced.value = false
  clearFile('source')
  clearFile('mask')
  error.value = ''
}

function onFileChange(kind: 'source' | 'mask', event: Event): void {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  const allowedTypes = new Set(['image/png', 'image/jpeg', 'image/webp'])
  if (!allowedTypes.has(file.type)) {
    error.value = t('tools.image.invalidImage')
    input.value = ''
    return
  }
  if (file.size > MAX_IMAGE_UPLOAD_BYTES) {
    error.value = t('tools.image.imageTooLarge', { max: 20 })
    input.value = ''
    return
  }
  if (kind === 'source') {
    if (sourcePreview.value) URL.revokeObjectURL(sourcePreview.value)
    sourceFile.value = file
    sourcePreview.value = URL.createObjectURL(file)
  } else {
    if (maskPreview.value) URL.revokeObjectURL(maskPreview.value)
    maskFile.value = file
    maskPreview.value = URL.createObjectURL(file)
  }
  error.value = ''
  input.value = ''
}

function clearFile(kind: 'source' | 'mask'): void {
  if (kind === 'source') {
    if (sourcePreview.value) URL.revokeObjectURL(sourcePreview.value)
    sourceFile.value = null
    sourcePreview.value = ''
    if (sourceInput.value) sourceInput.value.value = ''
  } else {
    if (maskPreview.value) URL.revokeObjectURL(maskPreview.value)
    maskFile.value = null
    maskPreview.value = ''
    if (maskInput.value) maskInput.value.value = ''
  }
}

onMounted(() => {
  void loadHistory()
  window.addEventListener('keydown', onGlobalKeydown)
})

function onGlobalKeydown(event: KeyboardEvent): void {
  if (event.key === 'Escape' && previewImage.value) closeImagePreview()
}

watch(previewImage, (image) => {
  document.body.style.overflow = image ? 'hidden' : ''
})

// 同一浏览器切换账号时，清掉上一个账号的图片历史，避免跨账号展示。
watch(
  () => authStore.user?.id,
  async (nextOwner, previousOwner) => {
    if (nextOwner === previousOwner) return
    history.value = []
    activeRecord.value = null
    if (nextOwner !== undefined && nextOwner !== null) await loadHistory()
  },
)

onBeforeUnmount(() => {
  stopGenerationTimer()
  window.removeEventListener('keydown', onGlobalKeydown)
  document.body.style.overflow = ''
  if (sourcePreview.value) URL.revokeObjectURL(sourcePreview.value)
  if (maskPreview.value) URL.revokeObjectURL(maskPreview.value)
})

function onImageLoaded(image: GeneratedImage, event: Event): void {
  const element = event.target as HTMLImageElement
  if (!element.naturalWidth || !element.naturalHeight) return
  if (image.width === element.naturalWidth && image.height === element.naturalHeight) return
  image.width = element.naturalWidth
  image.height = element.naturalHeight
  if ((!image.bytes || image.bytes <= 0) && image.isDataUrl) {
    const comma = image.src.indexOf(',')
    if (comma >= 0) image.bytes = base64Bytes(image.src.slice(comma + 1))
  }
}

function imageDimensionsLabel(image: GeneratedImage, fallbackSize?: string): string {
  if (image.width && image.height) return `${image.width} × ${image.height}`
  const fallback = dimensionsFromSize(fallbackSize)
  if (fallback) return `${fallback.width} × ${fallback.height}`
  return '—'
}

function imageFormatLabel(image: GeneratedImage): string {
  const subtype = (image.mime || 'image/png').split('/')[1] || 'png'
  return subtype.toUpperCase()
}

function imageBytesLabel(image: GeneratedImage): string {
  let bytes = image.bytes
  if ((!bytes || bytes <= 0) && image.isDataUrl) {
    const comma = image.src.indexOf(',')
    if (comma >= 0) bytes = base64Bytes(image.src.slice(comma + 1))
  }
  if (!bytes || bytes <= 0) return '—'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}

function imageInfoLabel(image: GeneratedImage, fallbackSize?: string): string {
  return `${imageDimensionsLabel(image, fallbackSize)} · ${imageFormatLabel(image)} · ${imageBytesLabel(image)}`
}

async function deleteHistoryImage(record: GenerationRecord, index: number): Promise<void> {
  if (index < 0 || index >= record.images.length) return
  const nextImages = record.images.filter((_, imageIndex) => imageIndex !== index)
  const nextRecord: GenerationRecord = { ...record, images: nextImages }
  history.value = history.value
    .map((item) => (item.id === record.id ? nextRecord : item))
    .filter((item) => item.images.length > 0)
  if (activeRecord.value?.id === record.id) {
    activeRecord.value = nextImages.length > 0 ? nextRecord : history.value[0] ?? null
  } else if (!activeRecord.value) {
    activeRecord.value = history.value[0] ?? null
  }
  persistLocalHistory()
  try {
    if (nextImages.length > 0) {
      await toolImageHistoryAPI.update(toApiItem(nextRecord))
    } else {
      await toolImageHistoryAPI.delete(record.id)
    }
  } catch (err) {
    console.warn('Failed to delete image history', err)
    appStore.showError(t('tools.image.imageHistorySaveFailed'))
  }
}

function openImagePreview(image: GeneratedImage, model = ''): void {
  previewActualSize.value = false
  previewModel.value = model
  previewImage.value = image
}

function closeImagePreview(): void {
  previewImage.value = null
  previewActualSize.value = false
}

function togglePreviewActualSize(): void {
  previewActualSize.value = !previewActualSize.value
}

function onPreviewLoaded(event: Event): void {
  const element = event.target as HTMLImageElement
  const image = previewImage.value
  if (!image || !element.naturalWidth || !element.naturalHeight) return
  image.width = element.naturalWidth
  image.height = element.naturalHeight
}

async function loadImageForEdit(image: GeneratedImage): Promise<void> {
  try {
    const response = await fetch(image.src)
    if (!response.ok) throw new Error(`image fetch failed: ${response.status}`)
    const blob = await response.blob()
    if (blob.size > MAX_IMAGE_UPLOAD_BYTES) {
      appStore.showError(t('tools.image.imageTooLarge', { max: 20 }))
      return
    }
    const mime = image.mime || blob.type || 'image/png'
    const extension = imageExtensionFromMime(mime)
    const file = new File([blob], `edit-source-${Date.now()}.${extension}`, { type: mime })
    clearFile('source')
    clearFile('mask')
    sourceFile.value = file
    sourcePreview.value = URL.createObjectURL(blob)
    mode.value = 'edit'
    error.value = ''
    closeImagePreview()
  } catch (err) {
    console.warn('Failed to load image for editing', err)
    appStore.showError(t('tools.image.loadToEditFailed'))
  }
}

async function downloadImage(image: GeneratedImage, model: string, index: number): Promise<void> {
  const link = document.createElement('a')
  const safeModel = model.replace(/[^a-zA-Z0-9._-]/g, '-') || 'image'
  const filename = `${safeModel}-${Date.now()}-${index + 1}.${imageExtensionFromMime(image.mime)}`
  if (image.isDataUrl) {
    link.href = image.src
    link.download = filename
    link.click()
    return
  }
  try {
    const response = await fetch(image.src)
    if (!response.ok) throw new Error(`image download failed: ${response.status}`)
    const blob = await response.blob()
    const objectUrl = URL.createObjectURL(blob)
    link.href = objectUrl
    link.download = filename
    link.click()
    setTimeout(() => URL.revokeObjectURL(objectUrl), 1000)
  } catch {
    // 跨域或上游链接不可直接读取时，退化为新标签页打开。
    link.href = image.src
    link.target = '_blank'
    link.rel = 'noopener'
    link.click()
  }
}
</script>