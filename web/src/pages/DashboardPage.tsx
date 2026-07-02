import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import type { Event } from '../api/types'
import { EventCard } from '../components/cards/EventCard'
import { AppShell } from '../components/layout/AppShell'
import { Icon } from '../components/ui/Icon'
import { useAuth } from '../context/AuthContext'
import { useLocale } from '../context/LocaleContext'
import { images } from '../lib/assets'

export function DashboardPage() {
  const { t } = useLocale()
  const { user, token } = useAuth()
  const navigate = useNavigate()
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) return
    api
      .getEvents(token, true)
      .then(setEvents)
      .catch(() => setEvents([]))
      .finally(() => setLoading(false))
  }, [token])

  const firstName = user?.full_name?.split(' ')[0] ?? t('صديقي', 'friend')

  return (
    <AppShell
      showFab
      onFabClick={() => navigate('/events/new')}
      leftSlot={
        <div className="w-10 h-10 rounded-full overflow-hidden border-2 border-primary-fixed">
          <img
            src={user?.profile_picture_url || images.userAvatar}
            alt=""
            className="w-full h-full object-cover"
          />
        </div>
      }
    >
      <section className="space-y-1">
        <p className="text-label-md text-on-surface-variant">
          {t(`أهلاً بك، ${firstName}`, `Welcome, ${firstName}`)}
        </p>
        <h2 className="text-headline-lg text-primary">{t('لوحة التحكم الخاصة بك', 'Your dashboard')}</h2>
      </section>

      <section className="space-y-stack-gap">
        <div className="flex justify-between items-end">
          <h3 className="text-headline-lg text-on-surface">{t('مناسباتي', 'My events')}</h3>
          <span className="text-label-md text-primary">{t('عرض الكل', 'View all')}</span>
        </div>

        {loading ? (
          <div className="text-body-md text-on-surface-variant text-center py-8">{t('جاري التحميل...', 'Loading...')}</div>
        ) : events.length > 0 ? (
          events.map((event) => <EventCard key={event.id} event={event} />)
        ) : (
          <div className="bg-surface-container-lowest rounded-lg p-8 text-center shadow-card space-y-4">
            <Icon name="celebration" className="text-primary text-4xl" />
            <p className="text-body-md text-on-surface-variant">
              {t('لا توجد مناسبات بعد. أنشئ أول مناسبة!', 'No events yet. Create your first one!')}
            </p>
            <button
              type="button"
              onClick={() => navigate('/events/new')}
              className="bg-primary text-on-primary px-6 py-3 rounded-full text-label-md"
            >
              {t('إضافة مناسبة', 'Add event')}
            </button>
          </div>
        )}
      </section>

      <section className="bg-primary-container text-on-primary-container p-6 rounded-lg space-y-stack-gap relative overflow-hidden">
        <div className="relative z-10 space-y-2">
          <h4 className="text-headline-lg">{t('شارك المناسبة', 'Share your event')}</h4>
          <p className="text-body-sm opacity-90 max-w-[240px]">
            {t(
              'اجعل من السهل على أحبائك اختيار الهدية المثالية من خلال مشاركة الرابط',
              'Make it easy for loved ones to pick the perfect gift by sharing your link',
            )}
          </p>
          <div className="pt-4">
            <button
              type="button"
              className="bg-[#25D366] text-white px-6 py-3 rounded-lg text-label-md flex items-center gap-3 hover:shadow-lg transition-shadow active:scale-95"
            >
              <Icon name="qr_code_2" />
              {t('مشاركة عبر واتساب', 'Share via WhatsApp')}
            </button>
          </div>
        </div>
        <div className="absolute -bottom-6 -start-6 w-32 h-32 bg-on-primary-container/10 rounded-full blur-2xl" />
      </section>

      <section className="grid grid-cols-2 gap-stack-gap">
        <div className="bg-surface-container-low p-5 rounded-lg space-y-2">
          <Icon name="favorite" className="text-primary" />
          <p className="text-headline-lg text-on-surface">١٢</p>
          <p className="text-label-md text-on-surface-variant">{t('مساهمين جدد', 'New contributors')}</p>
        </div>
        <div className="bg-surface-container-low p-5 rounded-lg space-y-2">
          <Icon name="redeem" className="text-primary" />
          <p className="text-headline-lg text-on-surface">٨</p>
          <p className="text-label-md text-on-surface-variant">{t('هدايا محجوزة', 'Reserved gifts')}</p>
        </div>
      </section>
    </AppShell>
  )
}
