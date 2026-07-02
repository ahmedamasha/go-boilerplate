import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { api } from '../api/client'
import type { Event, Gift } from '../api/types'
import { GiftCard } from '../components/cards/GiftCard'
import { TopAppBar } from '../components/layout/TopAppBar'
import { Button } from '../components/ui/Button'
import { Icon } from '../components/ui/Icon'
import { useLocale } from '../context/LocaleContext'
import { formatDate } from '../lib/format'
import { images } from '../lib/assets'

export function GifterExperiencePage() {
  const { id } = useParams<{ id: string }>()
  const { t } = useLocale()
  const navigate = useNavigate()
  const [event, setEvent] = useState<Event | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!id) return
    api
      .getEvent(id)
      .then(setEvent)
      .catch(() => setEvent(null))
      .finally(() => setLoading(false))
  }, [id])

  function handleContribute(gift: Gift) {
    navigate(`/events/${id}/checkout`, { state: { gift, event } })
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-body-md text-on-surface-variant">{t('جاري التحميل...', 'Loading...')}</p>
      </div>
    )
  }

  if (!event) {
    return (
      <div className="min-h-screen flex items-center justify-center px-container-margin text-center">
        <p className="text-body-md text-on-surface-variant">{t('المناسبة غير موجودة', 'Event not found')}</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen pb-32">
      <TopAppBar title="Refda" showBack showNotifications={false} />

      <div className="mt-24">
        <div className="relative h-56 overflow-hidden">
          <img
            src={event.hero_image_url || images.weddingHero}
            alt={event.title}
            className="w-full h-full object-cover"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-background via-transparent to-transparent" />
          <div className="absolute bottom-0 px-container-margin pb-6">
            <h2 className="text-headline-xl text-on-surface mb-2">{event.title}</h2>
            <p className="text-body-sm text-on-surface-variant flex items-center gap-1">
              <Icon name="calendar_today" className="text-[16px]" />
              {formatDate(event.event_date)}
            </p>
          </div>
        </div>

        <main className="px-container-margin space-y-section-gap mt-section-gap">
          <section>
            <h2 className="text-headline-xl text-primary mb-4">
              {t('تجربة المهدي', 'Gifter experience')}
            </h2>

            <div className="bg-secondary-container rounded-lg p-5 mb-section-gap">
              <h3 className="text-headline-lg text-on-secondary-container mb-1">
                {t('مساهمة نقدية مباشرة', 'Direct cash contribution')}
              </h3>
              <p className="text-body-sm text-on-secondary-container/80 mb-4">
                {t('ادعم أحباءك بمساهمة نقدية مرنة', 'Support loved ones with a flexible cash gift')}
              </p>
              <Button
                variant="primary"
                onClick={() =>
                  handleContribute({
                    id: 'cash',
                    name: t('مساهمة نقدية', 'Cash gift'),
                    type: 'cash',
                    status: 'available',
                  })
                }
              >
                {t('إهداء نقدي', 'Send cash')}
              </Button>
            </div>
          </section>

          <section className="flex gap-2 overflow-x-auto pb-2">
            {['الكل', 'متاح', 'قيد التجميع', 'محجوز'].map((filter, i) => (
              <button
                key={filter}
                type="button"
                className={`px-4 py-2 rounded-full text-label-md whitespace-nowrap ${
                  i === 0 ? 'bg-primary text-on-primary' : 'bg-surface-container-low text-on-surface-variant'
                }`}
              >
                {filter}
              </button>
            ))}
          </section>

          <section className="grid gap-stack-gap">
            {event.gifts.map((gift) => (
              <GiftCard key={gift.id} gift={gift} onContribute={handleContribute} />
            ))}
          </section>

          <section className="bg-primary-container text-on-primary-container rounded-lg p-6 relative overflow-hidden">
            <h2 className="text-headline-xl mb-3">{t('شاركهم فرحة البداية', 'Share their joy')}</h2>
            <p className="text-body-sm opacity-90 mb-4">
              {t('كل هدية تقربهم من بداية حياتهم الجديدة', 'Every gift brings them closer to their new beginning')}
            </p>
            <Button variant="primary" className="bg-on-primary text-primary">
              {t('إهداء الآن', 'Gift now')}
            </Button>
          </section>
        </main>
      </div>
    </div>
  )
}
