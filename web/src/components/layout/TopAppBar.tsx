import { useNavigate } from 'react-router-dom'
import { Icon } from '../ui/Icon'

interface TopAppBarProps {
  title?: string
  showBack?: boolean
  showNotifications?: boolean
  leftSlot?: React.ReactNode
  rightSlot?: React.ReactNode
  sticky?: boolean
}

export function TopAppBar({
  title = 'رفدة',
  showBack = false,
  showNotifications = true,
  leftSlot,
  rightSlot,
  sticky = true,
}: TopAppBarProps) {
  const navigate = useNavigate()

  return (
    <header
      className={`bg-surface z-50 flex justify-between items-center w-full px-container-margin py-4 ${sticky ? 'fixed top-0 left-0 right-0' : 'sticky top-0'}`}
    >
      <div className="flex items-center gap-3 min-w-10">
        {showBack ? (
          <button
            type="button"
            onClick={() => navigate(-1)}
            className="flex items-center justify-center active:scale-95 transition-transform"
          >
            <Icon name="arrow_back" className="text-primary text-headline-lg rtl-flip" />
          </button>
        ) : (
          leftSlot
        )}
      </div>
      <h1 className="text-headline-xl font-bold text-primary">{title}</h1>
      <div className="flex items-center gap-3 min-w-10 justify-end">
        {rightSlot ??
          (showNotifications && (
            <button
              type="button"
              className="w-10 h-10 flex items-center justify-center rounded-full hover:bg-surface-container-high transition-colors active:scale-95"
            >
              <Icon name="notifications" className="text-primary" />
            </button>
          ))}
      </div>
    </header>
  )
}
