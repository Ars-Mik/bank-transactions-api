import { apiClient } from './client'

import type {
  Account,
  TransferResult,
} from '../types/banking'

export async function depositToAccount(
  accountId: number,
  amountKopecks: number,
): Promise<Account> {
  const response =
    await apiClient.post<Account>(
      `/accounts/${accountId}/deposit`,
      {
        amount_kopecks: amountKopecks,
      },
    )

  return response.data
}

export async function transferBetweenAccounts(
  fromAccountId: number,
  toAccountId: number,
  amountKopecks: number,
): Promise<TransferResult> {
  const response =
    await apiClient.post<TransferResult>(
      '/transfers',
      {
        from_account_id: fromAccountId,
        to_account_id: toAccountId,
        amount_kopecks: amountKopecks,
      },
    )

  return response.data
}