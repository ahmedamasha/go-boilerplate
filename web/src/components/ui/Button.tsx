import type { ButtonHTMLAttributes, ReactNode } from 'react'

type Variant = 'primary' | 'secondary' | 'ghost' | 'danger'
type Size = 'sm' | 'md' | 'lg'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant
  size?: Size
  fullWidth?: boolean
  children: ReactNode
}

const variants: Record<Variant, string> = {
  primary: 'bg-primary text-on-primary hover:bg-primary-container shadow-md',
  secondary: 'bg-secondary-container text-primary hover:bg-outline-variant/30',
  ghost: 'bg-transparent text-primary hover:bg-surface-container-low',
  danger: 'bg-error-container text-on-error-container hover:bg-error/10',
}

const sizes: Record<Size, string> = {
  sm: 'px-4 py-2 text-label-md rounded-lg',
  md: 'px-6 py-3 text-body-md rounded-lg',
  lg: 'px-8 py-4 text-headline-lg rounded-full',
}

export function Button({
  variant = 'primary',
  size = 'md',
  fullWidth = false,
  className = '',
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      className={`inline-flex items-center justify-center gap-2 font-semibold transition-all active:scale-95 disabled:opacity-50 disabled:pointer-events-none ${variants[variant]} ${sizes[size]} ${fullWidth ? 'w-full' : ''} ${className}`}
      {...props}
    >
      {children}
    </button>
  )
}
