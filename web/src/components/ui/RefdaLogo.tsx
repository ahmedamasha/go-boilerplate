import { Icon } from './Icon'

interface RefdaLogoProps {
  size?: 'sm' | 'md' | 'lg'
  showWordmark?: boolean
  className?: string
}

const sizes = {
  sm: { box: 'w-12 h-12', icon: 'text-2xl', word: 'text-headline-lg' },
  md: { box: 'w-16 h-16', icon: 'text-4xl', word: 'text-headline-xl' },
  lg: { box: 'w-20 h-20', icon: 'text-5xl', word: 'text-[2rem]' },
}

export function RefdaLogo({ size = 'md', showWordmark = true, className = '' }: RefdaLogoProps) {
  const s = sizes[size]

  return (
    <div className={`flex flex-col items-center gap-3 ${className}`}>
      <div
        className={`${s.box} bg-primary rounded-xl flex items-center justify-center shadow-lg -rotate-3 ring-4 ring-primary-fixed/40`}
      >
        <Icon name="card_giftcard" className={`text-on-primary ${s.icon} filled`} />
      </div>
      {showWordmark && (
        <span className={`${s.word} font-bold text-primary tracking-tight`}>رفدة</span>
      )}
    </div>
  )
}
