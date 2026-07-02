import { NavLink } from 'react-router-dom'
import { Icon } from '../ui/Icon'
import { useLocale } from '../../context/LocaleContext'

const items = [
  { to: '/', icon: 'dashboard', labelAr: 'الرئيسية', labelEn: 'Home' },
  { to: '/events/new', icon: 'calendar_today', labelAr: 'مناسبات', labelEn: 'Events' },
  { to: '/gifts', icon: 'redeem', labelAr: 'هدايا', labelEn: 'Gifts' },
  { to: '/profile', icon: 'person', labelAr: 'حسابي', labelEn: 'Profile' },
]

export function BottomNav() {
  const { t } = useLocale()

  return (
    <nav className="fixed bottom-0 left-0 w-full z-50 flex justify-around items-center px-container-margin py-stack-gap bg-surface shadow-nav rounded-t-lg">
      {items.map((item) => (
        <NavLink
          key={item.to}
          to={item.to}
          end={item.to === '/'}
          className={({ isActive }) =>
            `flex flex-col items-center justify-center px-4 py-2 rounded-lg transition-all active:scale-95 ${
              isActive
                ? 'bg-primary-container text-on-primary-container'
                : 'text-on-surface-variant hover:bg-surface-container-low'
            }`
          }
        >
          <Icon name={item.icon} />
          <span className="text-label-md mt-1">{t(item.labelAr, item.labelEn)}</span>
        </NavLink>
      ))}
    </nav>
  )
}
