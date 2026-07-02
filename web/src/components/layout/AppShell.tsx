import type { ReactNode } from 'react'
import { BottomNav } from './BottomNav'
import { TopAppBar } from './TopAppBar'

interface AppShellProps {
  children: ReactNode
  title?: string
  showBack?: boolean
  showNav?: boolean
  showFab?: boolean
  onFabClick?: () => void
  leftSlot?: React.ReactNode
}

export function AppShell({
  children,
  title,
  showBack,
  showNav = true,
  showFab,
  onFabClick,
  leftSlot,
}: AppShellProps) {
  return (
    <div className="min-h-screen pb-32">
      <TopAppBar title={title} showBack={showBack} leftSlot={leftSlot} />
      <main className="mt-24 px-container-margin space-y-section-gap">{children}</main>
      {showFab && (
        <button
          type="button"
          onClick={onFabClick}
          className="fixed bottom-24 start-6 w-16 h-16 bg-primary text-on-primary rounded-full shadow-xl flex items-center justify-center hover:scale-105 active:scale-95 transition-transform z-40"
        >
          <span className="material-symbols-outlined text-[32px]">add</span>
        </button>
      )}
      {showNav && <BottomNav />}
    </div>
  )
}
