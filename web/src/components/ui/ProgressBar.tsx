interface ProgressBarProps {
  value: number
  max: number
  labelStart?: string
  labelEnd?: string
  hint?: string
}

export function ProgressBar({ value, max, labelStart, labelEnd, hint }: ProgressBarProps) {
  const pct = max > 0 ? Math.min(100, Math.round((value / max) * 100)) : 0

  return (
    <div className="space-y-2">
      {(labelStart || labelEnd) && (
        <div className="flex justify-between text-label-md text-primary">
          {labelStart && <span>{labelStart}</span>}
          {labelEnd && <span>{labelEnd}</span>}
        </div>
      )}
      <div className="h-3 w-full bg-secondary-container rounded-full overflow-hidden">
        <div className="h-full bg-primary rounded-full transition-all" style={{ width: `${pct}%` }} />
      </div>
      {hint && <p className="text-body-sm text-on-surface-variant text-center pt-1">{hint}</p>}
    </div>
  )
}
