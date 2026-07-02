import { Navigate, Outlet, Route, Routes } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { LoginPage } from '../pages/auth/LoginPage'
import { RegisterPage } from '../pages/auth/RegisterPage'
import { VerifyOtpPage } from '../pages/auth/VerifyOtpPage'
import { CheckoutPage } from '../pages/CheckoutPage'
import { CreateEventPage } from '../pages/CreateEventPage'
import { DashboardPage } from '../pages/DashboardPage'
import { EditProfilePage } from '../pages/EditProfilePage'
import { GifterExperiencePage } from '../pages/GifterExperiencePage'
import { ProfilePage } from '../pages/ProfilePage'

function ProtectedRoute() {
  const { token, loading } = useAuth()
  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-body-md text-on-surface-variant">...</p>
      </div>
    )
  }
  if (!token) return <Navigate to="/auth/login" replace />
  return <Outlet />
}

function PublicOnlyRoute() {
  const { token, loading } = useAuth()
  if (loading) return null
  if (token) return <Navigate to="/" replace />
  return <Outlet />
}

export function AppRoutes() {
  return (
    <Routes>
      <Route element={<PublicOnlyRoute />}>
        <Route path="/auth/login" element={<LoginPage />} />
        <Route path="/auth/register" element={<RegisterPage />} />
        <Route path="/auth/verify" element={<VerifyOtpPage />} />
      </Route>

      <Route path="/events/:id" element={<GifterExperiencePage />} />
      <Route path="/events/:id/checkout" element={<CheckoutPage />} />

      <Route element={<ProtectedRoute />}>
        <Route path="/" element={<DashboardPage />} />
        <Route path="/events/new" element={<CreateEventPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/profile/edit" element={<EditProfilePage />} />
        <Route path="/gifts" element={<Navigate to="/" replace />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
