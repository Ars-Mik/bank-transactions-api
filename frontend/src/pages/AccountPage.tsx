import {
  useEffect,
  useState,
} from 'react'

import {
  Link,
  useParams,
} from 'react-router-dom'

import {
  getAccount,
  getAccountTransactions,
} from '../api/accounts'

import {
  getApiErrorMessage,
} from '../api/errors'

import type {
  Account,
  Transaction,
} from '../types/banking'

import {
  formatKopecks,
} from '../utils/money'

export function AccountPage() {
  const { id } = useParams()

  const accountId = Number(id)

  const [account, setAccount] =
    useState<Account | null>(null)

  const [
    transactions,
    setTransactions,
  ] = useState<Transaction[]>([])

  const [loading, setLoading] =
    useState(true)

  const [error, setError] =
    useState<string | null>(null)

  useEffect(() => {
    if (
      !Number.isInteger(accountId) ||
      accountId <= 0
    ) {
      setError(
        'Некорректный идентификатор счёта',
      )

      setLoading(false)

      return
    }

    async function loadAccount() {
      try {
        setLoading(true)
        setError(null)

        const [
          accountData,
          transactionData,
        ] = await Promise.all([
          getAccount(accountId),
          getAccountTransactions(
            accountId,
          ),
        ])

        setAccount(accountData)
        setTransactions(
          transactionData,
        )
      } catch (error) {
        setError(
          getApiErrorMessage(error),
        )
      } finally {
        setLoading(false)
      }
    }

    void loadAccount()
  }, [accountId])

  if (loading) {
    return (
      <main className="app">
        <section className="dashboard">
          <div className="state-message">
            Загружаем счёт...
          </div>
        </section>
      </main>
    )
  }

  if (error || !account) {
    return (
      <main className="app">
        <section className="dashboard">
          <Link
            className="back-link"
            to="/"
          >
            ← Назад к счетам
          </Link>

          <div
            className="error-message"
            role="alert"
          >
            {error ??
              'Счёт не найден'}
          </div>
        </section>
      </main>
    )
  }

  return (
    <main className="app">
      <section className="dashboard">
        <Link
          className="back-link"
          to="/"
        >
          ← Все счета
        </Link>

        <header className="account-page__header">
          <div>
            <p className="dashboard__eyebrow">
              Банковский счёт
            </p>

            <h1>
              Счёт #{account.id}
            </h1>
          </div>

          <div className="account-page__balance">
            <span>
              Текущий баланс
            </span>

            <strong>
              {formatKopecks(
                account.balance_kopecks,
              )}
            </strong>
          </div>
        </header>

        <section className="transactions-section">
          <div className="section-heading">
            <div>
              <h2>
                История операций
              </h2>

              <p>
                Все пополнения и переводы
                по счёту
              </p>
            </div>

            <span className="transactions-count">
              {transactions.length}
            </span>
          </div>

          {transactions.length === 0 ? (
            <div className="state-message">
              По этому счёту пока нет
              операций.
            </div>
          ) : (
            <div className="transactions-list">
              {transactions.map(
                (transaction) => (
                  <TransactionRow
                    key={transaction.id}
                    transaction={
                      transaction
                    }
                  />
                ),
              )}
            </div>
          )}
        </section>
      </section>
    </main>
  )
}

interface TransactionRowProps {
  transaction: Transaction
}

function TransactionRow({
  transaction,
}: TransactionRowProps) {
  const incoming =
    transaction.direction ===
    'incoming'

  const title =
    getTransactionTitle(
      transaction,
    )

  return (
    <article className="transaction-row">
      <div className="transaction-row__icon">
        {incoming ? '↓' : '↑'}
      </div>

      <div className="transaction-row__details">
        <strong>
          {title}
        </strong>

        <span>
          {formatTransactionDate(
            transaction.created_at,
          )}
        </span>
      </div>

      <div
        className={[
          'transaction-row__amount',
          incoming
            ? 'transaction-row__amount--incoming'
            : 'transaction-row__amount--outgoing',
        ].join(' ')}
      >
        {incoming ? '+' : '−'}
        {formatKopecks(
          transaction.amount_kopecks,
        )}
      </div>
    </article>
  )
}

function getTransactionTitle(
  transaction: Transaction,
): string {
  if (
    transaction.kind === 'deposit'
  ) {
    return 'Пополнение счёта'
  }

  if (
    transaction.direction ===
    'incoming'
  ) {
    return transaction.from_account_id
      ? `Перевод со счёта #${transaction.from_account_id}`
      : 'Входящий перевод'
  }

  return transaction.to_account_id
    ? `Перевод на счёт #${transaction.to_account_id}`
    : 'Исходящий перевод'
}

function formatTransactionDate(
  value: string,
): string {
  return new Intl.DateTimeFormat(
    'ru-RU',
    {
      dateStyle: 'medium',
      timeStyle: 'short',
    },
  ).format(
    new Date(value),
  )
}