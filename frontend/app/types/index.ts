export interface Subscription {
  id: string
  service_name: string
  price: number
  user_id: string
  start_date: string
  end_date?: string
  created_at: string
  updated_at: string
}

export interface SubscriptionListResponse {
  data: Subscription[]
  total: number
}

export interface TotalCostResponse {
  total_cost: number
}

export interface ServicesResponse {
  services: string[]
}
