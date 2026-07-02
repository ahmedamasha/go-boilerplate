type ChipVariant = 'active' | 'success' | 'warning' | 'reserved' | 'premium'

const styles: Record<ChipVariant, string> = {
  active: 'bg-primary/90 text-on-primary',
  success: 'bg-primary-fixed/80 text-on-primary-fixed',
  warning: 'bg-tertiary-fixed/80 text-on-tertiary-fixed',
  reserved: 'bg-on-surface/80 text-white',
  premium: 'bg-secondary-container text-on-secondary-container',
}

interface StatusChipProps {
  label: string
  variant?: ChipVariant
}

export function StatusChip({ label, variant = 'active' }: StatusChipProps) {
  return (
    <span className={`px-3 py-1 rounded-full text-label-md backdrop-blur-sm ${styles[variant]}`}>
      {label}
    </span>
  )
}
