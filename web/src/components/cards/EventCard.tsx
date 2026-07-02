import { Link } from 'react-router-dom'
import { useLocale } from '../../context/LocaleContext'
import { formatCurrency, formatDate } from '../../lib/format'
import type { Event, EventStatus } from '../../api/types'
import { Icon } from '../ui/Icon'
import { ProgressBar } from '../ui/ProgressBar'
import { StatusChip } from '../ui/StatusChip'
import { images } from '../../lib/assets'

interface EventCardProps {
  event: Event
}

const statusLabels: Record<EventStatus, { ar: string; en: string; variant: 'active' | 'success' | 'warning' | 'reserved' | 'premium' }> = {
  pending_review: { ar: 'قيد المراجعة', en: 'Pending review', variant: 'warning' },
  open: { ar: 'نشط', en: 'Open', variant: 'active' },
  collected: { ar: 'مكتمل', en: 'Collected', variant: 'success' },
  partial_withdrawn: { ar: 'سحب جزئي', en: 'Partial withdrawal', variant: 'premium' },
  fully_withdrawn: { ar: 'تم السحب', en: 'Fully withdrawn', variant: 'reserved' },
  rejected: { ar: 'مرفوض', en: 'Rejected', variant: 'reserved' },
}

export function EventCard({ event }: EventCardProps) {
  const { t } = useLocale()
  const received = event.amount_collected ?? event.gifts.reduce((sum, g) => sum + (g.amount_received ?? 0), 0)
  const target = event.target_amount ?? (received || 1)
  const pct = Math.round((received / target) * 100)
  const status = statusLabels[event.status] ?? statusLabels.open

  return (
    <div className="bg-surface-container-lowest rounded-lg overflow-hidden shadow-card group">
      <div className="relative h-48 w-full overflow-hidden">
        <img
          src={event.hero_image_url || images.weddingHero}
          alt={event.title}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700"
        />
        <div className="absolute top-4 end-4">
          <StatusChip label={t(status.ar, status.en)} variant={status.variant} />
        </div>
      </div>
      <div className="p-6 space-y-stack-gap">
        <div className="space-y-1">
          <h4 className="text-headline-lg text-on-surface">{event.title}</h4>
          <p className="text-body-sm text-on-surface-variant flex items-center gap-1">
            <Icon name="calendar_today" className="text-[16px]" />
            {formatDate(event.event_date)}
          </p>
        </div>
        <ProgressBar
          value={received}
          max={target}
          labelStart={t(`المبلغ المحصل: ${formatCurrency(received)}`, `Collected: ${formatCurrency(received)}`)}
          labelEnd={t(`الهدف: ${formatCurrency(target)}`, `Goal: ${formatCurrency(target)}`)}
          hint={t(`تم تحقيق ${pct}٪ من الهدف`, `${pct}% of goal achieved`)}
        />
        {(event.amount_available ?? 0) > 0 && (
          <p className="text-body-sm text-primary font-medium">
            {t(
              `متاح للسحب: ${formatCurrency(event.amount_available ?? 0)}`,
              `Available to withdraw: ${formatCurrency(event.amount_available ?? 0)}`,
            )}
          </p>
        )}
        <div className="flex gap-stack-gap pt-2">
          <Link
            to={`/events/${event.id}`}
            className="flex-1 bg-primary text-on-primary text-label-md py-4 rounded-lg flex items-center justify-center gap-2 hover:opacity-90 active:scale-95 transition-all"
          >
            {t('إدارة القائمة', 'Manage registry')}
          </Link>
          <button
            type="button"
            className="flex items-center justify-center w-14 h-14 rounded-lg bg-secondary-container text-primary hover:bg-outline-variant/30 active:scale-95 transition-all"
          >
            <Icon name="share" />
          </button>
        </div>
      </div>
    </div>
  )
}
