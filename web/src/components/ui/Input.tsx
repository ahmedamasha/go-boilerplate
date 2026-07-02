import type { InputHTMLAttributes, ReactNode } from 'react'

interface InputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'prefix'> {
  label?: string
  startAdornment?: ReactNode
  /** Use "left" for dir="ltr" phone fields inside RTL pages */
  adornmentSide?: 'left' | 'right'
  error?: string
}

export function Input({
  label,
  startAdornment,
  adornmentSide = 'right',
  error,
  className = '',
  ...props
}: InputProps) {
  const adornmentClass =
    adornmentSide === 'left'
      ? 'absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none gap-2'
      : 'absolute inset-y-0 right-0 pr-4 flex items-center pointer-events-none gap-2'

  const inputPadClass =
    adornmentSide === 'left' ? 'pl-[7.5rem] pr-4 text-left' : 'pr-[7.5rem] pl-4'

  return (
    <div className="flex flex-col gap-1">
      {label && (
        <label className="text-label-md text-on-surface-variant uppercase tracking-widest">
          {label}
        </label>
      )}
      <div className="relative">
        {startAdornment && <div className={adornmentClass}>{startAdornment}</div>}
        <input
          className={`w-full py-4 bg-surface-container-low border-none rounded-DEFAULT text-on-surface text-body-md focus:ring-2 focus:ring-primary-container transition-all placeholder:text-outline-variant ${startAdornment ? inputPadClass : 'px-inline-padding'} ${className}`}
          {...props}
        />
      </div>
      {error && <p className="text-body-sm text-error">{error}</p>}
    </div>
  )
}
