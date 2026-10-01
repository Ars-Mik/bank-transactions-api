import { apiClient } from './client'
import type {
  Account,
  AccountsResponse,
} from '../types/banking'

export async function getAccounts(): Promise<Account[]> {
  const response =
    await apiClient.get<AccountsResponse>(
      '/accounts',
    )

  return response.data.accounts
}

export async function createAccount(): Promise<Account> {
  const response =
    await apiClient.post<Account>(
      '/accounts',
    )

  return response.data
}