export interface Account {
  id: number
  balance_kopecks: number
}

export interface AccountsResponse {
  accounts: Account[]
}

export type TransactionKind =
  | 'deposit'
  | 'transfer'

export type TransactionDirection =
  | 'incoming'
  | 'outgoing'

export interface Transaction {
  id: number
  kind: TransactionKind
  direction: TransactionDirection
  amount_kopecks: number
  from_account_id: number | null
  to_account_id: number | null
  created_at: string
}

export interface TransactionsResponse {
  account_id: number
  transactions: Transaction[]
}

export interface ApiErrorResponse {
  error: {
    code: string
    message: string
  }
}

export interface TransferResult {
  transaction_id: number
  from_account: Account
  to_account: Account
}