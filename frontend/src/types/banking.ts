export interface Account {
  id: number
  balance_kopecks: number
}

export interface AccountsResponse {
  accounts: Account[]
}

export interface ApiErrorResponse {
  error: {
    code: string
    message: string
  }
}