import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import { AppShell } from '../components/layout/AppShell'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { useAuth } from '../context/AuthContext'
import { useLocale } from '../context/LocaleContext'
import { images } from '../lib/assets'

export function EditProfilePage() {
  const { t } = useLocale()
  const { user, token, refreshProfile } = useAuth()
  const navigate = useNavigate()
  const [fullName, setFullName] = useState(user?.full_name ?? '')
  const [location, setLocation] = useState(user?.location ?? '')
  const [loading, setLoading] = useState(false)

  async function handleSave() {
    if (!token) return
    setLoading(true)
    try {
      await api.updateProfile(token, { full_name: fullName, location })
      await refreshProfile()
      navigate('/profile')
    } finally {
      setLoading(false)
    }
  }

  return (
    <AppShell title={t('تعديل الملف الشخصي', 'Edit profile')} showBack showNav={false}>
      <section className="flex flex-col items-center">
        <div className="w-28 h-28 rounded-full border-4 border-secondary-container overflow-hidden">
          <img src={user?.profile_picture_url || images.profileAvatar} alt="" className="w-full h-full object-cover" />
        </div>
        <h2 className="text-headline-lg text-primary mt-4">{fullName}</h2>
      </section>

      <section className="space-y-stack-gap">
        <Input
          label={t('الاسم الكامل', 'Full name')}
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
        />
        <Input
          label={t('رقم الجوال', 'Phone')}
          value={user?.phone ?? ''}
          disabled
          dir="ltr"
        />
        <Input
          label={t('الموقع', 'Location')}
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          placeholder={t('الرياض، السعودية', 'Riyadh, KSA')}
        />
      </section>

      <Button size="lg" fullWidth onClick={handleSave} disabled={loading}>
        {t('حفظ التغييرات', 'Save changes')}
      </Button>
    </AppShell>
  )
}
