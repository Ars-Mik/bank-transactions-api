
CREATE TABLE IF NOT EXISTS accounts (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    balance_kopecks BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT accounts_balance_nonnegative
        CHECK (balance_kopecks >= 0)
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    kind VARCHAR(20) NOT NULL,

    amount_kopecks BIGINT NOT NULL,

    from_account_id BIGINT REFERENCES accounts(id),

    to_account_id BIGINT REFERENCES accounts(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT transactions_kind_valid
        CHECK (kind IN ('deposit', 'transfer')),

    CONSTRAINT transactions_amount_positive
        CHECK (amount_kopecks > 0),

    CONSTRAINT transactions_accounts_valid
        CHECK (
            (
                kind = 'deposit'
                AND from_account_id IS NULL
                AND to_account_id IS NOT NULL
            )
            OR
            (
                kind = 'transfer'
                AND from_account_id IS NOT NULL
                AND to_account_id IS NOT NULL
                AND from_account_id <> to_account_id
            )
        )
);

CREATE INDEX IF NOT EXISTS idx_transactions_from_account
ON transactions (from_account_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_transactions_to_account
ON transactions (to_account_id, created_at DESC);
