import {
  Link,
} from 'react-router-dom'

import type {
  Account,
} from '../types/banking'

import {
  formatKopecks,
} from '../utils/money'

interface AccountCardProps {
  account: Account
}

export function AccountCard({
  account,
}: AccountCardProps) {
  return (
    <Link
      className="account-card-link"
      to={`/accounts/${account.id}`}
    >
      <article className="account-card">
        <div className="account-card__header">
          <span className="account-card__label">
            Банковский счёт
          </span>

          <span className="account-card__id">
            #{account.id}
          </span>
        </div>

        <div className="account-card__balance">
          {formatKopecks(
            account.balance_kopecks,
          )}
        </div>

        <div className="account-card__caption">
          Текущий баланс
        </div>
      </article>
    </Link>
  )
}