import { useLocale } from '../../context/LocaleContext'
import { formatCurrency } from '../../lib/format'
import type { Gift } from '../../api/types'
import { Button } from '../ui/Button'
import { StatusChip } from '../ui/StatusChip'
import { ProgressBar } from '../ui/ProgressBar'

interface GiftCardProps {
  gift: Gift
  onContribute?: (gift: Gift) => void
}

export function GiftCard({ gift, onContribute }: GiftCardProps) {
  const { t } = useLocale()
  const received = gift.amount_received ?? 0
  const target = gift.target_amount ?? 0
  const image = gift.image_urls?.[0]

  if (gift.reserved) {
    return (
      <div className="bg-surface-container-lowest rounded-lg overflow-hidden shadow-card relative">
        {image && (
          <div className="relative h-40">
            <img src={image} alt={gift.name} className="w-full h-full object-cover opacity-60" />
            <div className="absolute inset-0 flex items-center justify-center">
              <StatusChip label={t('محجوز', 'Reserved')} variant="reserved" />
            </div>
          </div>
        )}
        <div className="p-4">
          <h3 className="text-headline-lg mb-1">{gift.name}</h3>
          <p className="text-body-sm text-on-surface-variant">{t('تم حجز هذه الهدية', 'This gift is reserved')}</p>
        </div>
      </div>
    )
  }

  const isGroup = gift.type === 'product' && target > 0 && received < target

  return (
    <div className="bg-surface-container-lowest rounded-lg overflow-hidden shadow-card">
      {image && (
        <div className="h-40 overflow-hidden">
          <img src={image} alt={gift.name} className="w-full h-full object-cover" />
        </div>
      )}
      <div className="p-4 space-y-3">
        <h3 className="text-headline-lg">{gift.name}</h3>
        {isGroup ? (
          <>
            <ProgressBar
              value={received}
              max={target}
              labelStart={formatCurrency(received)}
              labelEnd={formatCurrency(target)}
            />
            <Button size="sm" fullWidth onClick={() => onContribute?.(gift)}>
              {t('ساهم جزئياً', 'Contribute')}
            </Button>
          </>
        ) : (
          <>
            <p className="text-body-sm text-on-surface-variant">
              {gift.type === 'cash'
                ? t('مساهمة نقدية', 'Cash gift')
                : formatCurrency(target)}
            </p>
            <Button size="sm" fullWidth onClick={() => onContribute?.(gift)}>
              {t('إهداء', 'Gift')}
            </Button>
          </>
        )}
      </div>
    </div>
  )
}
