import { useCallback, useEffect, useRef, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { api } from '../../api/client'
import { Button } from '../../components/ui/Button'
import { useAuth } from '../../context/AuthContext'
import { useLocale } from '../../context/LocaleContext'
import { OTP_LENGTH } from '../../lib/constants'
import { normalizeDigits } from '../../lib/format'

interface VerifyState {
  phone: string
  fullName?: string
  mode?: string
}

const emptyDigits = () => Array.from({ length: OTP_LENGTH }, () => '')

const RESEND_SECONDS = 59

export function VerifyOtpPage() {
  const { t, dir } = useLocale()
  const { setSession } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const state = location.state as VerifyState | null

  const [digits, setDigits] = useState(emptyDigits)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [countdown, setCountdown] = useState(RESEND_SECONDS)
  const inputs = useRef<(HTMLInputElement | null)[]>([])
  const verifying = useRef(false)

  const phone = state?.phone ?? ''
  const fullName = state?.fullName
  const mode = state?.mode
  const code = digits.join('')
  const codeComplete = code.length === OTP_LENGTH

  useEffect(() => {
    if (!state?.phone) navigate('/auth/login', { replace: true })
  }, [state, navigate])

  useEffect(() => {
    if (state?.phone) inputs.current[0]?.focus()
  }, [state?.phone])

  useEffect(() => {
    if (countdown <= 0) return
    const id = setInterval(() => setCountdown((c) => c - 1), 1000)
    return () => clearInterval(id)
  }, [countdown])

  const verify = useCallback(async (otpCode: string) => {
    if (!phone || otpCode.length !== OTP_LENGTH || verifying.current) return
    verifying.current = true
    setError('')
    setLoading(true)
    try {
      const res = await api.verify(phone, otpCode, fullName)
      setSession(res.token, res.user)
      navigate('/', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : t('رمز غير صحيح', 'Invalid code'))
      verifying.current = false
    } finally {
      setLoading(false)
    }
  }, [phone, fullName, setSession, navigate, t])

  useEffect(() => {
    if (codeComplete && !loading) verify(code)
  }, [code, codeComplete, loading, verify])

  if (!phone) return null

  function handleChange(index: number, raw: string) {
    const value = normalizeDigits(raw).replace(/\D/g, '').slice(-1)
    if (raw && !value) return
    const next = [...digits]
    next[index] = value
    setDigits(next)
    if (value && index < OTP_LENGTH - 1) inputs.current[index + 1]?.focus()
  }

  function handleKeyDown(index: number, e: React.KeyboardEvent) {
    if (e.key === 'Backspace' && !digits[index] && index > 0) {
      inputs.current[index - 1]?.focus()
    }
  }

  function handlePaste(e: React.ClipboardEvent) {
    e.preventDefault()
    const pasted = normalizeDigits(e.clipboardData.getData('text')).replace(/\D/g, '').slice(0, OTP_LENGTH)
    if (!pasted) return
    const next = emptyDigits()
    for (let i = 0; i < pasted.length; i++) next[i] = pasted[i]!
    setDigits(next)
    inputs.current[Math.min(pasted.length, OTP_LENGTH - 1)]?.focus()
  }

  async function handleResend() {
    if (countdown > 0 || loading) return
    setError('')
    setDigits(emptyDigits())
    verifying.current = false
    try {
      if (mode === 'register' && fullName) {
        await api.register(phone, fullName)
      } else {
        await api.login(phone)
      }
      setCountdown(RESEND_SECONDS)
      inputs.current[0]?.focus()
    } catch (err) {
      setError(err instanceof Error ? err.message : t('حدث خطأ', 'Something went wrong'))
    }
  }

  return (
    <div className="min-h-screen flex flex-col px-container-margin">
      <header className="py-4 flex justify-center">
        <div className="text-headline-lg text-primary">Refda</div>
      </header>

      <main className="flex-1 flex flex-col items-center justify-center max-w-md mx-auto w-full space-y-section-gap">
        <div className="text-center space-y-2">
          <h1 className="text-headline-xl text-primary">{t('التحقق من الهاتف', 'Verify phone')}</h1>
          <p className="text-body-md text-on-surface-variant">
            {t(
              `أدخل الرمز المكوّن من ${OTP_LENGTH} أرقام المرسل إلى ${phone}`,
              `Enter the ${OTP_LENGTH}-digit code sent to ${phone}`,
            )}
          </p>
        </div>

        <div
          className={`flex gap-4 justify-center ${dir === 'rtl' ? 'flex-row-reverse' : ''}`}
          dir="ltr"
        >
          {digits.map((d, i) => (
            <input
              key={i}
              ref={(el) => { inputs.current[i] = el }}
              className="otp-input w-16 h-20 text-center text-headline-xl bg-surface-container-lowest border-2 border-outline-variant rounded-lg focus:outline-none focus:border-primary transition-all"
              maxLength={1}
              value={d}
              onChange={(e) => handleChange(i, e.target.value)}
              onKeyDown={(e) => handleKeyDown(i, e)}
              onPaste={handlePaste}
              inputMode="numeric"
              autoComplete={i === 0 ? 'one-time-code' : 'off'}
              aria-label={t(`الرقم ${i + 1}`, `Digit ${i + 1}`)}
            />
          ))}
        </div>

        {error && <p className="text-body-sm text-error text-center">{error}</p>}

        <p className="text-body-sm text-on-surface-variant text-center">
          {countdown > 0 ? (
            t(`لم يصلك الرمز؟ أعد الإرسال خلال 0:${String(countdown).padStart(2, '0')}`, `Didn't receive it? Resend in 0:${String(countdown).padStart(2, '0')}`)
          ) : (
            <button type="button" onClick={handleResend} className="text-primary font-bold hover:underline">
              {t('إعادة إرسال الرمز', 'Resend code')}
            </button>
          )}
        </p>

        <Button size="lg" fullWidth onClick={() => verify(code)} disabled={loading || !codeComplete}>
          {mode === 'register' ? t('إنشاء الحساب', 'Create account') : t('تسجيل الدخول', 'Sign in')}
        </Button>
      </main>
    </div>
  )
}
