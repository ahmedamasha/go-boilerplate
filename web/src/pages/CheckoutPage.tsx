import { useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'
import { api } from '../api/client'
import type { Event, Gift } from '../api/types'
import { TopAppBar } from '../components/layout/TopAppBar'
import { Button } from '../components/ui/Button'
import { Icon } from '../components/ui/Icon'
import { useAuth } from '../context/AuthContext'
import { useLocale } from '../context/LocaleContext'
import { formatCurrency } from '../lib/format'

export function CheckoutPage() {
  const { id } = useParams<{ id: string }>()
  const { t } = useLocale()
  const { token } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const { gift, event } = (location.state as { gift?: Gift; event?: Event }) || {}

  const [amount, setAmount] = useState(gift?.target_amount ? Math.min(500, gift.target_amount) : 500)
  const [hideAmount, setHideAmount] = useState(false)
  const [message, setMessage] = useState('')
  const [paymentMethod, setPaymentMethod] = useState<'mada' | 'apple_pay'>('mada')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [reference, setReference] = useState('')

  if (!gift || !event) {
    navigate(`/events/${id}`)
    return null
  }

  async function handlePay() {
    setLoading(true)
    try {
      const contribution = await api.contribute(
        {
          event_id: event!.id,
          gift_id: gift!.id !== 'cash' ? gift!.id : undefined,
          amount,
          hide_amount: hideAmount,
          message,
        },
        token,
      )
      const payment = await api.pay(contribution.id, paymentMethod)
      setReference(payment.reference_number)
      setSuccess(true)
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center px-container-margin text-center space-y-section-gap">
        <div className="w-20 h-20 bg-primary-fixed rounded-full flex items-center justify-center">
          <Icon name="check_circle" className="text-primary text-5xl filled" />
        </div>
        <h1 className="text-headline-xl text-primary">{t('تم الإهداء بنجاح!', 'Gift sent successfully!')}</h1>
        <p className="text-body-md text-on-surface-variant">
          {t(`رقم المرجع: ${reference}`, `Reference: ${reference}`)}
        </p>
        <div className="bg-surface-container-lowest rounded-lg p-6 shadow-card w-full max-w-sm">
          <h4 className="text-headline-lg font-bold mb-1">{t('بطاقة إهداء', 'Gift card')}</h4>
          <p className="text-body-sm text-on-surface-variant">{message || t('مع أطيب التمنيات', 'With best wishes')}</p>
        </div>
        <Button onClick={() => navigate(`/events/${id}`)}>{t('العودة للمناسبة', 'Back to event')}</Button>
      </div>
    )
  }

  return (
    <div className="min-h-screen pb-8">
      <TopAppBar title="Refda" showBack />

      <main className="mt-24 px-container-margin space-y-section-gap">
        <section>
          <h2 className="text-headline-lg text-primary mb-4">{t('تفاصيل الهدية', 'Gift details')}</h2>
          <div className="bg-surface-container-lowest rounded-lg p-4 shadow-card flex gap-4">
            <div className="flex-1">
              <h3 className="text-headline-lg font-bold">{gift.name}</h3>
              <p className="text-body-sm text-on-surface-variant">{event.title}</p>
            </div>
            <p className="text-headline-lg text-primary">{formatCurrency(amount)}</p>
          </div>
        </section>

        {gift.type === 'cash' && (
          <section>
            <label className="text-label-md text-on-surface-variant">{t('المبلغ', 'Amount')}</label>
            <input
              type="number"
              min={50}
              step={50}
              value={amount}
              onChange={(e) => setAmount(Number(e.target.value))}
              className="w-full mt-2 py-4 px-4 bg-surface-container-low rounded-lg text-headline-lg"
            />
          </section>
        )}

        <section>
          <label className="text-label-md text-on-surface-variant">{t('رسالة إهداء', 'Gift message')}</label>
          <textarea
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            rows={3}
            className="w-full mt-2 py-4 px-4 bg-surface-container-low rounded-lg text-body-md resize-none"
            placeholder={t('اكتب رسالتك هنا...', 'Write your message...')}
          />
        </section>

        <label className="flex items-center gap-3">
          <input
            type="checkbox"
            checked={hideAmount}
            onChange={(e) => setHideAmount(e.target.checked)}
            className="w-5 h-5 rounded accent-primary"
          />
          <span className="text-body-sm">{t('إخفاء المبلغ عن المستفيد', 'Hide amount from beneficiary')}</span>
        </label>

        <section>
          <h2 className="text-headline-lg text-primary mb-4">{t('طريقة الدفع', 'Payment method')}</h2>
          <div className="space-y-3">
            {(['mada', 'apple_pay'] as const).map((method) => (
              <button
                key={method}
                type="button"
                onClick={() => setPaymentMethod(method)}
                className={`w-full flex items-center gap-4 p-4 rounded-lg border-2 transition-all ${
                  paymentMethod === method ? 'border-primary bg-primary-fixed/20' : 'border-outline-variant'
                }`}
              >
                <Icon name={method === 'mada' ? 'credit_card' : 'phone_iphone'} className="text-2xl" />
                <span className="text-body-md font-semibold">
                  {method === 'mada' ? 'mada' : 'Apple Pay'}
                </span>
              </button>
            ))}
          </div>
        </section>

        <Button size="lg" fullWidth onClick={handlePay} disabled={loading}>
          {t(`ادفع ${formatCurrency(amount)}`, `Pay ${formatCurrency(amount)}`)}
        </Button>
      </main>
    </div>
  )
}
