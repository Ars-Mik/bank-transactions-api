import { apiClient } from './client'

import type {
  Account,
  AccountsResponse,
  Transaction,
  TransactionsResponse,
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

export async function getAccount(
  accountId: number,
): Promise<Account> {
  const response =
    await apiClient.get<Account>(
      `/accounts/${accountId}`,
    )

  return response.data
}

export async function getAccountTransactions(
  accountId: number,
): Promise<Transaction[]> {
  const response =
    await apiClient.get<TransactionsResponse>(
      `/accounts/${accountId}/transactions`,
    )

  return response.data.transactions
}