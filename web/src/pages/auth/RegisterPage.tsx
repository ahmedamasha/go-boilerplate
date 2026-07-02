import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../../api/client'
import { Button } from '../../components/ui/Button'
import { Icon } from '../../components/ui/Icon'
import { Input } from '../../components/ui/Input'
import { useLocale } from '../../context/LocaleContext'
import { normalizeSaudiPhone } from '../../lib/format'

export function RegisterPage() {
  const { t } = useLocale()
  const navigate = useNavigate()
  const [fullName, setFullName] = useState('')
  const [phone, setPhone] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const normalized = normalizeSaudiPhone(phone)
      await api.register(normalized, fullName)
      navigate('/auth/verify', { state: { phone: normalized, fullName, mode: 'register' } })
    } catch (err) {
      setError(err instanceof Error ? err.message : t('حدث خطأ', 'Something went wrong'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex flex-col">
      <header className="px-container-margin py-4 flex justify-between items-center">
        <Link to="/auth/login">
          <Icon name="arrow_back" className="text-primary rtl-flip" />
        </Link>
        <span className="text-headline-xl text-primary">Refda</span>
        <div className="w-6" />
      </header>

      <main className="flex-1 px-container-margin pb-8">
        <div className="h-40 rounded-lg bg-primary-container/30 mb-section-gap desert-sand-gradient" />
        <h1 className="text-headline-xl text-primary mb-2">{t('انضم لرفدة', 'Join Refda')}</h1>
        <p className="text-body-md text-on-surface-variant mb-section-gap">
          {t('أنشئ حسابك وابدأ باستقبال الهدايا', 'Create your account and start receiving gifts')}
        </p>

        <form onSubmit={handleSubmit} className="space-y-stack-gap">
          <Input
            label={t('الاسم الكامل', 'Full name')}
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            required
          />
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
                <span className="text-on-surface font-semibold text-body-md">+966</span>
                <div className="h-6 w-px bg-outline-variant" />
              </>
            }
            required
          />

          {error && <p className="text-body-sm text-error">{error}</p>}

          <Button type="submit" size="lg" fullWidth disabled={loading}>
            {t('متابعة', 'Continue')}
            <Icon name="arrow_forward" className="rtl-flip" />
          </Button>
        </form>
      </main>
    </div>
  )
}
