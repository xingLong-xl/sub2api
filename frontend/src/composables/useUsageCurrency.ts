import { computed } from 'vue'
import { getActivePinia } from 'pinia'
import { useAppStore } from '@/stores/app'
import { getLocale } from '@/i18n'

/**
 * 用量/额度金额显示币种。
 * 数据库金额统一为 CNY；USD 仅做展示换算，不影响任何计费、支付或订单金额。
 */
export function useUsageCurrency() {
  // 单测可能没有挂 Pinia；此时按默认 CNY 展示，不影响生产。
  const appStore = getActivePinia() ? useAppStore() : null

  const currency = computed<'CNY' | 'USD'>(() =>
    appStore?.cachedPublicSettings?.usage_currency === 'USD' ? 'USD' : 'CNY',
  )
  const usdCnyRate = computed(() => {
    const value = Number(appStore?.cachedPublicSettings?.usd_cny_rate)
    return Number.isFinite(value) && value > 0 ? value : 7.2
  })
  const currencySymbol = computed(() => (currency.value === 'USD' ? '$' : '¥'))

  function convert(amountCNY: number | null | undefined): number {
    const value = Number(amountCNY ?? 0)
    if (!Number.isFinite(value)) return 0
    return currency.value === 'USD' ? value / usdCnyRate.value : value
  }

  function formatUsageMoney(amountCNY: number | null | undefined, fractionDigits = 2): string {
    const value = convert(amountCNY)
    return new Intl.NumberFormat(getLocale(), {
      style: 'currency',
      currency: currency.value,
      currencyDisplay: 'narrowSymbol',
      minimumFractionDigits: fractionDigits,
      maximumFractionDigits: fractionDigits,
    }).format(value)
  }

  return { currency, usdCnyRate, currencySymbol, convert, formatUsageMoney }
}
