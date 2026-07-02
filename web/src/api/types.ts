export interface User {
  id: string
  phone: string
  full_name: string
  location?: string
  profile_picture_url?: string
}

export interface AuthResponse {
  token: string
  user: User
  is_new_user: boolean
}

export interface Gift {
  id: string
  name: string
  type: 'product' | 'cash'
  image_urls?: string[]
  target_amount?: number
  amount_received?: number
  status: string
  reserved?: boolean
}

export type EventStatus =
  | 'pending_review'
  | 'open'
  | 'collected'
  | 'partial_withdrawn'
  | 'fully_withdrawn'
  | 'rejected'

export interface Event {
  id: string
  title: string
  event_type: string
  event_date: string
  event_time?: string
  target_amount?: number
  privacy: string
  status: EventStatus
  amount_collected?: number
  amount_withdrawn?: number
  amount_available?: number
  hero_image_url?: string
  gifts: Gift[]
  beneficiary?: User
}

export interface Contribution {
  id: string
  amount: number
  status: string
  reference_number: string
  hide_amount: boolean
}

export interface PaymentResponse {
  reference_number: string
  status: string
  gift_card?: Record<string, unknown>
  contribution: Contribution
}
