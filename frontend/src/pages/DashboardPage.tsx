import {
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react'

import {
  createAccount,
  getAccounts,
} from '../api/accounts'

import {
  getApiErrorMessage,
} from '../api/errors'

import {
  AccountCard,
} from '../components/AccountCard'

import type {
  Account,
} from '../types/banking'

import {
  formatKopecks,
} from '../utils/money'

export function DashboardPage() {
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
          getApiErrorMessage(error),
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
        (total, account) =>
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
        getApiErrorMessage(error),
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
            onClick={handleCreateAccount}
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