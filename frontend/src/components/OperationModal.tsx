import {
  useEffect,
  useMemo,
  useState,
} from 'react'

import type {
  FormEvent,
} from 'react'

import {
  depositToAccount,
  transferBetweenAccounts,
} from '../api/operations'

import {
  getApiErrorMessage,
} from '../api/errors'

import type {
  Account,
} from '../types/banking'

import {
  formatKopecks,
  parseRublesToKopecks,
} from '../utils/money'

export type OperationMode =
  | 'deposit'
  | 'transfer'

interface OperationModalProps {
  mode: OperationMode
  account: Account
  accounts: Account[]
  onClose: () => void
  onCompleted: () => Promise<void>
}

export function OperationModal({
  mode,
  account,
  accounts,
  onClose,
  onCompleted,
}: OperationModalProps) {
  const [amount, setAmount] =
    useState('')

  const [
    destinationAccountId,
    setDestinationAccountId,
  ] = useState('')

  const [submitting, setSubmitting] =
    useState(false)

  const [error, setError] =
    useState<string | null>(null)

  const destinationAccounts =
    useMemo(
      () =>
        accounts.filter(
          (item) =>
            item.id !== account.id,
        ),
      [accounts, account.id],
    )

  useEffect(() => {
    function handleKeyDown(
      event: KeyboardEvent,
    ) {
      if (
        event.key === 'Escape' &&
        !submitting
      ) {
        onClose()
      }
    }

    window.addEventListener(
      'keydown',
      handleKeyDown,
    )

    return () => {
      window.removeEventListener(
        'keydown',
        handleKeyDown,
      )
    }
  }, [onClose, submitting])

  async function handleSubmit(
    event: FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault()

    const amountKopecks =
      parseRublesToKopecks(amount)

    if (amountKopecks === null) {
      setError(
        'Введите корректную сумму больше нуля, не более двух знаков после запятой.',
      )
      return
    }

    try {
      setSubmitting(true)
      setError(null)

      if (mode === 'deposit') {
        await depositToAccount(
          account.id,
          amountKopecks,
        )
      } else {
        const toAccountId = Number(
          destinationAccountId,
        )

        if (
          !Number.isInteger(
            toAccountId,
          ) ||
          toAccountId <= 0 ||
          toAccountId === account.id
        ) {
          setError(
            'Выберите счёт получателя.',
          )
          return
        }

        await transferBetweenAccounts(
          account.id,
          toAccountId,
          amountKopecks,
        )
      }

      await onCompleted()

      onClose()
    } catch (error) {
      setError(
        getApiErrorMessage(error),
      )
    } finally {
      setSubmitting(false)
    }
  }

  const isDeposit =
    mode === 'deposit'

  return (
    <div
      className="modal-backdrop"
      role="presentation"
      onMouseDown={(event) => {
        if (
          event.target ===
            event.currentTarget &&
          !submitting
        ) {
          onClose()
        }
      }}
    >
      <section
        className="operation-modal"
        role="dialog"
        aria-modal="true"
        aria-labelledby="operation-modal-title"
      >
        <header className="operation-modal__header">
          <div>
            <p className="dashboard__eyebrow">
              Счёт #{account.id}
            </p>

            <h2 id="operation-modal-title">
              {isDeposit
                ? 'Пополнение счёта'
                : 'Перевод средств'}
            </h2>
          </div>

          <button
            className="modal-close-button"
            type="button"
            aria-label="Закрыть"
            disabled={submitting}
            onClick={onClose}
          >
            ×
          </button>
        </header>

        <div className="operation-modal__balance">
          <span>
            Доступный баланс
          </span>

          <strong>
            {formatKopecks(
              account.balance_kopecks,
            )}
          </strong>
        </div>

        <form
          className="operation-form"
          onSubmit={handleSubmit}
        >
          {!isDeposit && (
            <label className="form-field">
              <span>
                Счёт получателя
              </span>

              <select
                value={
                  destinationAccountId
                }
                disabled={submitting}
                onChange={(event) =>
                  setDestinationAccountId(
                    event.target.value,
                  )
                }
              >
                <option value="">
                  Выберите счёт
                </option>

                {destinationAccounts.map(
                  (destinationAccount) => (
                    <option
                      key={
                        destinationAccount.id
                      }
                      value={
                        destinationAccount.id
                      }
                    >
                      Счёт #
                      {
                        destinationAccount.id
                      }
                    </option>
                  ),
                )}
              </select>
            </label>
          )}

          <label className="form-field">
            <span>
              Сумма, ₽
            </span>

            <input
              autoFocus
              type="text"
              inputMode="decimal"
              autoComplete="off"
              placeholder="Например, 500,00"
              value={amount}
              disabled={submitting}
              onChange={(event) =>
                setAmount(
                  event.target.value,
                )
              }
            />
          </label>

          {error && (
            <div
              className="error-message"
              role="alert"
            >
              {error}
            </div>
          )}

          {!isDeposit &&
            destinationAccounts.length ===
              0 && (
              <div className="info-message">
                Для перевода нужен ещё
                один банковский счёт.
              </div>
            )}

          <div className="operation-modal__actions">
            <button
              className="secondary-button"
              type="button"
              disabled={submitting}
              onClick={onClose}
            >
              Отмена
            </button>

            <button
              className="primary-button"
              type="submit"
              disabled={
                submitting ||
                (!isDeposit &&
                  destinationAccounts.length ===
                    0)
              }
            >
              {submitting
                ? 'Выполняем...'
                : isDeposit
                  ? 'Пополнить'
                  : 'Перевести'}
            </button>
          </div>
        </form>
      </section>
    </div>
  )
}