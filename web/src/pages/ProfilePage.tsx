import { useNavigate } from 'react-router-dom'
import { AppShell } from '../components/layout/AppShell'
import { Icon } from '../components/ui/Icon'
import { StatusChip } from '../components/ui/StatusChip'
import { useAuth } from '../context/AuthContext'
import { useLocale } from '../context/LocaleContext'
import { images } from '../lib/assets'

const menuItems = [
  { icon: 'person', labelAr: 'المعلومات الشخصية', labelEn: 'Personal info', path: '/profile/edit' },
  { icon: 'credit_card', labelAr: 'طرق الدفع', labelEn: 'Payment methods' },
  { icon: 'notifications', labelAr: 'الإشعارات', labelEn: 'Notifications' },
  { icon: 'help', labelAr: 'المساعدة والدعم', labelEn: 'Help & support' },
  { icon: 'info', labelAr: 'عن رفدة', labelEn: 'About Refda' },
]

export function ProfilePage() {
  const { t } = useLocale()
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  return (
    <AppShell title="Refda" showBack>
      <section className="flex flex-col items-center mt-stack-gap">
        <div className="relative">
          <div className="w-32 h-32 rounded-full border-4 border-secondary-container p-1 overflow-hidden">
            <img
              src={user?.profile_picture_url || images.profileAvatar}
              alt=""
              className="w-full h-full object-cover rounded-full"
            />
          </div>
          <div className="absolute bottom-1 end-1 bg-primary text-white p-1.5 rounded-full border-2 border-white">
            <Icon name="verified" className="text-sm filled" />
          </div>
        </div>
        <div className="text-center mt-stack-gap">
          <h2 className="text-headline-lg text-primary">{user?.full_name || t('مستخدم', 'User')}</h2>
          <div className="flex items-center justify-center gap-1 mt-1">
            <StatusChip label={t('عضو مميز', 'Premium Member')} variant="premium" />
          </div>
          <p className="text-body-sm text-outline mt-1">{user?.location || t('الرياض، السعودية', 'Riyadh, KSA')}</p>
        </div>
      </section>

      <section className="grid grid-cols-2 gap-stack-gap">
        <div className="bg-surface-container-lowest shadow-card p-stack-gap rounded-lg flex flex-col items-center text-center">
          <span className="text-headline-lg text-primary">٢٤</span>
          <span className="text-label-md text-outline">{t('هدايا مرسلة', 'Gifts sent')}</span>
        </div>
        <div className="bg-surface-container-lowest shadow-card p-stack-gap rounded-lg flex flex-col items-center text-center">
          <span className="text-headline-lg text-primary">٠٨</span>
          <span className="text-label-md text-outline">{t('مناسبات', 'Events')}</span>
        </div>
      </section>

      <section className="bg-surface-container-lowest rounded-lg overflow-hidden shadow-card">
        {menuItems.map((item, i) => (
          <button
            key={item.icon}
            type="button"
            onClick={() => item.path && navigate(item.path)}
            className={`w-full flex items-center gap-4 px-5 py-4 hover:bg-surface-container-low transition-colors ${
              i < menuItems.length - 1 ? 'border-b border-surface-container' : ''
            }`}
          >
            <Icon name={item.icon} className="text-primary" />
            <span className="text-body-md flex-1 text-start">{t(item.labelAr, item.labelEn)}</span>
            <Icon name="chevron_left" className="text-outline rtl-flip" />
          </button>
        ))}
      </section>

      <button
        type="button"
        onClick={() => { logout(); navigate('/auth/login') }}
        className="w-full py-4 text-error text-body-md font-semibold flex items-center justify-center gap-2"
      >
        <Icon name="logout" />
        {t('تسجيل الخروج', 'Log out')}
      </button>
    </AppShell>
  )
}
