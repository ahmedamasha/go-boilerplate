import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import { AppShell } from '../components/layout/AppShell'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Icon } from '../components/ui/Icon'
import { useAuth } from '../context/AuthContext'
import { useLocale } from '../context/LocaleContext'

const eventTypes = [
  { id: 'wedding', labelAr: 'زواج', labelEn: 'Wedding', icon: 'favorite' },
  { id: 'newborn', labelAr: 'مولود', labelEn: 'Newborn', icon: 'child_care' },
  { id: 'birthday', labelAr: 'عيد ميلاد', labelEn: 'Birthday', icon: 'cake' },
  { id: 'other', labelAr: 'أخرى', labelEn: 'Other', icon: 'celebration' },
]

export function CreateEventPage() {
  const { t } = useLocale()
  const { token } = useAuth()
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [eventType, setEventType] = useState('wedding')
  const [eventDate, setEventDate] = useState('')
  const [privacy, setPrivacy] = useState<'public' | 'private'>('public')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!token) return
    setLoading(true)
    try {
      const event = await api.createEvent(token, {
        title,
        event_type: eventType,
        event_date: eventDate,
        privacy,
        gifts: [],
      })
      navigate(`/events/${event.id}`)
    } finally {
      setLoading(false)
    }
  }

  return (
    <AppShell title={t('إضافة مناسبة', 'Add event')} showBack>
      <div className="relative h-40 rounded-lg overflow-hidden bg-primary-container flex items-center justify-center">
        <p className="text-on-primary-container text-headline-lg">{t('وثّق لحظاتك السعيدة', 'Capture your happy moments')}</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-section-gap">
        <section>
          <h2 className="text-headline-lg text-primary mb-4">{t('نوع المناسبة', 'Event type')}</h2>
          <div className="flex gap-3 overflow-x-auto pb-2">
            {eventTypes.map((type) => (
              <button
                key={type.id}
                type="button"
                onClick={() => setEventType(type.id)}
                className={`flex flex-col items-center gap-2 min-w-[80px] p-4 rounded-lg transition-all ${
                  eventType === type.id
                    ? 'bg-primary text-on-primary'
                    : 'bg-surface-container-low text-on-surface-variant'
                }`}
              >
                <Icon name={type.icon} />
                <span className="text-label-md">{t(type.labelAr, type.labelEn)}</span>
              </button>
            ))}
          </div>
        </section>

        <Input
          label={t('اسم المناسبة', 'Event name')}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder={t('زواج منى وأحمد', 'Mona & Ahmed Wedding')}
          required
        />

        <Input
          label={t('تاريخ المناسبة', 'Event date')}
          type="date"
          value={eventDate}
          onChange={(e) => setEventDate(e.target.value)}
          required
        />

        <section className="space-y-3">
          <h3 className="text-headline-lg text-primary">{t('الخصوصية', 'Privacy')}</h3>
          <div className="flex gap-3">
            {(['public', 'private'] as const).map((p) => (
              <button
                key={p}
                type="button"
                onClick={() => setPrivacy(p)}
                className={`flex-1 py-3 rounded-lg text-label-md ${
                  privacy === p ? 'bg-primary text-on-primary' : 'bg-surface-container-low'
                }`}
              >
                {p === 'public' ? t('عام', 'Public') : t('خاص', 'Private')}
              </button>
            ))}
          </div>
        </section>

        <Button type="submit" size="lg" fullWidth disabled={loading}>
          <Icon name="add" />
          {t('إنشاء المناسبة', 'Create event')}
        </Button>
      </form>
    </AppShell>
  )
}
