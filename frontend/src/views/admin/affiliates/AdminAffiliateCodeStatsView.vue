<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <div class="relative w-full md:w-80">
            <Icon name="search" size="md" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              v-model="filters.search"
              type="text"
              class="input pl-10"
              :placeholder="t('admin.affiliates.codeStats.searchPlaceholder')"
              @input="debounceLoad"
            />
          </div>
          <DateRangePicker
            :start-date="filters.start_at"
            :end-date="filters.end_at"
            @update:start-date="(value) => (filters.start_at = value)"
            @update:end-date="(value) => (filters.end_at = value)"
            @change="reloadFromFirstPage"
          />
          <button
            v-if="filters.start_at || filters.end_at"
            type="button"
            class="btn btn-secondary px-3"
            :title="t('common.clear')"
            @click="clearDates"
          >
            {{ t('common.clear') }}
          </button>
          <button class="btn btn-secondary px-2 md:px-3" :disabled="loading" :title="t('common.refresh')" @click="loadStats">
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
        </div>
      </template>

      <template #table>
        <DataTable :columns="columns" :data="records" :loading="loading">
          <template #cell-aff_code="{ row }">
            <span class="font-mono text-sm font-medium text-gray-900 dark:text-white">{{ row.aff_code }}</span>
          </template>
          <template #cell-inviter="{ row }">
            <div class="min-w-0">
              <div class="truncate text-sm font-medium text-gray-900 dark:text-white">
                {{ row.inviter_email || '-' }}
              </div>
              <div v-if="row.inviter_username" class="truncate text-xs text-gray-500 dark:text-dark-400">
                {{ row.inviter_username }}
              </div>
            </div>
          </template>
          <template #cell-invited_count="{ row }">
            <span class="text-sm text-gray-900 dark:text-white">{{ row.invited_count }}</span>
          </template>
          <template #cell-recharge_user_count="{ row }">
            <span class="text-sm font-semibold text-primary-600 dark:text-primary-400">
              {{ row.recharge_user_count }}
            </span>
          </template>
          <template #cell-recharge_order_count="{ row }">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ row.recharge_order_count }}</span>
          </template>
          <template #cell-recharge_amount="{ row }">
            <span class="text-sm font-medium text-green-600 dark:text-green-400">
              {{ formatCurrency(row.recharge_amount) }}
            </span>
          </template>
          <template #cell-rebate_amount="{ row }">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ formatCurrency(row.rebate_amount) }}</span>
          </template>
          <template #cell-last_recharge_at="{ row }">
            <span v-if="row.last_recharge_at" class="text-sm text-gray-700 dark:text-gray-300">
              {{ formatDateTime(row.last_recharge_at) }}
            </span>
            <span v-else class="text-sm text-gray-400 dark:text-dark-500">-</span>
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import Icon from '@/components/icons/Icon.vue'
import type { Column } from '@/components/common/types'
import { affiliatesAPI, type AffiliateCodeStats } from '@/api/admin/affiliates'
import { formatCurrency, formatDateTime } from '@/utils/format'

const { t } = useI18n()
const loading = ref(false)
const records = ref<AffiliateCodeStats[]>([])
const filters = reactive({ search: '', start_at: '', end_at: '' })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const columns = computed<Column[]>(() => [
  { key: 'aff_code', label: t('admin.affiliates.codeStats.affCode') },
  { key: 'inviter', label: t('admin.affiliates.records.inviter') },
  { key: 'invited_count', label: t('admin.affiliates.codeStats.invitedCount') },
  { key: 'recharge_user_count', label: t('admin.affiliates.codeStats.rechargeUserCount') },
  { key: 'recharge_order_count', label: t('admin.affiliates.codeStats.rechargeOrderCount') },
  { key: 'recharge_amount', label: t('admin.affiliates.codeStats.rechargeAmount') },
  { key: 'rebate_amount', label: t('admin.affiliates.records.rebateAmount') },
  { key: 'last_recharge_at', label: t('admin.affiliates.codeStats.lastRechargeAt') },
])

function timezone(): string {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone
  } catch {
    return 'UTC'
  }
}

async function loadStats() {
  loading.value = true
  try {
    const res = await affiliatesAPI.listCodeStats({
      page: pagination.page,
      page_size: pagination.page_size,
      search: filters.search.trim() || undefined,
      start_at: filters.start_at || undefined,
      end_at: filters.end_at || undefined,
      timezone: timezone(),
    })
    records.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to load affiliate code stats', error)
    records.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

function debounceLoad() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(reloadFromFirstPage, 300)
}

function clearDates() {
  filters.start_at = ''
  filters.end_at = ''
  reloadFromFirstPage()
}

function reloadFromFirstPage() {
  pagination.page = 1
  loadStats()
}

function handlePageChange(page: number) {
  pagination.page = page
  loadStats()
}

function handlePageSizeChange(pageSize: number) {
  pagination.page_size = pageSize
  pagination.page = 1
  loadStats()
}

onMounted(loadStats)
onUnmounted(() => {
  if (debounceTimer) clearTimeout(debounceTimer)
})
</script>
