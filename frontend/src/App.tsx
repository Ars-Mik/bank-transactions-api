import {
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react'

import axios from 'axios'

import {
  createAccount,
  getAccounts,
} from './api/accounts'

import { AccountCard } from './components/AccountCard'

import type {
  Account,
  ApiErrorResponse,
} from './types/banking'

import { formatKopecks } from './utils/money'

import './App.css'

function App() {
  const [accounts, setAccounts] =
    useState<Account[]>([])

  const [loading, setLoading] =
    useState(true)

  const [creating, setCreating] =
    useState(false)

  const [error, setError] =
    useState<string | null>(null)

  const loadAccounts = useCallback(
    async () => {
      try {
        setError(null)

        const data =
          await getAccounts()

        setAccounts(data)
      } catch (error) {
        setError(
          getErrorMessage(error),
        )
      } finally {
        setLoading(false)
      }
    },
    [],
  )

  useEffect(() => {
    void loadAccounts()
  }, [loadAccounts])

  const totalBalance = useMemo(
    () =>
      accounts.reduce(
        (
          total,
          account,
        ) =>
          total +
          account.balance_kopecks,
        0,
      ),
    [accounts],
  )

  async function handleCreateAccount() {
    if (creating) {
      return
    }

    try {
      setCreating(true)
      setError(null)

      const account =
        await createAccount()

      setAccounts(
        (currentAccounts) => [
          ...currentAccounts,
          account,
        ],
      )
    } catch (error) {
      setError(
        getErrorMessage(error),
      )
    } finally {
      setCreating(false)
    }
  }

  return (
    <main className="app">
      <section className="dashboard">
        <header className="dashboard__header">
          <div>
            <p className="dashboard__eyebrow">
              Bank Transactions
            </p>

            <h1>
              Мои счета
            </h1>

            <p className="dashboard__subtitle">
              Управление банковскими счетами
              и операциями
            </p>
          </div>

          <button
            className="primary-button"
            type="button"
            disabled={creating}
            onClick={
              handleCreateAccount
            }
          >
            {creating
              ? 'Создание...'
              : '+ Создать счёт'}
          </button>
        </header>

        <section className="summary-card">
          <span>
            Общий баланс
          </span>

          <strong>
            {formatKopecks(
              totalBalance,
            )}
          </strong>

          <small>
            Счетов: {accounts.length}
          </small>
        </section>

        {error && (
          <div
            className="error-message"
            role="alert"
          >
            {error}
          </div>
        )}

        {loading ? (
          <div className="state-message">
            Загружаем счета...
          </div>
        ) : accounts.length === 0 ? (
          <div className="state-message">
            У вас пока нет банковских
            счетов.
          </div>
        ) : (
          <section className="accounts-grid">
            {accounts.map(
              (account) => (
                <AccountCard
                  key={account.id}
                  account={account}
                />
              ),
            )}
          </section>
        )}
      </section>
    </main>
  )
}

function getErrorMessage(
  error: unknown,
): string {
  if (axios.isAxiosError<ApiErrorResponse>(
    error,
  )) {
    return (
      error.response?.data.error.message ??
      'Не удалось выполнить запрос'
    )
  }

  return 'Произошла неизвестная ошибка'
}

export default App