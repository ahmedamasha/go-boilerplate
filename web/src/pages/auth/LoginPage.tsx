import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../../api/client'
import { Button } from '../../components/ui/Button'
import { Icon } from '../../components/ui/Icon'
import { Input } from '../../components/ui/Input'
import { useLocale } from '../../context/LocaleContext'
import { normalizeSaudiPhone } from '../../lib/format'
import { RefdaLogo } from '../../components/ui/RefdaLogo'

export function LoginPage() {
  const { t } = useLocale()
  const navigate = useNavigate()
  const [phone, setPhone] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const normalized = normalizeSaudiPhone(phone)
      await api.login(normalized)
      navigate('/auth/verify', { state: { phone: normalized, mode: 'login' } })
    } catch (err) {
      setError(err instanceof Error ? err.message : t('حدث خطأ', 'Something went wrong'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-container-margin">
      <header className="w-full max-w-md flex flex-col items-center mb-section-gap">
        <RefdaLogo size="lg" />
        <div className="text-center mt-4 space-y-1">
          <h2 className="text-headline-lg text-on-surface">{t('مرحباً بك', 'Welcome')}</h2>
          <p className="text-body-md text-on-surface-variant">{t('سجل الدخول للمتابعة', 'Sign in to continue')}</p>
        </div>
      </header>

      <main className="w-full max-w-md">
        <form
          onSubmit={handleSubmit}
          className="bg-surface-container-lowest shadow-card rounded-lg p-8 flex flex-col gap-stack-gap border border-surface-container-high"
        >
          <div className="w-full rounded-DEFAULT overflow-hidden desert-sand-gradient border border-secondary-container/60 py-8 flex flex-col items-center justify-center gap-2">
            <RefdaLogo size="md" showWordmark={false} className="!gap-0" />
            <p className="text-headline-lg font-bold text-primary">رفدة</p>
            <p className="text-body-sm text-on-surface-variant">{t('الإهداء الاجتماعي بثقة', 'Social gifting, trusted')}</p>
          </div>

          <Input
            label={t('رقم الجوال', 'Phone number')}
            type="tel"
            dir="ltr"
            adornmentSide="left"
            placeholder="5X XXX XXXX"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
            startAdornment={
              <>
                <div className="w-6 h-4 bg-[#006C35] rounded-sm" />
                <span className="text-on-surface font-semibold text-body-md">+966</span>
                <div className="h-6 w-px bg-outline-variant" />
              </>
            }
          />

          {error && <p className="text-body-sm text-error">{error}</p>}

          <Button type="submit" size="lg" fullWidth disabled={loading}>
            <span>{t('إرسال الرمز', 'Send code')}</span>
            <Icon name="arrow_forward" className="rtl-flip" />
          </Button>

          <div className="flex flex-col items-center gap-stack-gap">
            <Link to="/auth/register" className="text-body-sm text-primary font-bold hover:underline">
              {t('ليس لديك حساب؟ سجل الآن', "Don't have an account? Register")}
            </Link>
          </div>
        </form>
      </main>

      <footer className="mt-section-gap flex flex-col items-center gap-2 text-center">
        <p className="text-label-md text-outline">{t('خاضع لرقابة ساما', 'Regulated by SAMA')}</p>
        <div className="flex gap-stack-gap opacity-40">
          <Icon name="verified_user" className="text-on-surface-variant" />
          <Icon name="lock" className="text-on-surface-variant" />
          <Icon name="security" className="text-on-surface-variant" />
        </div>
      </footer>
    </div>
  )
}
