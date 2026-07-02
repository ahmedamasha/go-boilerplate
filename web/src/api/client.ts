import type { AuthResponse, Event, EventStatus, PaymentResponse, User } from './types'

const API_BASE = '/api/v1'

class ApiError extends Error {
  status: number
  messageAr?: string

  constructor(message: string, status: number, messageAr?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.messageAr = messageAr
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {},
  token?: string | null,
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  if (token) headers.Authorization = `Bearer ${token}`

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })
  const data = await res.json().catch(() => ({}))

  if (!res.ok) {
    const err = data.error ?? data
    throw new ApiError(
      err.message || 'Request failed',
      res.status,
      err.message_ar,
    )
  }
  return data.data ?? data
}

export const api = {
  login: (phone: string) =>
    request<{ message: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ phone }),
    }),

  register: (phone: string, fullName: string) =>
    request<{ message: string }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ phone, full_name: fullName }),
    }),

  verify: (phone: string, code: string, fullName?: string | undefined) =>
    request<AuthResponse>('/auth/verify', {
      method: 'POST',
      body: JSON.stringify({ phone, code, ...(fullName ? { full_name: fullName } : {}) }),
    }),

  getProfile: (token: string) =>
    request<User>('/users/me', {}, token),

  updateProfile: (token: string, data: { full_name?: string; location?: string }) =>
    request<User>('/users/me', { method: 'PUT', body: JSON.stringify(data) }, token),

  getEvents: (token?: string | null, mine = false) =>
    request<Event[]>(`/events${mine ? '?mine=true' : ''}`, {}, token),

  getEvent: (id: string) => request<Event>(`/events/${id}`),

  createEvent: (token: string, data: Record<string, unknown>) =>
    request<Event>('/events', { method: 'POST', body: JSON.stringify(data) }, token),

  withdraw: (token: string, eventId: string, amount: number) =>
    request<{
      withdrawal_id: string
      amount: number
      status: string
      event_status: EventStatus
      amount_collected: number
      amount_withdrawn: number
      amount_available: number
    }>(`/events/${eventId}/withdraw`, {
      method: 'POST',
      body: JSON.stringify({ amount }),
    }, token),

  contribute: (data: Record<string, unknown>, token?: string | null) =>
    request<{ id: string }>('/gifts/contribute', {
      method: 'POST',
      body: JSON.stringify(data),
    }, token),

  pay: (contributionId: string, paymentMethod: 'mada' | 'apple_pay') =>
    request<PaymentResponse>('/gifts/pay', {
      method: 'POST',
      body: JSON.stringify({ contribution_id: contributionId, payment_method: paymentMethod }),
    }),
}

export { ApiError }
