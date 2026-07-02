export function formatCurrency(amount: number, locale = 'ar-SA'): string {
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency: 'SAR',
    maximumFractionDigits: 0,
  }).format(amount)
}

export function formatDate(date: string, locale = 'ar-SA'): string {
  return new Intl.DateTimeFormat(locale, {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }).format(new Date(date))
}

export function formatPhone(phone: string): string {
  const digits = phone.replace(/\D/g, '')
  if (digits.startsWith('966')) return `+${digits}`
  if (digits.startsWith('05')) return `+966${digits.slice(1)}`
  return `+966${digits}`
}

export function normalizeSaudiPhone(input: string): string {
  const digits = input.replace(/\D/g, '')
  if (digits.startsWith('966')) return `+${digits}`
  if (digits.startsWith('05')) return `+966${digits.slice(1)}`
  if (digits.startsWith('5') && digits.length === 9) return `+966${digits}`
  return `+966${digits}`
}

/** Convert Arabic-Indic / Persian digits to Western 0-9 */
export function normalizeDigits(input: string): string {
  const arabic = '٠١٢٣٤٥٦٧٨٩'
  const persian = '۰۱۲۳۴۵۶۷۸۹'
  return input
    .split('')
    .map((ch) => {
      const ai = arabic.indexOf(ch)
      if (ai >= 0) return String(ai)
      const pi = persian.indexOf(ch)
      if (pi >= 0) return String(pi)
      return ch
    })
    .join('')
}
